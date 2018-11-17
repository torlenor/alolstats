// Package statsrunner is the statistics calculation part of ALoLStats and provides
// the calculation facilities and the API endpoints for retreiving the data.
package statsrunner

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/api"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/storage"
)

type stats struct {
	handledRequests uint64
}

// StatsRunner calculates statistics and provides endpoints for the API
type StatsRunner struct {
	storage *storage.Storage
	log     *logrus.Entry
	stats   stats
}

type championStats struct {
	ChampionID      uint64 `json:"championid"`
	ChampionName    string `json:"championname"`
	SampleSize      uint64 `json:"samplesize"`
	GameVersion     string `json:"gameversion"`
	LanesPercentage struct {
		Mid    float64 `json:"mid"`
		Top    float64 `json:"top"`
		Jungle float64 `json:"jungle"`
		Bot    float64 `json:"bot"`
	} `json:"lanespercentage"`
}

// NewStatsRunner creates a new LoL StatsRunner
func NewStatsRunner(storage *storage.Storage) (*StatsRunner, error) {
	sr := &StatsRunner{
		storage: storage,
		log:     logging.Get("StatsRunner"),
	}

	return sr, nil
}

// RegisterAPI registers all endpoints from StatsRunner to the RestAPI
func (sr *StatsRunner) RegisterAPI(api *api.API) {
	api.AttachModuleGet("/stats/champion", sr.championEndpoint)
}

func (sr *StatsRunner) championEndpoint(w http.ResponseWriter, r *http.Request) {
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

	matches, err := sr.storage.GetStoredMatchesByGameVersion(gameVersion)
	if len(matches.Matches) == 0 {
		sr.log.Warnf("Error in getting matches for game version = %s", gameVersion)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var top, bot, mid, jungle, total uint64
	for _, match := range matches.Matches {
		if match.MapID != 11 {
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
					bot++
				default:
					// no data
					continue
				}
				total++
			}
		}
	}

	championStats := championStats{}
	championStats.ChampionID = champID
	championStats.SampleSize = total
	championStats.GameVersion = gameVersion
	championStats.LanesPercentage.Mid = float64(mid) / float64(total)
	championStats.LanesPercentage.Top = float64(top) / float64(total)
	championStats.LanesPercentage.Jungle = float64(jungle) / float64(total)
	championStats.LanesPercentage.Bot = float64(bot) / float64(total)

	champions := sr.storage.GetChampions()
	for _, val := range champions.Champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			championStats.ChampionName = val.Name
			break
		}
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

// GetHandeledRequests gets the total number of api requests handeled by the StatsRunner since creating it
func (sr *StatsRunner) GetHandeledRequests() uint64 {
	return atomic.LoadUint64(&sr.stats.handledRequests)
}
