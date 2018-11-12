package storage

import (
	"testing"

	"github.com/torlenor/alolstats/config"
)

func TestCreatingNewStorage(t *testing.T) {
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
}
