package memorybackend

import (
	"testing"

	"github.com/torlenor/alolstats/riotclient"
)

func TestStoreAndRetreiveMatch(t *testing.T) {
	backend, err := NewBackend()
	if err != nil || backend == nil {
		t.Fatalf("Could not get a new Memory Backend: %s", err)
	}
	match := &riotclient.MatchDTO{}
	match.GameID = 1234
	match.GameDuration = 3434

	err = backend.StoreMatch(match)
	if err != nil {
		t.Errorf("Failed storing a match: %s", err)
	}

	_, err = backend.GetMatch(uint64(3456))
	if err == nil {
		t.Errorf("Should have gotten an error because Match cannot exist in storage")
	}

	retreivedMatch, err := backend.GetMatch(uint64(match.GameID))
	if err != nil {
		t.Fatalf("Could not retrieve Match which should exist: %s", err)
	}
	if retreivedMatch.GameID != match.GameID {
		t.Error("GameID does not match")
	}
	if retreivedMatch.GameDuration != match.GameDuration {
		t.Error("GameDuration does not match")
	}
}
