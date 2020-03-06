// Package statsrunner is the statistics calculation part of ALoLStats and provides
// the calculation facilities and the API endpoints for retreiving the data.
package statsrunner

import (
	"fmt"
	"sync"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
	"git.abyle.org/hps/alolstats/storage"
	"github.com/sirupsen/logrus"
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

	isStarted              bool
	workersWG              sync.WaitGroup
	stopWorkers            chan struct{}
	shouldWorkersStopMutex sync.RWMutex
	shouldWorkersStop      bool

	calculationMutex sync.Mutex
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

// GetHandeledRequests gets the total number of api requests handeled by the StatsRunner since creating it
func (sr *StatsRunner) GetHandeledRequests() uint64 {
	return atomic.LoadUint64(&sr.stats.handledRequests)
}

// Start starts the StatsRunner and its workers
func (sr *StatsRunner) Start() {
	if !sr.isStarted {
		sr.log.Println("Starting StatsRunner")
		sr.shouldWorkersStopMutex.Lock()
		sr.shouldWorkersStop = false
		sr.shouldWorkersStopMutex.Unlock()
		sr.stopWorkers = make(chan struct{})
		if sr.config.RunRScripts {
			sr.workersWG.Add(1)
			go sr.rScriptWorker()
		} else {
			sr.log.Info("Not running R scripts (deactivated in config)")
		}

		if sr.config.ChampionsStats.Enabled {
			sr.workersWG.Add(1)
			go sr.matchAnalysisWorker()
		}
		if sr.config.ItemsStats.Enabled {
			sr.workersWG.Add(1)
			go sr.itemWinRateWorker()
		}
		if sr.config.SummonerSpellsStats.Enabled {
			sr.workersWG.Add(1)
			go sr.summonerSpellsWorker()
		}
		if sr.config.RunesReforgedStats.Enabled {
			sr.workersWG.Add(1)
			go sr.runesReforgedWorker()
		}

		sr.isStarted = true
	} else {
		sr.log.Println("StatsRunner already running")
	}
}

// Stop stops the StatsRunner and its workers
func (sr *StatsRunner) Stop() {
	if sr.isStarted {
		sr.log.Println("Stopping StatsRunner")
		sr.shouldWorkersStopMutex.Lock()
		sr.shouldWorkersStop = true
		sr.shouldWorkersStopMutex.Unlock()
		close(sr.stopWorkers)
		sr.workersWG.Wait()
		sr.isStarted = false
	} else {
		sr.log.Println("StatsRunner already stopped")
	}
}
