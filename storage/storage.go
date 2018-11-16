// Package storage is a main component of ALoLStats as it provides the interface between Storage Backend, Riot Client and serves
// as an interface for StatsRunner and FetchRunner
package storage

import (
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/api"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/matchfilereader"
	"github.com/torlenor/alolstats/riotclient"
)

// Backend defines the interface for a storage backend like sqlite
type Backend interface {
	GetChampions() (riotclient.ChampionList, error)
	GetChampionsTimeStamp() time.Time
	StoreChampions(championList riotclient.ChampionList) error

	GetFreeRotation() (riotclient.FreeRotation, error)
	GetFreeRotationTimeStamp() time.Time
	StoreFreeRotation(freeRotation riotclient.FreeRotation) error

	GetMatch(id uint64) (riotclient.Match, error)
	StoreMatch(data *riotclient.Match) error

	// Specialized fetching functions
	GetMatchesByGameVersion(gameVersion string) (riotclient.Matches, error)
}

type stats struct {
	handledRequests uint64
}

// Storage LoL data storage
type Storage struct {
	config     config.LoLStorage
	riotClient riotclient.Client
	log        *logrus.Entry
	stats      stats
	backend    Backend
}

// NewStorage creates a new Riot LoL API client
func NewStorage(cfg config.LoLStorage, riotClient riotclient.Client, backend Backend) (*Storage, error) {
	s := &Storage{
		config:     cfg,
		riotClient: riotClient,
		log:        logging.Get("Storage"),
		backend:    backend,
	}

	return s, nil
}

// RegisterAPI registers all endpoints from storage to the RestAPI
func (s *Storage) RegisterAPI(api *api.API) {
	api.AttachModuleGet("/champions", s.championsEndpoint)
	api.AttachModuleGet("/champion-rotations", s.freeRotationEndpoint)
	api.AttachModuleGet("/match", s.getMatchEndpoint)
	api.AttachModuleGet("/matches", s.getMatchesEndpoint)
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
	// TODO
}

// Stop stops the storage runners
func (s *Storage) Stop() {
	s.log.Println("Stopping Storage")
	// TODO
}

// GetHandeledRequests gets the total number of api requests handeled by the storage since creating it
func (s *Storage) GetHandeledRequests() uint64 {
	return atomic.LoadUint64(&s.stats.handledRequests)
}
