package memorybackend

import (
	"testing"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

func TestStoreAndRetreiveFreeRotation(t *testing.T) {
	backend, err := NewBackend()
	if err != nil || backend == nil {
		t.Fatalf("Could not get a new Memory Backend: %s", err)
	}
	freeRotation := riotclient.FreeRotation{}
	freeRotation.FreeChampionIds = append(freeRotation.FreeChampionIds, 345)
	freeRotation.FreeChampionIds = append(freeRotation.FreeChampionIds, 35)

	freeRotation.MaxNewPlayerLevel = 5
	freeRotation.FreeChampionIdsForNewPlayers = append(freeRotation.FreeChampionIdsForNewPlayers, 7834)
	freeRotation.FreeChampionIdsForNewPlayers = append(freeRotation.FreeChampionIdsForNewPlayers, 9379493)

	t1, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	freeRotation.Timestamp = t1

	backend.StoreFreeRotation(freeRotation)
	if err != nil {
		t.Fatalf("Failed storing FreeRotation: %s", err)
	}

	if backend.GetFreeRotationTimeStamp() != t1 {
		t.Error("GetFreeRotationTimeStamp() TimeStamp does not match")
	}

	retreivedFreeRotation, err := backend.GetFreeRotation()
	if err != nil {
		t.Fatalf("Could not retrieve FreeRotation: %s", err)
	}

	if len(retreivedFreeRotation.FreeChampionIds) != len(freeRotation.FreeChampionIds) {
		t.Error("Number of FreeChampionIds not as expected")
	}
	if retreivedFreeRotation.FreeChampionIds[0] != freeRotation.FreeChampionIds[0] {
		t.Error("FreeChampionIds does not match")
	}
	if retreivedFreeRotation.FreeChampionIds[1] != freeRotation.FreeChampionIds[1] {
		t.Error("FreeChampionIds does not match")
	}

	if retreivedFreeRotation.MaxNewPlayerLevel != freeRotation.MaxNewPlayerLevel {
		t.Error("MaxNewPlayerLevel does not match")
	}
	if len(retreivedFreeRotation.FreeChampionIdsForNewPlayers) != len(freeRotation.FreeChampionIdsForNewPlayers) {
		t.Error("Number of FreeChampionIdsForNewPlayers not as expected")
	}
	if retreivedFreeRotation.FreeChampionIdsForNewPlayers[0] != freeRotation.FreeChampionIdsForNewPlayers[0] {
		t.Error("FreeChampionIdsForNewPlayers does not match")
	}
	if retreivedFreeRotation.FreeChampionIdsForNewPlayers[1] != freeRotation.FreeChampionIdsForNewPlayers[1] {
		t.Error("FreeChampionIdsForNewPlayers does not match")
	}

	if retreivedFreeRotation.Timestamp != t1 {
		t.Error("Timestamp does not match")
	}
}
