package storage

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

type mockBackend struct {
	failChampions        bool
	champions            riotclient.ChampionList
	championsTimeStamp   time.Time
	championsRetrieved   bool
	championsStored      riotclient.ChampionList
	championsWhereStored bool

	failFreeRotation        bool
	freeRotation            riotclient.FreeRotation
	freeRotationRetrieved   bool
	freeRotationStored      riotclient.FreeRotation
	freeRotationWhereStored bool

	failSummoner         bool
	summoner             riotclient.Summoner
	wasSummonerRetrieved bool
	storedSummoner       riotclient.Summoner
	wasSummonerStored    bool
}

func (b *mockBackend) reset() {
	b.failChampions = false
	b.champions = riotclient.ChampionList{}
	b.championsTimeStamp = time.Time{}
	b.championsRetrieved = false
	b.championsStored = riotclient.ChampionList{}
	b.championsWhereStored = false

	b.failFreeRotation = false
	b.freeRotation = riotclient.FreeRotation{}
	b.freeRotationRetrieved = false
	b.freeRotationStored = riotclient.FreeRotation{}
	b.freeRotationWhereStored = false

	b.failSummoner = false
	b.summoner = riotclient.Summoner{}
	b.wasSummonerRetrieved = false
	b.storedSummoner = riotclient.Summoner{}
	b.wasSummonerStored = false
}

//
// Champions
//

func (b *mockBackend) setChampions(champions riotclient.ChampionList) {
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

func (b *mockBackend) GetChampions() (riotclient.ChampionList, error) {
	b.championsRetrieved = true

	if b.failChampions {
		return riotclient.ChampionList{}, fmt.Errorf("Error retreiving champions")
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

func (b *mockBackend) StoreChampions(championList riotclient.ChampionList) error {
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

func (b *mockBackend) GetFreeRotation() (riotclient.FreeRotation, error) {
	b.freeRotationRetrieved = true

	if b.failFreeRotation {
		return riotclient.FreeRotation{}, fmt.Errorf("Error retreiving Free Rotation")
	}

	return b.freeRotation, nil
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

func (b *mockBackend) StoreFreeRotation(freeRotation riotclient.FreeRotation) error {
	b.freeRotationStored = freeRotation
	b.freeRotationWhereStored = true
	return nil
}

//
// Summoner
//

func (b *mockBackend) setSummoner(data riotclient.Summoner) {
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

func (b *mockBackend) GetSummonerByName(name string) (riotclient.Summoner, error) {
	b.wasSummonerRetrieved = true

	if b.failSummoner {
		return riotclient.Summoner{}, fmt.Errorf("Error retreiving Summoner")
	}

	if name == b.summoner.Name {
		return b.summoner, nil
	}

	return riotclient.Summoner{}, fmt.Errorf("Summoner not found")
}

func (b *mockBackend) GetSummonerByNameTimeStamp(name string) time.Time {
	return b.summoner.Timestamp
}

func (b *mockBackend) GetSummonerBySummonerID(summonerID uint64) (riotclient.Summoner, error) {
	b.wasSummonerRetrieved = true

	if b.failSummoner {
		return riotclient.Summoner{}, fmt.Errorf("Error retreiving Summoner")
	}

	if summonerID == b.summoner.ID {
		return b.summoner, nil
	}

	return riotclient.Summoner{}, fmt.Errorf("Summoner not found")
}

func (b *mockBackend) GetSummonerBySummonerIDTimeStamp(summonerID uint64) time.Time {
	return b.summoner.Timestamp
}

func (b *mockBackend) GetSummonerByAccountID(accountID uint64) (riotclient.Summoner, error) {
	b.wasSummonerRetrieved = true

	if b.failSummoner {
		return riotclient.Summoner{}, fmt.Errorf("Error retreiving Summoner")
	}

	if accountID == b.summoner.AccountID {
		return b.summoner, nil
	}

	return riotclient.Summoner{}, fmt.Errorf("Summoner not found")
}

func (b *mockBackend) GetSummonerByAccountIDTimeStamp(accountID uint64) time.Time {
	return b.summoner.Timestamp
}

func (b *mockBackend) StoreSummoner(data *riotclient.Summoner) error {
	b.summoner = *data
	b.wasSummonerStored = true
	return nil
}

//
// Match
//

func (b *mockBackend) GetMatch(id uint64) (riotclient.Match, error) {
	return riotclient.Match{}, nil
}

func (b *mockBackend) StoreMatch(data *riotclient.Match) error {
	return nil
}

func (b *mockBackend) GetMatchesByGameVersion(gameVersion string) (riotclient.Matches, error) {
	return riotclient.Matches{}, nil
}

func (b *mockBackend) GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (riotclient.Matches, error) {
	return riotclient.Matches{}, nil
}

func (b *mockBackend) GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queue uint64) (riotclient.Matches, error) {
	return riotclient.Matches{}, nil
}

func (b *mockBackend) GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (riotclient.Matches, error) {
	return riotclient.Matches{}, nil
}

func (b *mockBackend) GetStorageSummary() (Summary, error) {
	return Summary{}, nil
}
