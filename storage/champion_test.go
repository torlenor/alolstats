package storage

import (
	"testing"
	"time"

	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/riotclient"
)

func TestGettingChampionList(t *testing.T) {
	config := config.LoLStorage{}
	riotClient := &mockClient{}
	backend := &mockBackend{}

	config.MaxAgeChampion = 120
	config.MaxAgeChampionRotation = 120
	config.MaxAgeSummoner = 120
	config.DefaultRiotClient = "euw1"

	storage, err := NewStorage(config, map[string]riotclient.Client{"euw1": riotClient}, backend)
	if err != nil || storage == nil {
		t.Fatalf("Could not get a new Storage: %s", err)
	}

	championsListBackend := make(riotclient.ChampionsList)
	championsListBackend["3"] = riotclient.Champion{Name: "BACKEND CHAMP",
		ID: "3", Timestamp: time.Now().Add(-time.Minute * 10)}
	championsListBackend["534"] = riotclient.Champion{Name: "BACKEND CHAMP 2",
		ID: "534", Timestamp: time.Now().Add(-time.Minute * 10)}

	championsListClient := make(riotclient.ChampionsList)
	championsListClient["43222"] = riotclient.Champion{Name: "CLIENT CHAMP",
		ID: "43222", Timestamp: time.Now()}
	championsListClient["123"] = riotclient.Champion{Name: "CLIENT CHAMP 2",
		ID: "123", Timestamp: time.Now()}

	riotClient.setChampions(championsListClient)

	// No stored list in backend should get it from client and store it
	actualChampions := storage.GetChampions(false)

	if backend.getChampionsRetrieved() != false {
		t.Errorf("Storage got Champions from backend even though it shouldn't")
	}
	if backend.getChampionsStored() != true {
		t.Errorf("Received Champions where not stored in backend")
	}

	if riotClient.getChampionsRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	if len(championsListClient) != len(actualChampions) {
		t.Error("Number of Champions not as expected")
	}
	if championsListClient["43222"].Name != actualChampions["43222"].Name {
		t.Error("Champion does not match")
	}
	if championsListClient["43222"].ID != actualChampions["43222"].ID {
		t.Error("Champion does not match")
	}
	if championsListClient["43222"].Timestamp != actualChampions["43222"].Timestamp {
		t.Error("Timestamp does not match")
	}
	if championsListClient["123"].Name != actualChampions["123"].Name {
		t.Error("Champion does not match")
	}
	if championsListClient["123"].ID != actualChampions["123"].ID {
		t.Error("Champion does not match")
	}
	if championsListClient["123"].Timestamp != actualChampions["123"].Timestamp {
		t.Error("Timestamp does not match")
	}

	// Stored list with valid TImeStamp in backend should get it from backend
	riotClient.reset()
	riotClient.setChampions(championsListClient)
	backend.reset()
	backend.setChampions(championsListBackend)
	backend.setChampionsTimeStamp(time.Now())

	actualChampions = storage.GetChampions(false)

	if backend.getChampionsRetrieved() != true {
		t.Errorf("Storage got did not get Champions from backend even though it should have")
	}
	if backend.getChampionsStored() != false {
		t.Errorf("Storage stored Champions in backend even though it shouldn't")
	}

	if riotClient.getChampionsRetrieved() != false {
		t.Errorf("Storage did get the CHampions from client")
	}

	if len(championsListBackend) != len(actualChampions) {
		t.Error("Number of Champions not as expected")
	}
	if championsListBackend["3"].Name != actualChampions["3"].Name {
		t.Error("Champion does not match")
	}
	if championsListBackend["3"].ID != actualChampions["3"].ID {
		t.Error("Champion does not match")
	}
	if championsListBackend["3"].Timestamp != actualChampions["3"].Timestamp {
		t.Error("Timestamp does not match")
	}
	if championsListBackend["534"].Name != actualChampions["534"].Name {
		t.Error("Champion does not match")
	}
	if championsListBackend["534"].ID != actualChampions["534"].ID {
		t.Error("Champion does not match")
	}
	if championsListBackend["534"].Timestamp != actualChampions["534"].Timestamp {
		t.Error("Timestamp does not match")
	}

	// Stored list with invalid TImeStamp in backend should get it from client
	riotClient.reset()
	riotClient.setChampions(championsListClient)
	backend.reset()
	championsListBackend["534"] = riotclient.Champion{Name: "BACKEND CHAMP 2",
		ID: "534", Timestamp: time.Now().Add(-time.Minute * time.Duration(config.MaxAgeChampion+1))}
	backend.setChampions(championsListBackend)
	backend.setChampionsTimeStamp(time.Now().Add(-time.Minute * time.Duration(config.MaxAgeChampion+1)))

	actualChampions = storage.GetChampions(false)

	if backend.getChampionsRetrieved() != false {
		t.Errorf("Storage got Champions from backend even though it shouldn't")
	}
	if backend.getChampionsStored() != true {
		t.Errorf("Received Champions where not stored in backend")
	}

	if riotClient.getChampionsRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	if len(championsListClient) != len(actualChampions) {
		t.Error("Number of Champions not as expected")
	}
	if championsListClient["43222"].Name != actualChampions["43222"].Name {
		t.Error("Champion does not match")
	}
	if championsListClient["43222"].ID != actualChampions["43222"].ID {
		t.Error("Champion does not match")
	}
	if championsListClient["43222"].Timestamp != actualChampions["43222"].Timestamp {
		t.Error("Timestamp does not match")
	}
	if championsListClient["123"].Name != actualChampions["123"].Name {
		t.Error("Champion does not match")
	}
	if championsListClient["123"].ID != actualChampions["123"].ID {
		t.Error("Champion does not match")
	}
	if championsListClient["123"].Timestamp != actualChampions["123"].Timestamp {
		t.Error("Timestamp does not match")
	}

	// Stored list with invalid TImeStamp in backend should still get it from backend when client fails
	riotClient.reset()
	riotClient.setChampions(championsListClient)
	backend.reset()
	backend.setChampions(championsListBackend)
	backend.setChampionsTimeStamp(time.Now().Add(-time.Minute * time.Duration(config.MaxAgeChampion+1)))
	riotClient.setFailChampions(true)

	actualChampions = storage.GetChampions(false)

	if backend.getChampionsRetrieved() != true {
		t.Errorf("Storage got did not get Champions from backend even though it should have")
	}
	if backend.getChampionsStored() != false {
		t.Errorf("Storage stored Champions in backend even though it shouldn't")
	}

	if riotClient.getChampionsRetrieved() != true {
		t.Errorf("Storage did not try to get the CHampions from client")
	}

	if len(championsListBackend) != len(actualChampions) {
		t.Error("Number of Champions not as expected")
	}
	if championsListBackend["3"].Name != actualChampions["3"].Name {
		t.Error("Champion does not match")
	}
	if championsListBackend["3"].ID != actualChampions["3"].ID {
		t.Error("Champion does not match")
	}
	if championsListBackend["3"].Timestamp != actualChampions["3"].Timestamp {
		t.Error("Timestamp does not match")
	}
	if championsListBackend["534"].Name != actualChampions["534"].Name {
		t.Error("Champion does not match")
	}
	if championsListBackend["534"].ID != actualChampions["534"].ID {
		t.Error("Champion does not match")
	}
	if championsListBackend["534"].Timestamp != actualChampions["534"].Timestamp {
		t.Error("Timestamp does not match")
	}

	// Should always get it from Client when backend fails
	riotClient.reset()
	riotClient.setChampions(championsListClient)
	backend.reset()
	championsListBackend["534"] = riotclient.Champion{Name: "BACKEND CHAMP 2",
		ID: "534", Timestamp: time.Now().Add(-time.Minute * time.Duration(5))}
	backend.setChampions(championsListBackend)
	backend.setChampionsTimeStamp(time.Now().Add(-time.Minute * time.Duration(5)))
	backend.setFailChampions(true)

	actualChampions = storage.GetChampions(false)

	if backend.getChampionsRetrieved() != false {
		t.Errorf("Storage got Champions from backend even though it shouldn't")
	}
	if backend.getChampionsStored() != true {
		t.Errorf("Received Champions where not stored in backend")
	}

	if riotClient.getChampionsRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	if len(championsListClient) != len(actualChampions) {
		t.Error("Number of Champions not as expected")
	}
	if championsListClient["43222"].Name != actualChampions["43222"].Name {
		t.Error("Champion does not match")
	}
	if championsListClient["43222"].ID != actualChampions["43222"].ID {
		t.Error("Champion does not match")
	}
	if championsListClient["43222"].Timestamp != actualChampions["43222"].Timestamp {
		t.Error("Timestamp does not match")
	}
	if championsListClient["123"].Name != actualChampions["123"].Name {
		t.Error("Champion does not match")
	}
	if championsListClient["123"].ID != actualChampions["123"].ID {
		t.Error("Champion does not match")
	}
	if championsListClient["123"].Timestamp != actualChampions["123"].Timestamp {
		t.Error("Timestamp does not match")
	}

	// Should return nil if everything fails
	riotClient.reset()
	riotClient.setChampions(championsListClient)
	backend.reset()
	backend.setChampions(championsListBackend)
	backend.setChampionsTimeStamp(time.Now().Add(-time.Minute * time.Duration(5)))
	backend.setFailChampions(true)
	riotClient.setFailChampions(true)

	actualChampions = storage.GetChampions(false)

	if backend.getChampionsRetrieved() != true {
		t.Errorf("Storage should have tried to get Champions from backend")
	}
	if backend.getChampionsStored() != false {
		t.Errorf("Even though client failed it stored data in Backend")
	}

	if riotClient.getChampionsRetrieved() != true {
		t.Errorf("Storage did not get the CHampions from client")
	}

	if actualChampions != nil {
		t.Error("Should have returned nil")
	}
}
