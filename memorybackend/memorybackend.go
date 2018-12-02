// Package memorybackend provides a storage backend which stores its data in memory.
// This means the data is not permanently saved but gives a nice testing framework
// when no database for the other storages backends is available.
package memorybackend

import (
	"fmt"
	"sync"
	"time"

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

// GetMatchesByGameVersionAndChampionID returns all matches specific to a certain game version and champion id
func (b *Backend) GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (riotclient.Matches, error) {
	matches := riotclient.Matches{}
	for _, match := range b.matches {
		if match.GameVersion == gameVersion {
			valid := false
			for _, participant := range match.Participants {
				if uint64(participant.ChampionID) == championID {
					valid = true
				}
			}
			if valid {
				matches.Matches = append(matches.Matches, match)
			}
		}
	}

	return matches, nil
}

// GetMatchesByGameVersionChampionIDMapQueue returns all matches specific to a certain game version, champion id, map id and queue id
func (b *Backend) GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queueID uint64) (riotclient.Matches, error) {
	matches := riotclient.Matches{}
	for _, match := range b.matches {
		if match.GameVersion == gameVersion && uint64(match.MapID) == mapID && uint64(match.QueueID) == queueID {
			valid := false
			for _, participant := range match.Participants {
				if uint64(participant.ChampionID) == championID {
					valid = true
				}
			}
			if valid {
				matches.Matches = append(matches.Matches, match)
			}
		}
	}

	return matches, nil
}

// GetMatchesByGameVersionChampionIDMapBetweenQueueIDs returns all matches specific to a certain game version, champion id, map id and queue ids between and equal to ltequeue <= queueid <= gtequeue
func (b *Backend) GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (riotclient.Matches, error) {
	matches := riotclient.Matches{}
	for _, match := range b.matches {
		if match.GameVersion == gameVersion && uint64(match.MapID) == mapID && (uint64(match.QueueID) >= gtequeue && uint64(match.QueueID) <= ltequeue) {
			valid := false
			for _, participant := range match.Participants {
				if uint64(participant.ChampionID) == championID {
					valid = true
				}
			}
			if valid {
				matches.Matches = append(matches.Matches, match)
			}
		}
	}

	return matches, nil
}

func (b *Backend) GetSummonerByName(name string) (riotclient.Summoner, error) {
	return riotclient.Summoner{}, fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerByNameTimeStamp(name string) time.Time {
	return time.Time{}
}

func (b *Backend) GetSummonerBySummonerID(summonerID uint64) (riotclient.Summoner, error) {
	return riotclient.Summoner{}, fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerBySummonerIDTimeStamp(summonerID uint64) time.Time {
	return time.Time{}
}

func (b *Backend) GetSummonerByAccountID(accountID uint64) (riotclient.Summoner, error) {
	return riotclient.Summoner{}, fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerByAccountIDTimeStamp(accountID uint64) time.Time {
	return time.Time{}
}

func (b *Backend) StoreSummoner(data *riotclient.Summoner) error {
	return nil
}
