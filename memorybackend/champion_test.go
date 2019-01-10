package memorybackend

import (
	"testing"
)

func TestStoreAndRetreiveChampionsList(t *testing.T) {
	// backend, err := NewBackend()
	// if err != nil || backend == nil {
	// 	t.Fatalf("Could not get a new Memory Backend: %s", err)
	// }

	// older, _ := time.Parse(
	// 	time.RFC3339,
	// 	"2018-11-01T22:08:41+00:00")
	// newer, _ := time.Parse(
	// 	time.RFC3339,
	// 	"2018-11-01T22:18:41+00:00")

	// championsList := riotclient.ChampionList{}
	// championsList.Champions = make(map[string]riotclient.Champion)
	// championsList.Champions["3"] = riotclient.Champion{Name: "TEST CHAMP",
	// 	ID: "3", Timestamp: newer}
	// championsList.Champions["534"] = riotclient.Champion{Name: "TEST CHAMP 2",
	// 	ID: "534", Timestamp: older}

	// backend.StoreChampions(&championsList)
	// if err != nil {
	// 	t.Fatalf("Failed storing ChampionList: %s", err)
	// }

	// if backend.GetChampionsTimeStamp() != older {
	// 	t.Error("GetChampionsTimeStamp() TimeStamp does not reflect the TimeStamp of the oldest entry")
	// }

	// retreivedChampionList, err := backend.GetChampions()
	// if err != nil {
	// 	t.Fatalf("Could not retrieve ChampionList: %s", err)
	// }

	// if len(retreivedChampionList.Champions) != len(championsList.Champions) {
	// 	t.Error("Number of Champions not as expected")
	// }
	// if championsList.Champions["3"].Name != retreivedChampionList.Champions["3"].Name {
	// 	t.Error("Champion does not match")
	// }
	// if championsList.Champions["3"].ID != retreivedChampionList.Champions["3"].ID {
	// 	t.Error("Champion does not match")
	// }
	// if retreivedChampionList.Champions["3"].Timestamp != newer {
	// 	t.Error("Timestamp does not match")
	// }
	// if championsList.Champions["534"].Name != retreivedChampionList.Champions["534"].Name {
	// 	t.Error("Champion does not match")
	// }
	// if championsList.Champions["534"].ID != retreivedChampionList.Champions["534"].ID {
	// 	t.Error("Champion does not match")
	// }
	// if retreivedChampionList.Champions["534"].Timestamp != older {
	// 	t.Error("Timestamp does not match")
	// }
}
