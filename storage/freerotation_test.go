package storage

import (
	"testing"
	"time"

	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/riotclient"
)

func TestGettingFreeRotation(t *testing.T) {
	config := config.LoLStorage{}
	riotClient := &mockClient{}
	backend := &mockBackend{}

	config.MaxAgeChampion = 120
	config.MaxAgeChampionRotation = 120
	config.MaxAgeSummoner = 120

	storage, err := NewStorage(config, riotClient, backend)
	if err != nil || storage == nil {
		t.Fatalf("Could not get a new Storage: %s", err)
	}

	freeRotationBackend := riotclient.FreeRotation{}
	freeRotationBackend.FreeChampionIds = append(freeRotationBackend.FreeChampionIds, 345)
	freeRotationBackend.FreeChampionIds = append(freeRotationBackend.FreeChampionIds, 35)
	freeRotationBackend.MaxNewPlayerLevel = 5
	freeRotationBackend.FreeChampionIdsForNewPlayers = append(freeRotationBackend.FreeChampionIdsForNewPlayers, 7834)
	freeRotationBackend.FreeChampionIdsForNewPlayers = append(freeRotationBackend.FreeChampionIdsForNewPlayers, 9379493)
	freeRotationBackend.Timestamp = time.Now().Add(-time.Minute * 10)

	freeRotationClient := riotclient.FreeRotation{}
	freeRotationClient.FreeChampionIds = append(freeRotationClient.FreeChampionIds, 975)
	freeRotationClient.FreeChampionIds = append(freeRotationClient.FreeChampionIds, 3522)
	freeRotationClient.FreeChampionIds = append(freeRotationClient.FreeChampionIds, 3)
	freeRotationClient.MaxNewPlayerLevel = 10
	freeRotationClient.FreeChampionIdsForNewPlayers = append(freeRotationClient.FreeChampionIdsForNewPlayers, 7)
	freeRotationClient.FreeChampionIdsForNewPlayers = append(freeRotationClient.FreeChampionIdsForNewPlayers, 56)
	freeRotationClient.FreeChampionIdsForNewPlayers = append(freeRotationClient.FreeChampionIdsForNewPlayers, 123)
	freeRotationClient.Timestamp = time.Now()

	riotClient.setFreeRotation(freeRotationClient)

	// No stored list in backend should get it from client and store it
	actualFreeRotation := storage.getFreeRotation()

	if backend.getFreeRotationRetrieved() != false {
		t.Errorf("Storage got FreeRotation from backend even though it shouldn't")
	}
	if backend.getFreeRotationStored() != true {
		t.Errorf("Received FreeRotation where not stored in backend")
	}

	if riotClient.getFreeRotationRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	if len(freeRotationClient.FreeChampionIds) != len(actualFreeRotation.FreeChampionIds) {
		t.Error("Number of FreeChampionIds not as expected")
	}
	for i := 0; i < len(freeRotationClient.FreeChampionIds); i++ {
		if freeRotationClient.FreeChampionIds[i] != actualFreeRotation.FreeChampionIds[i] {
			t.Error("Champion does not match")
		}
	}
	if len(freeRotationClient.FreeChampionIdsForNewPlayers) != len(actualFreeRotation.FreeChampionIdsForNewPlayers) {
		t.Error("Number of FreeChampionIdsForNewPlayers not as expected")
	}
	for i := 0; i < len(freeRotationClient.FreeChampionIdsForNewPlayers); i++ {
		if freeRotationClient.FreeChampionIdsForNewPlayers[i] != actualFreeRotation.FreeChampionIdsForNewPlayers[i] {
			t.Error("Champion does not match")
		}
	}
	if freeRotationClient.MaxNewPlayerLevel != actualFreeRotation.MaxNewPlayerLevel {
		t.Error("MaxNewPlayerLevel does not match")
	}
	if freeRotationClient.Timestamp != actualFreeRotation.Timestamp {
		t.Error("Timestamp does not match")
	}

	// Stored list with valid TImeStamp in backend should get it from backend
	riotClient.reset()
	riotClient.setFreeRotation(freeRotationClient)
	backend.reset()
	backend.setFreeRotation(freeRotationBackend)

	actualFreeRotation = storage.getFreeRotation()

	if backend.getFreeRotationRetrieved() != true {
		t.Errorf("Storage got did not get FreeRotation from backend even though it should have")
	}
	if backend.getFreeRotationStored() != false {
		t.Errorf("Storage stored FreeRotation in backend even though it shouldn't")
	}

	if riotClient.getFreeRotationRetrieved() != false {
		t.Errorf("Storage did get the CHampions from client")
	}

	if len(freeRotationBackend.FreeChampionIds) != len(actualFreeRotation.FreeChampionIds) {
		t.Error("Number of FreeChampionIds not as expected")
	}
	for i := 0; i < len(freeRotationBackend.FreeChampionIds); i++ {
		if freeRotationBackend.FreeChampionIds[i] != actualFreeRotation.FreeChampionIds[i] {
			t.Errorf("Free Rotation Entry number %d does not match", i+1)
		}
	}
	if len(freeRotationBackend.FreeChampionIdsForNewPlayers) != len(actualFreeRotation.FreeChampionIdsForNewPlayers) {
		t.Error("Number of FreeChampionIdsForNewPlayers not as expected")
	}
	for i := 0; i < len(freeRotationBackend.FreeChampionIdsForNewPlayers); i++ {
		if freeRotationBackend.FreeChampionIdsForNewPlayers[i] != actualFreeRotation.FreeChampionIdsForNewPlayers[i] {
			t.Error("Champion does not match")
		}
	}
	if freeRotationBackend.MaxNewPlayerLevel != actualFreeRotation.MaxNewPlayerLevel {
		t.Error("MaxNewPlayerLevel does not match")
	}
	if freeRotationBackend.Timestamp != actualFreeRotation.Timestamp {
		t.Error("Timestamp does not match")
	}

	// Stored list with invalid TImeStamp in backend should get it from client
	riotClient.reset()
	riotClient.setFreeRotation(freeRotationClient)
	backend.reset()
	freeRotationBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeChampion+1))
	backend.setFreeRotation(freeRotationBackend)

	actualFreeRotation = storage.getFreeRotation()

	if backend.getFreeRotationRetrieved() != false {
		t.Errorf("Storage got FreeRotation from backend even though it shouldn't")
	}
	if backend.getFreeRotationStored() != true {
		t.Errorf("Received FreeRotation where not stored in backend")
	}

	if riotClient.getFreeRotationRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	if len(freeRotationClient.FreeChampionIds) != len(actualFreeRotation.FreeChampionIds) {
		t.Error("Number of FreeChampionIds not as expected")
	}
	for i := 0; i < len(freeRotationClient.FreeChampionIds); i++ {
		if freeRotationClient.FreeChampionIds[i] != actualFreeRotation.FreeChampionIds[i] {
			t.Errorf("Free Rotation Entry number %d does not match", i+1)
		}
	}
	if len(freeRotationClient.FreeChampionIdsForNewPlayers) != len(actualFreeRotation.FreeChampionIdsForNewPlayers) {
		t.Error("Number of FreeChampionIdsForNewPlayers not as expected")
	}
	for i := 0; i < len(freeRotationClient.FreeChampionIdsForNewPlayers); i++ {
		if freeRotationClient.FreeChampionIdsForNewPlayers[i] != actualFreeRotation.FreeChampionIdsForNewPlayers[i] {
			t.Error("Champion does not match")
		}
	}
	if freeRotationClient.MaxNewPlayerLevel != actualFreeRotation.MaxNewPlayerLevel {
		t.Error("MaxNewPlayerLevel does not match")
	}
	if freeRotationClient.Timestamp != actualFreeRotation.Timestamp {
		t.Error("Timestamp does not match")
	}

	// Stored list with invalid TImeStamp in backend should still get it from backend when client fails
	riotClient.reset()
	riotClient.setFreeRotation(freeRotationClient)
	backend.reset()
	freeRotationBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeChampion+1))
	backend.setFreeRotation(freeRotationBackend)
	riotClient.setFailFreeRotation(true)

	actualFreeRotation = storage.getFreeRotation()

	if backend.getFreeRotationRetrieved() != true {
		t.Errorf("Storage got did not get FreeRotation from backend even though it should have")
	}
	if backend.getFreeRotationStored() != false {
		t.Errorf("Storage stored FreeRotation in backend even though it shouldn't")
	}

	if riotClient.getFreeRotationRetrieved() != true {
		t.Errorf("Storage did not try to get the CHampions from client")
	}

	if len(freeRotationBackend.FreeChampionIds) != len(actualFreeRotation.FreeChampionIds) {
		t.Error("Number of FreeChampionIds not as expected")
	}
	for i := 0; i < len(freeRotationBackend.FreeChampionIds); i++ {
		if freeRotationBackend.FreeChampionIds[i] != actualFreeRotation.FreeChampionIds[i] {
			t.Errorf("Free Rotation Entry number %d does not match", i+1)
		}
	}
	if len(freeRotationBackend.FreeChampionIdsForNewPlayers) != len(actualFreeRotation.FreeChampionIdsForNewPlayers) {
		t.Error("Number of FreeChampionIdsForNewPlayers not as expected")
	}
	for i := 0; i < len(freeRotationBackend.FreeChampionIdsForNewPlayers); i++ {
		if freeRotationBackend.FreeChampionIdsForNewPlayers[i] != actualFreeRotation.FreeChampionIdsForNewPlayers[i] {
			t.Error("Champion does not match")
		}
	}
	if freeRotationBackend.MaxNewPlayerLevel != actualFreeRotation.MaxNewPlayerLevel {
		t.Error("MaxNewPlayerLevel does not match")
	}
	if freeRotationBackend.Timestamp != actualFreeRotation.Timestamp {
		t.Error("Timestamp does not match")
	}

	// Should always get it from Client when backend fails
	riotClient.reset()
	riotClient.setFreeRotation(freeRotationClient)
	backend.reset()
	freeRotationBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(5))
	backend.setFreeRotation(freeRotationBackend)
	backend.setFailFreeRotation(true)

	actualFreeRotation = storage.getFreeRotation()

	if backend.getFreeRotationRetrieved() != false {
		t.Errorf("Storage got FreeRotation from backend even though it shouldn't")
	}
	if backend.getFreeRotationStored() != true {
		t.Errorf("Received FreeRotation where not stored in backend")
	}

	if riotClient.getFreeRotationRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	if len(freeRotationClient.FreeChampionIds) != len(actualFreeRotation.FreeChampionIds) {
		t.Error("Number of FreeChampionIds not as expected")
	}
	for i := 0; i < len(freeRotationClient.FreeChampionIds); i++ {
		if freeRotationClient.FreeChampionIds[i] != actualFreeRotation.FreeChampionIds[i] {
			t.Errorf("Free Rotation Entry number %d does not match", i+1)
		}
	}
	if len(freeRotationClient.FreeChampionIdsForNewPlayers) != len(actualFreeRotation.FreeChampionIdsForNewPlayers) {
		t.Error("Number of FreeChampionIdsForNewPlayers not as expected")
	}
	for i := 0; i < len(freeRotationClient.FreeChampionIdsForNewPlayers); i++ {
		if freeRotationClient.FreeChampionIdsForNewPlayers[i] != actualFreeRotation.FreeChampionIdsForNewPlayers[i] {
			t.Error("Champion does not match")
		}
	}
	if freeRotationClient.MaxNewPlayerLevel != actualFreeRotation.MaxNewPlayerLevel {
		t.Error("MaxNewPlayerLevel does not match")
	}
	if freeRotationClient.Timestamp != actualFreeRotation.Timestamp {
		t.Error("Timestamp does not match")
	}

	// Should return empty list if everything fails
	riotClient.reset()
	riotClient.setFreeRotation(freeRotationClient)
	backend.reset()
	freeRotationBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(5))
	backend.setFreeRotation(freeRotationBackend)
	backend.setFailFreeRotation(true)
	riotClient.setFailFreeRotation(true)

	actualFreeRotation = storage.getFreeRotation()

	if backend.getFreeRotationRetrieved() != true {
		t.Errorf("Storage should have tried to get FreeRotation from backend")
	}
	if backend.getFreeRotationStored() != false {
		t.Errorf("Even though client failed it stored data in Backend")
	}

	if riotClient.getFreeRotationRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	empty := riotclient.FreeRotation{}

	if len(empty.FreeChampionIds) != len(actualFreeRotation.FreeChampionIds) {
		t.Error("Number of FreeRotation not as expected")
	}
	if len(empty.FreeChampionIdsForNewPlayers) != len(actualFreeRotation.FreeChampionIdsForNewPlayers) {
		t.Error("Number of FreeRotation not as expected")
	}
	if empty.MaxNewPlayerLevel != actualFreeRotation.MaxNewPlayerLevel {
		t.Error("MaxNewPlayerLevel of FreeRotation not as expected")
	}
	if empty.Timestamp != actualFreeRotation.Timestamp {
		t.Error("Timestamp does not match")
	}
}
