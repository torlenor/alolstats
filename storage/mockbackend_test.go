package storage

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

type mockBackend struct {
	failChampions        bool
	champions            riotclient.ChampionsList
	championsTimeStamp   time.Time
	championsRetrieved   bool
	championsStored      riotclient.ChampionsList
	championsWhereStored bool

	failFreeRotation        bool
	freeRotation            riotclient.FreeRotation
	freeRotationRetrieved   bool
	freeRotationStored      riotclient.FreeRotation
	freeRotationWhereStored bool

	failSummoner         bool
	summoner             Summoner
	wasSummonerRetrieved bool
	storedSummoner       Summoner
	wasSummonerStored    bool
}

func (b *mockBackend) reset() {
	b.failChampions = false
	b.champions = make(riotclient.ChampionsList)
	b.championsTimeStamp = time.Time{}
	b.championsRetrieved = false
	b.championsStored = make(riotclient.ChampionsList)
	b.championsWhereStored = false

	b.failFreeRotation = false
	b.freeRotation = riotclient.FreeRotation{}
	b.freeRotationRetrieved = false
	b.freeRotationStored = riotclient.FreeRotation{}
	b.freeRotationWhereStored = false

	b.failSummoner = false
	b.summoner = Summoner{}
	b.wasSummonerRetrieved = false
	b.storedSummoner = Summoner{}
	b.wasSummonerStored = false
}

//
// Champions
//

func (b *mockBackend) setChampions(champions riotclient.ChampionsList) {
	b.champions = champions
}
func (b *mockBackend) setChampionsTimeStamp(time time.Time) {
	b.championsTimeStamp = time
}

func (b *mockBackend) setFailChampions(fail bool) {
	b.failChampions = fail
}

func (b *mockBackend) getChampionsRetrieved() bool {
	return b.championsRetrieved
}

func (b *mockBackend) getChampionsStored() bool {
	return b.championsWhereStored
}

func (b *mockBackend) GetChampions() (riotclient.ChampionsList, error) {
	b.championsRetrieved = true

	if b.failChampions {
		return nil, fmt.Errorf("Error retreiving champions")
	}

	return b.champions, nil
}

func (b *mockBackend) GetChampionsTimeStamp() time.Time {
	if b.failChampions {
		t1, _ := time.Parse(
			time.RFC3339,
			"2010-11-01T00:08:41+00:00")
		return t1
	}
	return b.championsTimeStamp
}

func (b *mockBackend) StoreChampions(championList riotclient.ChampionsList) error {
	b.championsStored = championList
	b.championsWhereStored = true
	return nil
}

//
// Free Rotation
//

func (b *mockBackend) setFreeRotation(freeRotation riotclient.FreeRotation) {
	b.freeRotation = freeRotation
}

func (b *mockBackend) setFailFreeRotation(fail bool) {
	b.failFreeRotation = fail
}

func (b *mockBackend) getFreeRotationRetrieved() bool {
	return b.freeRotationRetrieved
}

func (b *mockBackend) getFreeRotationStored() bool {
	return b.freeRotationWhereStored
}

func (b *mockBackend) GetFreeRotation() (*riotclient.FreeRotation, error) {
	b.freeRotationRetrieved = true

	if b.failFreeRotation {
		return nil, fmt.Errorf("Error retreiving Free Rotation")
	}

	return &b.freeRotation, nil
}

func (b *mockBackend) GetFreeRotationTimeStamp() time.Time {
	if b.failFreeRotation {
		t1, _ := time.Parse(
			time.RFC3339,
			"2010-11-01T00:08:41+00:00")
		return t1
	}
	return b.freeRotation.Timestamp
}

func (b *mockBackend) StoreFreeRotation(freeRotation *riotclient.FreeRotation) error {
	b.freeRotationStored = *freeRotation
	b.freeRotationWhereStored = true
	return nil
}

//
// Summoner
//

func (b *mockBackend) setSummoner(data Summoner) {
	b.summoner = data
}

func (b *mockBackend) setFailSummoner(fail bool) {
	b.failSummoner = fail
}

func (b *mockBackend) getWasSummonerRetrieved() bool {
	return b.wasSummonerRetrieved
}

func (b *mockBackend) getWasSummonerStored() bool {
	return b.wasSummonerStored
}

func (b *mockBackend) GetSummonerByName(name string) (*Summoner, error) {
	b.wasSummonerRetrieved = true

	if b.failSummoner {
		return nil, fmt.Errorf("Error retreiving Summoner")
	}

	if name == b.summoner.SummonerName {
		return &b.summoner, nil
	}

	return nil, fmt.Errorf("Summoner not found")
}

func (b *mockBackend) GetSummonerByNameTimeStamp(name string) time.Time {
	return b.summoner.SummonerDTO.Timestamp
}

func (b *mockBackend) GetSummonerBySummonerID(summonerID string) (*Summoner, error) {
	b.wasSummonerRetrieved = true

	if b.failSummoner {
		return nil, fmt.Errorf("Error retreiving Summoner")
	}

	if summonerID == b.summoner.SummonerDTO.ID {
		return &b.summoner, nil
	}

	return nil, fmt.Errorf("Summoner not found")
}

func (b *mockBackend) GetSummonerBySummonerIDTimeStamp(summonerID string) time.Time {
	return b.summoner.SummonerDTO.Timestamp
}

func (b *mockBackend) GetSummonerByAccountID(accountID string) (*Summoner, error) {
	b.wasSummonerRetrieved = true

	if b.failSummoner {
		return nil, fmt.Errorf("Error retreiving Summoner")
	}

	if accountID == b.summoner.SummonerDTO.AccountID {
		return &b.summoner, nil
	}

	return nil, fmt.Errorf("Summoner not found")
}

func (b *mockBackend) GetSummonerByAccountIDTimeStamp(accountID string) time.Time {
	return b.summoner.SummonerDTO.Timestamp
}

func (b *mockBackend) StoreSummoner(data *Summoner) error {
	b.summoner = *data
	b.wasSummonerStored = true
	return nil
}

//
// Match
//

func (b *mockBackend) GetMatch(id uint64) (*riotclient.MatchDTO, error) {
	return &riotclient.MatchDTO{}, nil
}

func (b *mockBackend) StoreMatch(data *riotclient.MatchDTO) error {
	return nil
}

func (b *mockBackend) GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (*riotclient.Matches, error) {
	return &riotclient.Matches{}, nil
}

func (b *mockBackend) GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queue uint64) (*riotclient.Matches, error) {
	return &riotclient.Matches{}, nil
}

func (b *mockBackend) GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (*riotclient.Matches, error) {
	return &riotclient.Matches{}, nil
}

func (b *mockBackend) GetStorageSummary() (Summary, error) {
	return Summary{}, nil
}

func (b *mockBackend) GetLeagueByQueue(league string, queue string) (*riotclient.LeagueListDTO, error) {
	return &riotclient.LeagueListDTO{}, nil
}

func (b *mockBackend) GetLeagueByQueueTimeStamp(league string, queue string) (time.Time, error) {
	return time.Time{}, nil
}

func (b *mockBackend) StoreLeague(*riotclient.LeagueListDTO) error {
	return nil
}

func (b *mockBackend) GetLeaguesForSummoner(summonerName string) (*SummonerLeagues, error) {
	return &SummonerLeagues{}, nil
}

func (b *mockBackend) GetLeaguesForSummonerBySummonerID(summonerID string) (*SummonerLeagues, error) {
	return &SummonerLeagues{}, nil
}

func (b *mockBackend) StoreLeaguesForSummoner(*SummonerLeagues) error {
	return nil
}

func (b *mockBackend) GetLeaguesForSummonerTimeStamp(summonerName string) (time.Time, error) {
	return time.Time{}, nil
}

func (b *mockBackend) GetLeaguesForSummonerBySummonerIDTimeStamp(summonerID string) (time.Time, error) {
	return time.Time{}, nil
}

func (b *mockBackend) GetMatchTimeLine(matchID uint64) (*riotclient.MatchTimelineDTO, error) {
	return &riotclient.MatchTimelineDTO{}, nil
}

func (b *mockBackend) StoreMatchTimeLine(data *riotclient.MatchTimelineDTO) error {
	return nil
}

func (b *mockBackend) GetSummonerByPUUID(PUUID string) (*Summoner, error) {
	return &Summoner{}, nil
}

func (b *mockBackend) GetSummonerByPUUIDTimeStamp(PUUID string) time.Time {
	return time.Time{}
}

func (b *mockBackend) StoreChampionStats(data *ChampionStatsStorage) error {
	return fmt.Errorf("Not implemented")
}

func (b *mockBackend) GetChampionStatsByChampionIDGameVersionTierQueue(championID string, gameVersion string, tier string, queue string) (*ChampionStatsStorage, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *mockBackend) GetKnownGameVersions() (*GameVersions, error) {
	return &GameVersions{}, nil
}

func (b *mockBackend) StoreKnownGameVersions(gameVersions *GameVersions) error {
	return fmt.Errorf("Not implemented")
}

func (b *mockBackend) GetMatchesCursorByGameVersion(gameVersion string) (QueryCursor, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *mockBackend) GetMatchesCursorByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (QueryCursor, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *mockBackend) GetMatchesCursorByGameVersionMapBetweenQueueIDs(gameVersion string, mapID uint64, ltequeue uint64, gtequeue uint64) (QueryCursor, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *mockBackend) StoreChampionStatsSummary(statsSummary *ChampionStatsSummaryStorage) error {
	return fmt.Errorf("Not implemented")
}

func (b *mockBackend) GetChampionStatsSummaryByGameVersionTier(gameVersion string, tier string) (*ChampionStatsSummaryStorage, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *mockBackend) GetItemStatsByChampionIDGameVersion(championID, gameVersion string) (*ItemStatsStorage, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *mockBackend) StoreItemStats(data *ItemStatsStorage) error {
	return fmt.Errorf("Not implemented")
}
