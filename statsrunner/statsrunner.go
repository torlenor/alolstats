// Package statsrunner is the statistics calculation part of ALoLStats and provides
// the calculation facilities and the API endpoints for retreiving the data.
package statsrunner

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/api"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/storage"
)

type stats struct {
	handledRequests uint64
}

// StatsRunner calculates statistics and provides endpoints for the API
type StatsRunner struct {
	config  config.StatsRunner
	storage *storage.Storage
	log     *logrus.Entry
	stats   stats

	isStarted         bool
	workersWG         sync.WaitGroup
	stopWorkers       chan struct{}
	shouldWorkersStop bool
}

// NewStatsRunner creates a new LoL StatsRunner
func NewStatsRunner(cfg config.StatsRunner, storage *storage.Storage) (*StatsRunner, error) {
	sr := &StatsRunner{
		storage:   storage,
		log:       logging.Get("StatsRunner"),
		isStarted: false,
		workersWG: sync.WaitGroup{},
	}

	if cfg.RunRScripts == true && cfg.RScriptsUpdateInterval <= 0 {
		return nil, fmt.Errorf("The specified RScriptsUpdateInterval is too small (%d min). Must be > 0 minutes or deactivate RunRScripts", cfg.RScriptsUpdateInterval)
	}
	sr.config = cfg

	return sr, nil
}

// RegisterAPI registers all endpoints from StatsRunner to the RestAPI
func (sr *StatsRunner) RegisterAPI(api *api.API) {
	// api.AttachModuleGet("/stats/champion/byid", sr.championByIDEndpoint)
	// api.AttachModuleGet("/stats/champion/byname", sr.championByNameEndpoint)

	// api.AttachModuleGet("/stats/plots/champion/byname", sr.championByNamePlotEndpoint)

}

// GetHandeledRequests gets the total number of api requests handeled by the StatsRunner since creating it
func (sr *StatsRunner) GetHandeledRequests() uint64 {
	return atomic.LoadUint64(&sr.stats.handledRequests)
}

// Start starts the StatsRunner and its workers
func (sr *StatsRunner) Start() {
	if !sr.isStarted {
		sr.log.Println("Starting StatsRunner")
		sr.shouldWorkersStop = false
		sr.stopWorkers = make(chan struct{})
		if sr.config.RunRScripts {
			go sr.rScriptWorker()
		} else {
			sr.log.Info("Not running R scripts (deactivated in config)")
		}

		go sr.matchAnalysisWorker()

		sr.isStarted = true
	} else {
		sr.log.Println("StatsRunner already running")
	}
}

// Stop stops the StatsRunner and its workers
func (sr *StatsRunner) Stop() {
	if sr.isStarted {
		sr.log.Println("Stopping StatsRunner")
		sr.shouldWorkersStop = true
		close(sr.stopWorkers)
		sr.workersWG.Wait()
		sr.isStarted = false
	} else {
		sr.log.Println("StatsRunner already stopped")
	}
}
