package memorybackend

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/riotclient"
)

// Backend represents the Memory Backend
type Backend struct {
	championList riotclient.ChampionList
	freeRotation riotclient.FreeRotation

	log   *logrus.Entry
	mutex sync.Mutex
}

// NewBackend creates a new Memory Backend
func NewBackend() (*Backend, error) {
	b := &Backend{
		log: logging.Get("Memory Storage Backend"),
	}
	return b, nil
}
