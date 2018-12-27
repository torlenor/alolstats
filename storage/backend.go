package storage

import (
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// BackendRiotAPI contains function which directly store/retreive Riot API data (DTO structs)
type BackendRiotAPI interface {
	GetMatch(matchID uint64) (*riotclient.MatchDTO, error)
	StoreMatch(data *riotclient.MatchDTO) error

	GetSummonerByName(name string) (*Summoner, error)
	GetSummonerByNameTimeStamp(name string) time.Time
	GetSummonerBySummonerID(summonerID string) (*Summoner, error)
	GetSummonerBySummonerIDTimeStamp(summonerID string) time.Time
	GetSummonerByAccountID(accountID string) (*Summoner, error)
	GetSummonerByAccountIDTimeStamp(accountID string) time.Time
	StoreSummoner(data *Summoner) error
}

// Backend defines the interface for a storage backend like mongodb
type Backend interface {
	GetChampions() (*riotclient.ChampionList, error)
	GetChampionsTimeStamp() time.Time
	StoreChampions(championList *riotclient.ChampionList) error

	BackendRiotAPI

	GetFreeRotation() (*riotclient.FreeRotation, error)
	GetFreeRotationTimeStamp() time.Time
	StoreFreeRotation(freeRotation *riotclient.FreeRotation) error

	// Specialized fetching functions
	GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (*riotclient.Matches, error)
	GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queue uint64) (*riotclient.Matches, error)
	GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (*riotclient.Matches, error)

	GetStorageSummary() (Summary, error)
}
