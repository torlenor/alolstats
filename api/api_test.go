package api

import (
	"testing"

	"github.com/torlenor/alolstats/config"
)

func TestCreatingNewAPI(t *testing.T) {
	config := config.API{}
	backend, err := NewAPI(config)
	if err != nil || backend == nil {
		t.Fatalf("Could not get a new Riot API: %s", err)
	}
}
