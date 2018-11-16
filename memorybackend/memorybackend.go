// Package memorybackend provides a storage backend which stores its data in memory.
// This means the data is not permanently saved but gives a nice testing framework
// when no database for the other storages backends is available.
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
	matches      map[uint64]riotclient.Match

	log   *logrus.Entry
	mutex sync.Mutex
}

// NewBackend creates a new Memory Backend
func NewBackend() (*Backend, error) {
	b := &Backend{
		log:     logging.Get("Memory Storage Backend"),
		matches: make(map[uint64]riotclient.Match),
	}
	return b, nil
}

// GetMatchesByGameVersion returns all matches specific to a certain game version
func (b *Backend) GetMatchesByGameVersion(gameVersion string) (riotclient.Matches, error) {
	matches := riotclient.Matches{}
	for _, match := range b.matches {
		if match.GameVersion == gameVersion {
			matches.Matches = append(matches.Matches, match)
		}
	}

	return matches, nil
}
