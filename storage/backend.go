package storage

import (
	"time"

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

	GetSummonerByName(name string) (riotclient.Summoner, error)
	GetSummonerByNameTimeStamp(name string) time.Time
	GetSummonerBySummonerID(summonerID uint64) (riotclient.Summoner, error)
	GetSummonerBySummonerIDTimeStamp(summonerID uint64) time.Time
	GetSummonerByAccountID(accountID uint64) (riotclient.Summoner, error)
	GetSummonerByAccountIDTimeStamp(accountID uint64) time.Time
	StoreSummoner(data *riotclient.Summoner) error

	// Specialized fetching functions
	GetMatchesByGameVersion(gameVersion string) (riotclient.Matches, error)
	GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (riotclient.Matches, error)
	GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queue uint64) (riotclient.Matches, error)
	GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (riotclient.Matches, error)
}
