package memorybackend

import (
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetChampions gets the champions list from storage
func (b *Backend) GetChampions() (riotclient.ChampionList, error) {
	b.log.Debugln("Getting Champions List from storage")

	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.championList, nil
}

// GetChampionsTimeStamp gets the timestamp of the stored champions list
func (b *Backend) GetChampionsTimeStamp() time.Time {
	b.log.Debugln("Getting Champions List TimeStamp from storage")

	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Find oldest champ time
	oldest := time.Now()
	for _, champion := range b.championList.Champions {
		if oldest.Sub(champion.Timestamp) > 0 {
			oldest = champion.Timestamp
		}
	}

	return oldest
}

// StoreChampions stores a new champions list
func (b *Backend) StoreChampions(championList riotclient.ChampionList) error {
	b.log.Debugln("Storing new Champions List in storage")

	b.mutex.Lock()
	b.championList = championList
	b.mutex.Unlock()

	return nil
}
