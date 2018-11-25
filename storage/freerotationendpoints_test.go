package storage

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/riotclient"
)

func TestFreeRotationEndpoint(t *testing.T) {
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
	backend.setFreeRotation(freeRotationBackend)
	riotClient.setFailFreeRotation(true)

	req := httptest.NewRequest("GET", "http://example.com/endpoint", nil)
	w := httptest.NewRecorder()
	storage.freeRotationEndpoint(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	receivedFreeRotation := riotclient.FreeRotation{}
	err = json.Unmarshal(body, &receivedFreeRotation)
	if err != nil {
		t.Fatalf("Decoding the json into the correct data struct was not possible: %s", err)
	}

	if !cmp.Equal(freeRotationBackend, receivedFreeRotation) {
		t.Error("Data not equal")
	}

	// Could not get any data and return empty data
	backend.setFailFreeRotation(true)

	req = httptest.NewRequest("GET", "http://example.com/endpoint", nil)
	w = httptest.NewRecorder()
	storage.freeRotationEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	receivedFreeRotation = riotclient.FreeRotation{}
	err = json.Unmarshal(body, &receivedFreeRotation)
	if err != nil {
		t.Fatalf("Decoding the json into the correct data struct was not possible: %s", err)
	}

	if !cmp.Equal(riotclient.FreeRotation{}, receivedFreeRotation) {
		t.Error("Data not equal")
	}
}
