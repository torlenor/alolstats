package storage

import (
	"testing"

	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/riotclient"
)

func TestCreatingNewStorage(t *testing.T) {
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
}
