// Package storage is a main component of ALoLStats as it provides the interface between Storage Backend, Riot Client and serves
// as an interface for StatsRunner and FetchRunner
package storage

import (
	"fmt"
	"path/filepath"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
	"git.abyle.org/hps/alolstats/matchfilereader"
	"git.abyle.org/hps/alolstats/riotclient"
	"github.com/sirupsen/logrus"
)

type stats struct {
	handledRequests uint64
}

// Storage LoL data storage
type Storage struct {
	config      config.LoLStorage
	riotClients map[string]riotclient.Client
	riotClient  riotclient.Client // Default riotclient, should be removed at some point
	log         *logrus.Entry
	stats       stats
	backend     Backend
}

// Summary gives an overview of the stored data in Storage/Backend
type Summary struct {
	NumberOfMatches   uint64 `json:"numberofmatches"`
	NumberOfSummoners uint64 `json:"numberofsummoners"`
	NumberOfChampions uint64 `json:"numberofchampions"`
}

// NewStorage creates a new Riot LoL API client
func NewStorage(cfg config.LoLStorage, riotClients map[string]riotclient.Client, backend Backend) (*Storage, error) {
	if client, ok := riotClients[cfg.DefaultRiotClient]; ok {
		s := &Storage{
			config:      cfg,
			riotClients: riotClients,
			riotClient:  client,
			log:         logging.Get("Storage"),
			backend:     backend,
		}

		return s, nil
	}
	return nil, fmt.Errorf("Error creating Storage. Requested default region RiotAPI does not exist: %s", cfg.DefaultRiotClient)
}

// Start starts the storage runners
func (s *Storage) Start() {
	s.log.Info("Starting Storage")
	if s.config.UseMatchFiles {
		s.log.Println("Reading match data from json files")
		files, err := filepath.Glob(s.config.MatchFileDir + "/*.json")
		if err != nil {
			s.log.Errorln("Error reading match files directory:", err)
		}
		for _, f := range files {
			matches, err := matchfilereader.ReadMatchesFile(f)
			if err != nil {
				s.log.Errorf("Error reading matches json file %s: %s", f, err)
			}
			for _, match := range matches.Matches {
				err = s.backend.StoreMatch(&match)
				if err != nil {
					s.log.Errorln("Error storing match data:", err)
				}
			}
		}
		s.log.Println("Finished reading match data from json files")
	}
}

// Stop stops the storage runners
func (s *Storage) Stop() {
	s.log.Println("Stopping Storage")
}

// GetHandeledRequests gets the total number of api requests handeled by the storage since creating it
func (s *Storage) GetHandeledRequests() uint64 {
	return atomic.LoadUint64(&s.stats.handledRequests)
}

// GameVersions struct is a list of game versions in the format major.minor, e.g., 8.24 or 9.1.
type GameVersions struct {
	Versions []string `json:"versions"`
}
