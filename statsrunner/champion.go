package statsrunner

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/torlenor/alolstats/utils"
)

type laneRolePercentage struct {
	Lane string `json:"lane"`
	Role string `json:"role"`

	Percentage float64 `json:"percentage"`
	Wins       uint32  `json:"wins"`
	NGames     uint32  `json:"ngames"`
}

type championStats struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`
	GameVersion    string `json:"gameversion"`

	Timestamp time.Time `json:"timestamp"`

	SampleSize uint64 `json:"samplesize"`

	AvgK    float64 `json:"averagekills"`
	StdDevK float64 `json:"stddevkills"`
	MedianK float64 `json:"mediankills"`

	AvgD    float64 `json:"averagedeaths"`
	StdDevD float64 `json:"stddevdeaths"`
	MedianD float64 `json:"mediandeaths"`

	AvgA    float64 `json:"averageassists"`
	StdDevA float64 `json:"stddevassists"`
	MedianA float64 `json:"medianassists"`

	WinLossRatio float64 `json:"winlossratio"`

	LaneRolePercentage []laneRolePercentage `json:"lanerolepercentage"`
}

func (sr *StatsRunner) getChampionStatsByID(champID uint64, gameVersion string) (*championStats, error) {

	start := time.Now()

	matches, err := sr.storage.GetStoredMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion, champID, 11, 440, 400)
	if err != nil {
		return nil, fmt.Errorf("Could not get GetStoredMatchesByGameVersionAndChampionID: %s", err)
	}
	if len(matches.Matches) == 0 {
		return nil, fmt.Errorf("Error in getting matches for game version = %s and Champion ID %d", gameVersion, champID)
	}

	elapsed := time.Since(start)
	sr.log.Debugf("Got Matches from storage for Champion Stats calculation. Took %s", elapsed)

	var total uint64
	var kills, deaths, assists []float64
	var wins, losses uint64

	var topWins uint32
	var topLosses uint32
	var jungleWins uint32
	var jungleLosses uint32
	var midWins uint32
	var midLosses uint32
	var botCarryWins uint32
	var botCarryLosses uint32
	var botSupWins uint32
	var botSupLosses uint32
	var botUnknownWins uint32
	var botUnknownLosses uint32
	var unknownWins uint32
	var unknownLosses uint32

	for _, match := range matches.Matches {
		if !(match.MapID == 11 && (match.QueueID >= 400 && match.QueueID < 450)) {
			continue
		}
		for _, participant := range match.Participants {
			if uint64(participant.ChampionID) == champID {

				kills = append(kills, float64(participant.Stats.Kills))
				deaths = append(deaths, float64(participant.Stats.Deaths))
				assists = append(assists, float64(participant.Stats.Assists))

				if participant.Stats.Win {
					wins++
				} else {
					losses++
				}

				total++

				if participant.Timeline.Lane == "TOP" {
					if participant.Stats.Win == true {
						topWins = topWins + 1
					} else {
						topLosses = topLosses + 1
					}
				} else if participant.Timeline.Lane == "JUNGLE" {
					if participant.Stats.Win == true {
						jungleWins = jungleWins + 1
					} else {
						jungleLosses = jungleLosses + 1
					}
				} else if participant.Timeline.Lane == "MIDDLE" {
					if participant.Stats.Win == true {
						midWins = midWins + 1
					} else {
						midLosses = midLosses + 1
					}
				} else if participant.Timeline.Lane == "BOTTOM" {
					if participant.Timeline.Role == "DUO_CARRY" {
						if participant.Stats.Win == true {
							botCarryWins = botCarryWins + 1
						} else {
							botCarryLosses = botCarryLosses + 1
						}
					} else if participant.Timeline.Role == "DUO_SUPPORT" {
						if participant.Stats.Win == true {
							botSupWins = botSupWins + 1
						} else {
							botSupLosses = botSupLosses + 1
						}
					} else {
						if participant.Stats.Win == true {
							botUnknownWins = botUnknownWins + 1
						} else {
							botUnknownLosses = botUnknownLosses + 1
						}
					}
				} else {
					if participant.Stats.Win == true {
						unknownWins = unknownWins + 1
					} else {
						unknownLosses = unknownLosses + 1
					}
				}

			}
		}
	}

	championStats := championStats{}
	championStats.ChampionID = champID
	championStats.GameVersion = gameVersion
	championStats.SampleSize = total

	championStats.AvgK, championStats.StdDevK = calcMeanStdDev(kills, nil)
	championStats.AvgD, championStats.StdDevD = calcMeanStdDev(deaths, nil)
	championStats.AvgA, championStats.StdDevA = calcMeanStdDev(assists, nil)

	championStats.MedianK = calcMedian(kills, nil)
	championStats.MedianD = calcMedian(deaths, nil)
	championStats.MedianA = calcMedian(assists, nil)

	if losses > 0 {
		championStats.WinLossRatio = float64(wins) / float64(losses)
	} else if losses == 0 && wins > 0 {
		championStats.WinLossRatio = 1.0
	} else {
		championStats.WinLossRatio = 0
	}

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		laneRolePercentage{
			Lane: "TOP",
			Role: "Solo",

			Percentage: float64(topWins+topLosses) / float64(total) * 100.0,
			Wins:       topWins,
			NGames:     topWins + topLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		laneRolePercentage{
			Lane: "MIDDLE",
			Role: "Solo",

			Percentage: float64(midWins+midLosses) / float64(total) * 100.0,
			Wins:       midWins,
			NGames:     midWins + midLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		laneRolePercentage{
			Lane: "JUNGLE",
			Role: "Solo",

			Percentage: float64(jungleWins+jungleLosses) / float64(total) * 100.0,
			Wins:       jungleWins,
			NGames:     jungleWins + jungleLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		laneRolePercentage{
			Lane: "BOT",
			Role: "Carry",

			Percentage: float64(botCarryWins+botCarryLosses) / float64(total) * 100.0,
			Wins:       botCarryWins,
			NGames:     botCarryWins + botCarryLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		laneRolePercentage{
			Lane: "BOT",
			Role: "Support",

			Percentage: float64(botSupWins+botSupLosses) / float64(total) * 100.0,
			Wins:       botSupWins,
			NGames:     botSupWins + botSupLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		laneRolePercentage{
			Lane: "BOT",
			Role: "Unknown",

			Percentage: float64(botUnknownWins+botUnknownLosses) / float64(total) * 100.0,
			Wins:       botUnknownWins,
			NGames:     botUnknownWins + botUnknownLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		laneRolePercentage{
			Lane: "UNKNOWN",
			Role: "Unknown",

			Percentage: float64(unknownWins+unknownLosses) / float64(total) * 100.0,
			Wins:       unknownWins,
			NGames:     unknownWins + unknownLosses,
		},
	)

	champions := sr.storage.GetChampions()
	for _, val := range champions.Champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			championStats.ChampionName = val.Name
			championStats.ChampionRealID = val.ID
			break
		}
	}

	championStats.Timestamp = time.Now()

	elapsed = time.Since(start)
	sr.log.Debugf("Finished Champion Stats calculation. Took %s", elapsed)

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
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "Bad Request - Could not get data for the gameversion and champion id specified"),
			http.StatusBadRequest)
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

func (sr *StatsRunner) championByNamePlotEndpoint(w http.ResponseWriter, r *http.Request) {
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
	var champRealID string
	for _, val := range champions.Champions {
		if strings.ToLower(val.ID) == strings.ToLower(championName) {
			id, err := strconv.ParseUint(val.Key, 10, 32)
			if err != nil {
				sr.log.Warnf("Could not convert value %s to Champion ID", val.Key)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			champID = id
			champRealID = val.ID
			continue
		}
	}

	if len(champRealID) == 0 {
		sr.log.Warnf("Could not convert provided champ name %s to Champion ID", championName)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var path = sr.config.RPlotsOutputPath
	var imageName = fmt.Sprintf("champion_role_%s_%d_%s.png", champRealID, champID, gameVersion)
	filePath := path + string(filepath.Separator) + imageName

	img, err := os.Open(filePath)
	if err != nil {
		sr.log.Warnf("Could not get plot for requested Champion %s and game version %s: %s", champRealID, gameVersion, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, img)

	atomic.AddUint64(&sr.stats.handledRequests, 1)
}

func (sr *StatsRunner) championStatsListingEndpoint(w http.ResponseWriter, r *http.Request) {
	t := template.New("champions list")
	t, _ = t.Parse(`
			<head>
			<h1>Champion Stats</h1>
			<style>
			// table {
			// 	width:100%;
			// }
			table, th, td {
				border: 1px solid black;
				border-collapse: collapse;
			}
			th, td {
				padding: 15px;
				text-align: left;
			}
			table#t01 tr:nth-child(even) {
				background-color: #eee;
			}
			table#t01 tr:nth-child(odd) {
			   background-color: #fff;
			}
			table#t01 th {
				background-color: black;
				color: white;
			}
			</style>
			</head>
			<body>
			<table id="t01">
			<tr>
			  <th>Champion</th>
			  <th>Role Distribution</th> 
			</tr>
			{{range .Champions}}
			<tr>
			  <td>{{.Name}}</td>
			  <td><a href="/v1/stats/plots/champion/byname?name={{.ID}}&gameversion=8.23">Path 8.23</a> <a href="/v1/stats/plots/champion/byname?name={{.ID}}&gameversion=8.24">Path 8.24</a> </td>
			</tr>
			{{end}}
		  </table>

			</body>
			`)

	champions := sr.storage.GetChampions()

	err := t.Execute(w, champions)
	if err != nil {
		sr.log.Errorf("Error in using the Template for building the Stats overview: %s", err)
	}
}
