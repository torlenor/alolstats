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
	"github.com/torlenor/alolstats/storage"
)

// Backend represents the Memory Backend
type Backend struct {
	championList riotclient.ChampionsList
	freeRotation riotclient.FreeRotation
	matches      map[uint64]riotclient.MatchDTO

	log   *logrus.Entry
	mutex sync.Mutex
}

// NewBackend creates a new Memory Backend
func NewBackend() (*Backend, error) {
	b := &Backend{
		log:          logging.Get("Memory Storage Backend"),
		matches:      make(map[uint64]riotclient.MatchDTO),
		championList: make(riotclient.ChampionsList),
	}
	return b, nil
}

// GetMatchesByGameVersionAndChampionID returns all matches specific to a certain game version and champion id
func (b *Backend) GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (*riotclient.Matches, error) {
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

	return &matches, nil
}

// GetMatchesByGameVersionChampionIDMapQueue returns all matches specific to a certain game version, champion id, map id and queue id
func (b *Backend) GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queueID uint64) (*riotclient.Matches, error) {
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

	return &matches, nil
}

// GetMatchesByGameVersionChampionIDMapBetweenQueueIDs returns all matches specific to a certain game version, champion id, map id and queue ids between and equal to ltequeue <= queueid <= gtequeue
func (b *Backend) GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (*riotclient.Matches, error) {
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

	return &matches, nil
}

func (b *Backend) GetMatchTimeLine(matchID uint64) (*riotclient.MatchTimelineDTO, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) StoreMatchTimeLine(data *riotclient.MatchTimelineDTO) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerByName(name string) (*storage.Summoner, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerByNameTimeStamp(name string) time.Time {
	return time.Time{}
}

func (b *Backend) GetSummonerBySummonerID(summonerID string) (*storage.Summoner, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerBySummonerIDTimeStamp(summonerID string) time.Time {
	return time.Time{}
}

func (b *Backend) GetSummonerByAccountID(accountID string) (*storage.Summoner, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerByAccountIDTimeStamp(accountID string) time.Time {
	return time.Time{}
}

func (b *Backend) GetSummonerByPUUID(PUUID string) (*storage.Summoner, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetSummonerByPUUIDTimeStamp(PUUID string) time.Time {
	return time.Time{}
}

func (b *Backend) StoreSummoner(data *storage.Summoner) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) GetLeagueByQueue(league string, queue string) (*riotclient.LeagueListDTO, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetLeagueByQueueTimeStamp(league string, queue string) (time.Time, error) {
	return time.Time{}, fmt.Errorf("Not implemented")
}

func (b *Backend) StoreLeague(*riotclient.LeagueListDTO) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) GetLeaguesForSummoner(summonerName string) (*storage.SummonerLeagues, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetLeaguesForSummonerBySummonerID(summonerID string) (*storage.SummonerLeagues, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) StoreLeaguesForSummoner(*storage.SummonerLeagues) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) GetLeaguesForSummonerTimeStamp(summonerName string) (time.Time, error) {
	return time.Time{}, fmt.Errorf("Not implemented")
}

func (b *Backend) GetLeaguesForSummonerBySummonerIDTimeStamp(summonerID string) (time.Time, error) {
	return time.Time{}, fmt.Errorf("Not implemented")
}

// GetStorageSummary returns stats about the stored elements in the Backend
func (b *Backend) GetStorageSummary() (storage.Summary, error) {
	summary := storage.Summary{}
	summary.NumberOfMatches = uint64(len(b.matches))
	summary.NumberOfChampions = uint64(len(b.championList))
	summary.NumberOfSummoners = 0 // not implemented

	return summary, nil
}

func (b *Backend) StoreChampionStats(data *storage.ChampionStatsStorage) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) GetChampionStatsByChampionIDGameVersion(championID string, gameVersion string) (*storage.ChampionStatsStorage, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetChampionStatsByChampionIDGameVersionTier(championID string, gameVersion string, tier string) (*storage.ChampionStatsStorage, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetKnownGameVersions() (*storage.GameVersions, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) StoreKnownGameVersions(gameVersions *storage.GameVersions) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) StoreChampionStatsSummary(statsSummary *storage.ChampionStatsSummaryStorage) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) GetChampionStatsSummaryByGameVersionTier(gameVersion string, tier string) (*storage.ChampionStatsSummaryStorage, error) {
	return nil, fmt.Errorf("Not implemented")
}
