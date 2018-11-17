package riotclient

import (
	"net/http"
	"testing"

	"github.com/torlenor/alolstats/config"
)

func TestRiotClient(t *testing.T) {
	client, err := NewClient(&http.Client{}, config.RiotClient{})
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though APIVersion is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1"})
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though Key is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1", Key: "abcd"})
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though Region is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1", Key: "abcd", Region: "euw"})
	if err != nil || client == nil {
		t.Fatalf("Could not get a new client")
	}

	if client.IsRunning() != false {
		t.Fatalf("Client claims to be running even though we did not start it")
	}

	client.Start()
	if client.IsRunning() != true {
		t.Fatalf("Client is not running even though we started it")
	}

	client.Stop()
	if client.IsRunning() != false {
		t.Fatalf("Client not stopped even though we stopped it")
	}
}
