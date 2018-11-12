package memorybackend

import (
	"testing"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

func TestStoreAndRetreiveChampionsList(t *testing.T) {
	backend, err := NewBackend()
	if err != nil || backend == nil {
		t.Fatalf("Could not get a new Memory Backend: %s", err)
	}
	championsList := riotclient.ChampionList{}
	championsList.Type = "TYPE"
	championsList.Version = "32"
	championsList.Champions = make(map[string]riotclient.Champion)
	championsList.Champions["3"] = riotclient.Champion{Name: "TEST CHAMP",
		ID: "3"}
	championsList.Champions["534"] = riotclient.Champion{Name: "TEST CHAMP 2",
		ID: "534"}

	t1, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	championsList.Timestamp = t1

	backend.StoreChampions(championsList)
	if err != nil {
		t.Fatalf("Failed storing ChampionList: %s", err)
	}

	if backend.GetChampionsTimeStamp() != t1 {
		t.Error("GetChampionsTimeStamp() TimeStamp does not match")
	}

	retreivedChampionList, err := backend.GetChampions()
	if err != nil {
		t.Fatalf("Could not retrieve ChampionList: %s", err)
	}

	if championsList.Type != retreivedChampionList.Type {
		t.Error("Type does not match")
	}
	if championsList.Version != retreivedChampionList.Version {
		t.Error("Version does not match")
	}

	if len(retreivedChampionList.Champions) != len(championsList.Champions) {
		t.Error("Number of Champions not as expected")
	}
	if championsList.Champions["3"].Name != retreivedChampionList.Champions["3"].Name {
		t.Error("Champion does not match")
	}
	if championsList.Champions["3"].ID != retreivedChampionList.Champions["3"].ID {
		t.Error("Champion does not match")
	}
	if championsList.Champions["534"].Name != retreivedChampionList.Champions["534"].Name {
		t.Error("Champion does not match")
	}
	if championsList.Champions["534"].ID != retreivedChampionList.Champions["534"].ID {
		t.Error("Champion does not match")
	}

	if retreivedChampionList.Timestamp != t1 {
		t.Error("Timestamp does not match")
	}
}
