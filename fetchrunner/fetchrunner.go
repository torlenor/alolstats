// Package fetchrunner is used to automatically trigger downloads of summoner and match data such that is lands in storage
package fetchrunner

import (
	"fmt"
	"sync"
	"time"

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
	config      config.FetchRunner
	storage     *storage.Storage
	log         *logrus.Entry
	stats       stats
	isStarted   bool
	workersWG   sync.WaitGroup
	stopWorkers chan struct{}
}

// NewFetchRunner creates a new FetchRunner
func NewFetchRunner(cfg config.FetchRunner, storage *storage.Storage) (*FetchRunner, error) {
	sr := &FetchRunner{
		storage:   storage,
		log:       logging.Get("FetchRunner"),
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
		f.log.Println("Starting FetchRunner")
		f.stopWorkers = make(chan struct{})
		go f.summonerMatchesWorker()
		f.isStarted = true
	} else {
		f.log.Println("FetchRunner already running")
	}
}

// Stop stops the FetchRunner and its workers
func (f *FetchRunner) Stop() {
	if f.isStarted {
		f.log.Println("Stopping FetchRunner")
		close(f.stopWorkers)
		f.workersWG.Wait()
		f.isStarted = false
	} else {
		f.log.Println("FetchRunner already stopped")
	}
}

func (f *FetchRunner) fetchSummonerMatches(accountID uint64) {
	// do work
	stop := false
	startIndex := uint32(0)
	endIndex := uint32(100)
	for !stop {
		matches, err := f.storage.GetMatchesByAccountID(accountID, startIndex, endIndex)
		if err != nil {
			f.log.Errorf("Error getting the current match list for Summoner: %s", err)
			continue
		}
		for _, match := range matches.Matches {
			f.storage.FetchAndStoreMatch(uint64(match.GameID))
		}
		if len(matches.Matches) == 0 || (endIndex+1) >= uint32(matches.TotalGames) {
			stop = true
		}
		startIndex += 100
		endIndex += 100
	}
}

func (f *FetchRunner) summonerMatchesWorker() {
	f.workersWG.Add(1)
	defer f.workersWG.Done()

	var nextUpdate time.Duration // := time.Duration(f.config.UpdateIntervalSummonerMatches)

	for {
		select {
		case <-f.stopWorkers:
			f.log.Printf("Stopping summonerMatchesWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			f.log.Debugf("Starting SummonerMatchesWorker run")

			start := time.Now()

			for _, accountID := range f.config.MatchesForSummonerAccountIDs {
				f.fetchSummonerMatches(accountID)
			}

			nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)

			elapsed := time.Since(start)
			f.log.Debugf("Finished SummonerMatchesWorker run. Took %s", elapsed)
		}
	}
}
