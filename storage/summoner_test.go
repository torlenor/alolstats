package storage

import (
	"testing"
	"time"

	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/riotclient"
)

func TestGettingSummonerByName(t *testing.T) {
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

	summonerBackend := riotclient.Summoner{
		AccountID:    112345,
		ID:           212345,
		Name:         "Backend Summoner",
		Level:        10,
		RevisionDate: 312345,
		Timestamp:    time.Now().Add(-time.Minute * 10),
	}

	summonerClient := riotclient.Summoner{
		AccountID:    412345,
		ID:           512345,
		Name:         "Client Summoner",
		Level:        20,
		RevisionDate: 612345,
		Timestamp:    time.Now(),
	}

	riotClient.setSummoner(summonerClient)

	// No stored list in backend should get it from client and store it
	actualSummoner, err := storage.GetSummonerByName("Test")
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data was not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with valid TImeStamp in backend should get it from backend
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now()
	backend.setSummoner(summonerBackend)

	actualSummoner, err = storage.GetSummonerByName("Test")
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get it from backend even though it should have")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Storage stored data in backend even though it shouldn't")
	}

	if riotClient.getWasSummonerRetrieved() != false {
		t.Errorf("Storage did get the data from client")
	}

	if summonerBackend != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with invalid TImeStamp in backend should get it from client
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner+1))
	backend.setSummoner(summonerBackend)

	actualSummoner, err = storage.GetSummonerByName("Test")
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data where not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with invalid TImeStamp in backend should still get it from backend when client fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner+1))
	backend.setSummoner(summonerBackend)
	riotClient.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerByName("Test")
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage got did not get it from backend even though it should have")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Storage stored data in backend even though it shouldn't")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not try to get the data from client")
	}

	if summonerBackend != actualSummoner {
		t.Error("Data does not match")
	}

	// Should always get it from Client when backend fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner-1))
	backend.setSummoner(summonerBackend)
	backend.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerByName("Test")
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data where not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Should return empty list if everything fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner-1))
	backend.setSummoner(summonerBackend)
	backend.setFailSummoner(true)
	riotClient.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerByName("Test")
	if err == nil {
		t.Errorf("There should have been an error, but wasn't")
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage should have tried to get data from backend")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Even though client failed it stored data in Backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	empty := riotclient.Summoner{}

	if empty != actualSummoner {
		t.Error("Data not equal empty data")
	}
}

func TestGettingSummonerBySummonerID(t *testing.T) {
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

	summonerBackend := riotclient.Summoner{
		AccountID:    112345,
		ID:           212345,
		Name:         "Backend Summoner",
		Level:        10,
		RevisionDate: 312345,
		Timestamp:    time.Now().Add(-time.Minute * 10),
	}

	summonerClient := riotclient.Summoner{
		AccountID:    412345,
		ID:           512345,
		Name:         "Client Summoner",
		Level:        20,
		RevisionDate: 612345,
		Timestamp:    time.Now(),
	}

	riotClient.setSummoner(summonerClient)

	// No stored list in backend should get it from client and store it
	actualSummoner, err := storage.GetSummonerBySummonerID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data was not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with valid TImeStamp in backend should get it from backend
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now()
	backend.setSummoner(summonerBackend)

	actualSummoner, err = storage.GetSummonerBySummonerID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get it from backend even though it should have")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Storage stored data in backend even though it shouldn't")
	}

	if riotClient.getWasSummonerRetrieved() != false {
		t.Errorf("Storage did get the data from client")
	}

	if summonerBackend != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with invalid TImeStamp in backend should get it from client
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner+1))
	backend.setSummoner(summonerBackend)

	actualSummoner, err = storage.GetSummonerBySummonerID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data where not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with invalid TImeStamp in backend should still get it from backend when client fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner+1))
	backend.setSummoner(summonerBackend)
	riotClient.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerBySummonerID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage got did not get it from backend even though it should have")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Storage stored data in backend even though it shouldn't")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not try to get the data from client")
	}

	if summonerBackend != actualSummoner {
		t.Error("Data does not match")
	}

	// Should always get it from Client when backend fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner-1))
	backend.setSummoner(summonerBackend)
	backend.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerBySummonerID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data where not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Should return empty list if everything fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner-1))
	backend.setSummoner(summonerBackend)
	backend.setFailSummoner(true)
	riotClient.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerBySummonerID(1234)
	if err == nil {
		t.Errorf("There should have been an error, but wasn't")
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage should have tried to get data from backend")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Even though client failed it stored data in Backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	empty := riotclient.Summoner{}

	if empty != actualSummoner {
		t.Error("Data not equal empty data")
	}
}

func TestGettingSummonerByAccountID(t *testing.T) {
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

	summonerBackend := riotclient.Summoner{
		AccountID:    112345,
		ID:           212345,
		Name:         "Backend Summoner",
		Level:        10,
		RevisionDate: 312345,
		Timestamp:    time.Now().Add(-time.Minute * 10),
	}

	summonerClient := riotclient.Summoner{
		AccountID:    412345,
		ID:           512345,
		Name:         "Client Summoner",
		Level:        20,
		RevisionDate: 612345,
		Timestamp:    time.Now(),
	}

	riotClient.setSummoner(summonerClient)

	// No stored list in backend should get it from client and store it
	actualSummoner, err := storage.GetSummonerByAccountID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data was not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with valid TImeStamp in backend should get it from backend
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now()
	backend.setSummoner(summonerBackend)

	actualSummoner, err = storage.GetSummonerByAccountID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get it from backend even though it should have")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Storage stored data in backend even though it shouldn't")
	}

	if riotClient.getWasSummonerRetrieved() != false {
		t.Errorf("Storage did get the data from client")
	}

	if summonerBackend != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with invalid TImeStamp in backend should get it from client
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner+1))
	backend.setSummoner(summonerBackend)

	actualSummoner, err = storage.GetSummonerByAccountID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data where not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Stored list with invalid TImeStamp in backend should still get it from backend when client fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner+1))
	backend.setSummoner(summonerBackend)
	riotClient.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerByAccountID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage got did not get it from backend even though it should have")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Storage stored data in backend even though it shouldn't")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not try to get the data from client")
	}

	if summonerBackend != actualSummoner {
		t.Error("Data does not match")
	}

	// Should always get it from Client when backend fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner-1))
	backend.setSummoner(summonerBackend)
	backend.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerByAccountID(1234)
	if err != nil {
		t.Errorf("There was an error: %s", err)
	}

	if backend.getWasSummonerRetrieved() != false {
		t.Errorf("Storage got it from backend even though it shouldn't")
	}
	if backend.getWasSummonerStored() != true {
		t.Errorf("Received data where not stored in backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	if summonerClient != actualSummoner {
		t.Error("Data does not match")
	}

	// Should return empty list if everything fails
	riotClient.reset()
	riotClient.setSummoner(summonerClient)
	backend.reset()
	summonerBackend.Timestamp = time.Now().Add(-time.Minute * time.Duration(config.MaxAgeSummoner-1))
	backend.setSummoner(summonerBackend)
	backend.setFailSummoner(true)
	riotClient.setFailSummoner(true)

	actualSummoner, err = storage.GetSummonerByAccountID(1234)
	if err == nil {
		t.Errorf("There should have been an error, but wasn't")
	}

	if backend.getWasSummonerRetrieved() != true {
		t.Errorf("Storage should have tried to get data from backend")
	}
	if backend.getWasSummonerStored() != false {
		t.Errorf("Even though client failed it stored data in Backend")
	}

	if riotClient.getWasSummonerRetrieved() != true {
		t.Errorf("Storage did not get the data from client")
	}

	empty := riotclient.Summoner{}

	if empty != actualSummoner {
		t.Error("Data not equal empty data")
	}
}
