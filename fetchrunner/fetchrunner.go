// Package fetchrunner is used to automatically trigger downloads of summoner and match data such that is lands in storage
package fetchrunner

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/storage"
)

type stats struct {
	handledRequests uint64
}

// FetchRunner automatically fetches summoner and match data based on specified criteras
type FetchRunner struct {
	config            config.FetchRunner
	storage           *storage.Storage
	log               *logrus.Entry
	stats             stats
	isStarted         bool
	workersWG         sync.WaitGroup
	stopWorkers       chan struct{}
	shouldWorkersStop bool
}

// NewFetchRunner creates a new FetchRunner
func NewFetchRunner(cfg config.FetchRunner, storage *storage.Storage) (*FetchRunner, error) {
	name := fmt.Sprintf("FetchRunner [%s]", cfg.Region)
	sr := &FetchRunner{
		storage:   storage,
		log:       logging.Get(name),
		isStarted: false,
		workersWG: sync.WaitGroup{},
	}
	if cfg.UpdateIntervalSummonerMatches <= 0 {
		return nil, fmt.Errorf("The specified UpdateIntervalSummonerMatches is too small (%d min). Must be > 0 minutes", cfg.UpdateIntervalSummonerMatches)
	}
	sr.config = cfg

	return sr, nil
}

// Start starts the FetchRunner and its workers
func (f *FetchRunner) Start() {
	if !f.isStarted {
		f.log.Print("Starting FetchRunner")
		f.shouldWorkersStop = false
		f.stopWorkers = make(chan struct{})
		go f.summonerMatchesWorker()
		if f.config.UpdateIntervalFreeRotation > 0 {
			go f.freeRotationWorker()
		}
		f.isStarted = true
	} else {
		f.log.Print("FetchRunner already running")
	}
}

// Stop stops the FetchRunner and its workers
func (f *FetchRunner) Stop() {
	if f.isStarted {
		f.log.Print("Stopping FetchRunner")
		f.shouldWorkersStop = true
		close(f.stopWorkers)
		f.workersWG.Wait()
		f.isStarted = false
	} else {
		f.log.Printf("FetchRunner already stopped")
	}
}
