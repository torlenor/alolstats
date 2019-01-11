package statsrunner

import (
	"fmt"
	"strconv"
	"time"

	"github.com/torlenor/alolstats/storage"
)

func (sr *StatsRunner) getChampionStatsByID(champID uint64, majorVersion uint32, minorVersion uint32) (*storage.ChampionStats, error) {

	majorMinor := fmt.Sprintf("%d\\.%d\\.", majorVersion, minorVersion)
	gameVersion := fmt.Sprintf("%d.%d", majorVersion, minorVersion)

	matches, err := sr.storage.GetStoredMatchesByGameVersionChampionIDMapBetweenQueueIDs(majorMinor, champID, 11, 440, 400)
	if err != nil {
		return nil, fmt.Errorf("Could not get GetStoredMatchesByGameVersionAndChampionID: %s", err)
	}
	if len(matches.Matches) == 0 {
		return nil, fmt.Errorf("Error in getting matches for game version = %s and Champion ID %d", majorMinor, champID)
	}

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

	championStats := storage.ChampionStats{}
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
		storage.LaneRolePercentage{
			Lane: "TOP",
			Role: "Solo",

			Percentage: float64(topWins+topLosses) / float64(total) * 100.0,
			Wins:       topWins,
			NGames:     topWins + topLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "MIDDLE",
			Role: "Solo",

			Percentage: float64(midWins+midLosses) / float64(total) * 100.0,
			Wins:       midWins,
			NGames:     midWins + midLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "JUNGLE",
			Role: "Solo",

			Percentage: float64(jungleWins+jungleLosses) / float64(total) * 100.0,
			Wins:       jungleWins,
			NGames:     jungleWins + jungleLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "BOT",
			Role: "Carry",

			Percentage: float64(botCarryWins+botCarryLosses) / float64(total) * 100.0,
			Wins:       botCarryWins,
			NGames:     botCarryWins + botCarryLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "BOT",
			Role: "Support",

			Percentage: float64(botSupWins+botSupLosses) / float64(total) * 100.0,
			Wins:       botSupWins,
			NGames:     botSupWins + botSupLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "BOT",
			Role: "Unknown",

			Percentage: float64(botUnknownWins+botUnknownLosses) / float64(total) * 100.0,
			Wins:       botUnknownWins,
			NGames:     botUnknownWins + botUnknownLosses,
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "UNKNOWN",
			Role: "Unknown",

			Percentage: float64(unknownWins+unknownLosses) / float64(total) * 100.0,
			Wins:       unknownWins,
			NGames:     unknownWins + unknownLosses,
		},
	)

	champions := sr.storage.GetChampions(false)
	for _, val := range champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			championStats.ChampionName = val.Name
			championStats.ChampionRealID = val.ID
			break
		}
	}

	// Plotly data
	// type laneRolePercentagePlotly struct {
	// 	X []string `json:"x"` // ['TOP', 'MIDDLE', 'JUNGLE', 'BOT', 'UNKNOWN'],
	// 	Y []string `json:"x"` // [2.2058823529411766, 2.941176470588235, 0.7352941176470588, 0, 0],

	// 	Name string `json:"name"` // 'Solo',
	// 	Type string `json:"type"` // 'bar'
	// }

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Uknown",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				0.0,
				0.0,
				0.0,
				float64(botUnknownWins+botUnknownLosses) / float64(total) * 100.0,
				float64(unknownWins+unknownLosses) / float64(total) * 100.0,
			},
		},
	)

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Support",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				0.0,
				0.0,
				0.0,
				float64(botSupWins+botSupLosses) / float64(total) * 100.0,
				0.0,
			},
		},
	)

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Carry",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				0.0,
				0.0,
				0.0,
				float64(botCarryWins+botCarryLosses) / float64(total) * 100.0,
				0.0,
			},
		},
	)

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Solo",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				float64(topWins+topLosses) / float64(total) * 100.0,
				float64(midWins+midLosses) / float64(total) * 100.0,
				float64(jungleWins+jungleLosses) / float64(total) * 100.0,
				0.0,
				0.0,
			},
		},
	)

	championStats.Timestamp = time.Now()

	// elapsed = time.Since(start)
	// sr.log.Debugf("Finished Champion Stats calculation. Took %s", elapsed)

	return &championStats, nil
}

// func (sr *StatsRunner) championByNamePlotEndpoint(w http.ResponseWriter, r *http.Request) {
// 	sr.log.Debugln("Received Rest API championByName request from", r.RemoteAddr)
// 	var gameVersion string

// 	var championName string
// 	if val, ok := r.URL.Query()["name"]; ok {
// 		if len(val) == 0 {
// 			sr.log.Warnf("name parameter was empty in request")
// 			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 			return
// 		}
// 		championName = val[0]
// 	}

// 	if val, ok := r.URL.Query()["gameversion"]; ok {
// 		if len(val) == 0 {
// 			sr.log.Warnf("gameversion parameter was empty in request.")
// 			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 			return
// 		}
// 		gameVersion = val[0]
// 	}

// 	champions := sr.storage.GetChampions(false)
// 	champID := uint64(0)
// 	var champRealID string
// 	for _, val := range champions {
// 		if strings.ToLower(val.ID) == strings.ToLower(championName) {
// 			id, err := strconv.ParseUint(val.Key, 10, 32)
// 			if err != nil {
// 				sr.log.Warnf("Could not convert value %s to Champion ID", val.Key)
// 				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 				return
// 			}
// 			champID = id
// 			champRealID = val.ID
// 			continue
// 		}
// 	}

// 	if len(champRealID) == 0 {
// 		sr.log.Warnf("Could not convert provided champ name %s to Champion ID", championName)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	var path = sr.config.RPlotsOutputPath
// 	var imageName = fmt.Sprintf("champion_role_%s_%d_%s.png", champRealID, champID, gameVersion)
// 	filePath := path + string(filepath.Separator) + imageName

// 	img, err := os.Open(filePath)
// 	if err != nil {
// 		sr.log.Warnf("Could not get plot for requested Champion %s and game version %s: %s", champRealID, gameVersion, err)
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}
// 	defer img.Close()
// 	w.Header().Set("Content-Type", "image/png")
// 	io.Copy(w, img)

// 	atomic.AddUint64(&sr.stats.handledRequests, 1)
// }
