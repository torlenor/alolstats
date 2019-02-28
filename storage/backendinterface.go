package storage

import (
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// BackendChampion defines an interface to store/retrieve Champion data from Storage Backend
type BackendChampion interface {
	GetChampions() (riotclient.ChampionsList, error)
	GetChampionsTimeStamp() time.Time

	StoreChampions(championsList riotclient.ChampionsList) error
}

// BackendFreeRotation defines an interfce to store/retrieve the Champions Free Rotation from Storage Backend
type BackendFreeRotation interface {
	GetFreeRotation() (*riotclient.FreeRotation, error)
	GetFreeRotationTimeStamp() time.Time
	StoreFreeRotation(freeRotation *riotclient.FreeRotation) error
}

// BackendLeague defines an interface to store/retrieve League data from Storage Backend
type BackendLeague interface {
	GetLeagueByQueue(league string, queue string) (*riotclient.LeagueListDTO, error)
	GetLeagueByQueueTimeStamp(league string, queue string) (time.Time, error)

	StoreLeague(*riotclient.LeagueListDTO) error

	GetLeaguesForSummoner(summonerName string) (*SummonerLeagues, error)
	GetLeaguesForSummonerTimeStamp(summonerName string) (time.Time, error)

	GetLeaguesForSummonerBySummonerID(summonerID string) (*SummonerLeagues, error)
	GetLeaguesForSummonerBySummonerIDTimeStamp(summonerID string) (time.Time, error)

	StoreLeaguesForSummoner(*SummonerLeagues) error
}

// BackendMatch defines an interface to store/retrieve Match data from Storage Backend
// Matches have no TimeStamp as they are always valid
type BackendMatch interface {
	GetMatch(matchID uint64) (*riotclient.MatchDTO, error)
	StoreMatch(data *riotclient.MatchDTO) error

	GetMatchTimeLine(matchID uint64) (*riotclient.MatchTimelineDTO, error)
	StoreMatchTimeLine(data *riotclient.MatchTimelineDTO) error

	// Specialized fetching functions
	GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (*riotclient.Matches, error)
	GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queue uint64) (*riotclient.Matches, error)
	GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (*riotclient.Matches, error)

	// Cursor fetching functions
	GetMatchesCursorByGameVersion(gameVersion string) (QueryCursor, error)
	GetMatchesCursorByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (QueryCursor, error)
	GetMatchesCursorByGameVersionMapBetweenQueueIDs(gameVersion string, mapID uint64, ltequeue uint64, gtequeue uint64) (QueryCursor, error)
}

// BackendSummoner defines an interface to store/retrieve Summoner data from Storage Backend
type BackendSummoner interface {
	GetSummonerByName(name string) (*Summoner, error)
	GetSummonerByNameTimeStamp(name string) time.Time

	GetSummonerBySummonerID(summonerID string) (*Summoner, error)
	GetSummonerBySummonerIDTimeStamp(summonerID string) time.Time

	GetSummonerByAccountID(accountID string) (*Summoner, error)
	GetSummonerByAccountIDTimeStamp(accountID string) time.Time

	GetSummonerByPUUID(PUUID string) (*Summoner, error)
	GetSummonerByPUUIDTimeStamp(PUUID string) time.Time

	StoreSummoner(data *Summoner) error
}

// BackendInternals defines an interface to retrieve internal infos from the Backend, e.g., number of stored elements
type BackendInternals interface {
	GetStorageSummary() (Summary, error)
}

// BackendStats defines an interface to retrieve stored statistics from Backend
type BackendStats interface {
	GetChampionStatsByChampionIDGameVersion(championID string, gameVersion string) (*ChampionStatsStorage, error)
	GetChampionStatsByChampionIDGameVersionTier(championID string, gameVersion string, tier string) (*ChampionStatsStorage, error)

	StoreChampionStats(stats *ChampionStatsStorage) error
}

// BackendMisc defines an interface to generic storages from Backend
type BackendMisc interface {
	GetKnownGameVersions() (*GameVersions, error)
	StoreKnownGameVersions(gameVersions *GameVersions) error
}

// Backend defines an interface to store/retrieve data from Storage Backend
type Backend interface {
	BackendChampion
	BackendFreeRotation
	BackendMatch
	BackendSummoner
	BackendLeague

	BackendInternals

	BackendStats

	BackendMisc
}
