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
	ChampionID    uint64 `json:"championid"`
	ChampionName  string `json:"championname"`
	GameVersion   string `json:"gameversion"`
	LaneRoleRatio struct {
		Mid        float64 `json:"mid"`
		Top        float64 `json:"top"`
		Jungle     float64 `json:"jungle"`
		BotAdc     float64 `json:"botadc"`
		BotSup     float64 `json:"botsup"`
		BotUnknown float64 `json:"botunknown"`
		Unknown    float64 `json:"unknown"`
		SampleSize uint64  `json:"samplesize"`
	} `json:"laneroleratio"`
}

func (sr *StatsRunner) getCHampionStatsByID(champID uint64, gameVersion string) (*championStats, error) {
	matches, err := sr.storage.GetStoredMatchesByGameVersion(gameVersion)
	if err != nil {
		return nil, fmt.Errorf("Could not get GetStoredMatchesByGameVersion: %s", err)
	}
	if len(matches.Matches) == 0 {
		return nil, fmt.Errorf("Error in getting matches for game version = %s", gameVersion)
	}

	var top, botAdc, botSup, botUnknown, mid, jungle, unknown, total uint64
	for _, match := range matches.Matches {
		if !(match.MapID == 11 && (match.QueueID == 420 || match.QueueID == 440)) {
			continue
		}
		for _, participant := range match.Participants {
			if uint64(participant.ChampionID) == champID {
				switch participant.Timeline.Lane {
				case "MID":
					fallthrough
				case "MIDDLE":
					mid++
				case "TOP":
					top++
				case "JUNGLE":
					jungle++
				case "BOT":
					fallthrough
				case "BOTTOM":
					switch participant.Timeline.Role {
					case "DUO_CARRY":
						botAdc++
					case "DUO_SUPPORT":
						botSup++
					default:
						botUnknown++
					}
				default:
					unknown++
				}
				total++
			}
		}
	}

	championStats := championStats{}
	championStats.ChampionID = champID
	championStats.GameVersion = gameVersion
	if total > 0 {
		championStats.LaneRoleRatio.Mid = float64(mid) / float64(total)
		championStats.LaneRoleRatio.Top = float64(top) / float64(total)
		championStats.LaneRoleRatio.Jungle = float64(jungle) / float64(total)
		championStats.LaneRoleRatio.BotAdc = float64(botAdc) / float64(total)
		championStats.LaneRoleRatio.BotSup = float64(botSup) / float64(total)
		championStats.LaneRoleRatio.BotUnknown = float64(botUnknown) / float64(total)
		championStats.LaneRoleRatio.Unknown = float64(unknown) / float64(total)
	}
	championStats.LaneRoleRatio.SampleSize = total

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
	sr.log.Debugln("Received Rest API Match request from", r.RemoteAddr)
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

	championStats, err := sr.getCHampionStatsByID(champID, gameVersion)
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
	sr.log.Debugln("Received Rest API Match request from", r.RemoteAddr)
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

	championStats, err := sr.getCHampionStatsByID(champID, gameVersion)
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
