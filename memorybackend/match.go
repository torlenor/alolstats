package memorybackend

import (
	"fmt"

	"github.com/torlenor/alolstats/riotclient"
)

// GetMatch retreives match data for given id
func (b *Backend) GetMatch(id uint64) (*riotclient.MatchDTO, error) {

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if match, ok := b.matches[id]; ok {
		return &match, nil
	}
	return nil, fmt.Errorf("Match with id=%d not found in storage backend", id)
}

// StoreMatch stores new match data
func (b *Backend) StoreMatch(data *riotclient.MatchDTO) error {
	b.log.Debugf("Storing Match id=%d in storage", data.GameID)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.matches[uint64(data.GameID)] = *data

	return nil
}
