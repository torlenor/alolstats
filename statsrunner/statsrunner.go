// Package statsrunner is the statistics calculation part of ALoLStats and provides
// the calculation facilities and the API endpoints for retreiving the data.
package statsrunner

import (
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
	api.AttachModuleGet("/stats/champion/byid", sr.championByIDEndpoint)
	api.AttachModuleGet("/stats/champion/byname", sr.championByNameEndpoint)

	api.AttachModuleGet("/stats/plots/champion/byname", sr.championByNamePlotEndpoint)

}

// GetHandeledRequests gets the total number of api requests handeled by the StatsRunner since creating it
func (sr *StatsRunner) GetHandeledRequests() uint64 {
	return atomic.LoadUint64(&sr.stats.handledRequests)
}
