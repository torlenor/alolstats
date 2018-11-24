package statsrunner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

type championStats struct {
	ChampionID   uint64 `json:"championid"`
	ChampionName string `json:"championname"`
	GameVersion  string `json:"gameversion"`

	SampleSize uint64 `json:"samplesize"`

	LaneRelativeFrequency     map[string]float64 `json:"lanerelativefrequency"`
	RoleRelativeFrequency     map[string]float64 `json:"rolerelativefrequency"`
	LaneRoleRelativeFrequency map[string]float64 `json:"lanerolerelativefrequency"`

	AvgK    float64 `json:"averagekills"`
	StdDevK float64 `json:"stddevkills"`

	AvgD    float64 `json:"averagedeaths"`
	StdDevD float64 `json:"stddevdeaths"`

	AvgA    float64 `json:"averageassists"`
	StdDevA float64 `json:"stddevassists"`

	WinLossRatio float64 `json:"winlossratio"`
}

func (sr *StatsRunner) getChampionStatsByID(champID uint64, gameVersion string) (*championStats, error) {
	matches, err := sr.storage.GetStoredMatchesByGameVersionAndChampionID(gameVersion, champID)
	if err != nil {
		return nil, fmt.Errorf("Could not get GetStoredMatchesByGameVersionAndChampionID: %s", err)
	}
	if len(matches.Matches) == 0 {
		return nil, fmt.Errorf("Error in getting matches for game version = %s and Champion ID %d", gameVersion, champID)
	}

	var total uint64
	var laneObs, roleObs, laneRoleObs []string
	var kills, deaths, assists []float64
	var wins, losses uint64
	for _, match := range matches.Matches {
		if !(match.MapID == 11 && (match.QueueID == 420 || match.QueueID == 440)) {
			continue
		}
		for _, participant := range match.Participants {
			if uint64(participant.ChampionID) == champID {
				laneObs = append(laneObs, participant.Timeline.Lane)
				roleObs = append(roleObs, participant.Timeline.Role)
				laneRoleObs = append(laneRoleObs, participant.Timeline.Lane+":"+participant.Timeline.Role)

				kills = append(kills, float64(participant.Stats.Kills))
				deaths = append(deaths, float64(participant.Stats.Deaths))
				assists = append(assists, float64(participant.Stats.Assists))

				if participant.Stats.Win {
					wins++
				} else {
					losses++
				}

				total++
			}
		}
	}

	championStats := championStats{}
	championStats.ChampionID = champID
	championStats.GameVersion = gameVersion
	championStats.SampleSize = total
	championStats.LaneRelativeFrequency = calcRelativeFrequency(laneObs)
	championStats.RoleRelativeFrequency = calcRelativeFrequency(roleObs)
	championStats.LaneRoleRelativeFrequency = calcRelativeFrequency(laneRoleObs)

	championStats.AvgK, championStats.StdDevK = calcMeanStdDev(kills, nil)
	championStats.AvgD, championStats.StdDevD = calcMeanStdDev(deaths, nil)
	championStats.AvgA, championStats.StdDevA = calcMeanStdDev(assists, nil)

	championStats.WinLossRatio = float64(wins) / float64(losses)

	champions := sr.storage.GetChampions()
	for _, val := range champions.Champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			championStats.ChampionName = val.Name
			break
		}
	}

	return &championStats, nil
}

func (sr *StatsRunner) championByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	sr.log.Debugln("Received Rest API championByID request from", r.RemoteAddr)
	var champID uint64
	var gameVersion string

	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			sr.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		var err error
		champID, err = strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			sr.log.Warnf("Could not convert value %s to GameID", val)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	if val, ok := r.URL.Query()["gameversion"]; ok {
		if len(val) == 0 {
			sr.log.Warnf("gameversion parameter was empty in request.")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		gameVersion = val[0]
	}

	championStats, err := sr.getChampionStatsByID(champID, gameVersion)
	if err != nil {
		sr.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(championStats)
	if err != nil {
		sr.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&sr.stats.handledRequests, 1)
}

func (sr *StatsRunner) championByNameEndpoint(w http.ResponseWriter, r *http.Request) {
	sr.log.Debugln("Received Rest API championByName request from", r.RemoteAddr)
	var gameVersion string

	var championName string
	if val, ok := r.URL.Query()["name"]; ok {
		if len(val) == 0 {
			sr.log.Warnf("name parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		championName = val[0]
	}

	if val, ok := r.URL.Query()["gameversion"]; ok {
		if len(val) == 0 {
			sr.log.Warnf("gameversion parameter was empty in request.")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		gameVersion = val[0]
	}

	champions := sr.storage.GetChampions()
	champID := uint64(0)
	for _, val := range champions.Champions {
		if strings.ToLower(val.ID) == strings.ToLower(championName) {
			id, err := strconv.ParseUint(val.Key, 10, 32)
			if err != nil {
				sr.log.Warnf("Could not convert value %s to Champion ID", val.Key)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			champID = id
			continue
		}
	}

	stats, err := sr.getChampionStatsByID(champID, gameVersion)
	if err != nil {
		sr.log.Errorf("Error in getting stats for champion: %s", err)
		io.WriteString(w, `{"Error": "Could not get data, check request"}`)
		return
	}

	out, err := json.Marshal(stats)
	if err != nil {
		sr.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&sr.stats.handledRequests, 1)
}
