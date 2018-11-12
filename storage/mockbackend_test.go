package storage

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

type mockBackend struct {
	failChampions        bool
	champions            riotclient.ChampionList
	championsRetrieved   bool
	championsStored      riotclient.ChampionList
	championsWhereStored bool
}

func (b *mockBackend) reset() {
	b.failChampions = false
	b.champions = riotclient.ChampionList{}
	b.championsRetrieved = false
	b.championsStored = riotclient.ChampionList{}
	b.championsWhereStored = false
}

func (b *mockBackend) setChampions(champions riotclient.ChampionList) {
	b.champions = champions
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
	return b.champions.Timestamp
}

func (b *mockBackend) StoreChampions(championList riotclient.ChampionList) error {
	b.championsStored = championList
	b.championsWhereStored = true
	return nil
}

func (b *mockBackend) GetFreeRotation() (riotclient.FreeRotation, error) {
	return riotclient.FreeRotation{}, nil
}

func (b *mockBackend) GetFreeRotationTimeStamp() time.Time {
	return time.Now()
}

func (b *mockBackend) StoreFreeRotation(freeRotation riotclient.FreeRotation) error {
	return nil
}

func (b *mockBackend) GetMatch(id uint64) (riotclient.Match, error) {
	return riotclient.Match{}, nil
}

func (b *mockBackend) StoreMatch(data *riotclient.Match) error {
	return nil
}
