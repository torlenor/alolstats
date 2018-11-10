package storage

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/api"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/riotclient"
)

// Backend defines the interface for a storage backend like sqlite
type Backend interface {
	GetChampions() riotclient.ChampionList
	GetChampionsTimeStamp() time.Time
	StoreChampions(championList riotclient.ChampionList)

	GetFreeRotation() riotclient.FreeRotation
	GetFreeRotationTimeStamp() time.Time
	StoreFreeRotation(freeRotation riotclient.FreeRotation)
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
}

// Start starts the storage runners
func (s *Storage) Start() {
	s.log.Println("Starting Storage")
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
