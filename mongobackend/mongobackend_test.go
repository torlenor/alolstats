package mongobackend

import (
	"testing"

	"github.com/torlenor/alolstats/config"
)

func TestCreatingNewMongoBackend(t *testing.T) {
	cfg := config.MongoBackend{URL: "mongodb://mongo/ut", Database: "ut"}
	backend, err := NewBackend(cfg)
	if err != nil || backend == nil {
		t.Fatalf("Could not get a new Mongo Backend: %s", err)
	}
}
