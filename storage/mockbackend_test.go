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

	failFreeRotation        bool
	freeRotation            riotclient.FreeRotation
	freeRotationRetrieved   bool
	freeRotationStored      riotclient.FreeRotation
	freeRotationWhereStored bool
}

func (b *mockBackend) reset() {
	b.failChampions = false
	b.champions = riotclient.ChampionList{}
	b.championsRetrieved = false
	b.championsStored = riotclient.ChampionList{}
	b.championsWhereStored = false

	b.failFreeRotation = false
	b.freeRotation = riotclient.FreeRotation{}
	b.freeRotationRetrieved = false
	b.freeRotationStored = riotclient.FreeRotation{}
	b.freeRotationWhereStored = false
}

//
// Champions
//

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
// Match
//

func (b *mockBackend) GetMatch(id uint64) (riotclient.Match, error) {
	return riotclient.Match{}, nil
}

func (b *mockBackend) StoreMatch(data *riotclient.Match) error {
	return nil
}
