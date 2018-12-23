package riotclientv3

import (
	"net/http"
	"testing"

	"github.com/torlenor/alolstats/riotclient/datadragon"
	"github.com/torlenor/alolstats/riotclient/ratelimit"

	"github.com/torlenor/alolstats/config"
)

func TestRiotClient(t *testing.T) {
	httpClient := &http.Client{}
	ddragon, _ := riotclientdd.New(httpClient, config.RiotClient{})
	rateLimit, _ := riotclientrl.New()
	client, err := NewClient(httpClient, config.RiotClient{}, ddragon, rateLimit)
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though APIVersion is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1"}, ddragon, rateLimit)
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though Key is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1", Key: "abcd"}, ddragon, rateLimit)
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though Region is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1", Key: "abcd", Region: "euw"}, ddragon, rateLimit)
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
	client.Start()
	if client.IsRunning() != true {
		t.Fatalf("Client should still be running if we start it again")
	}

	client.Stop()
	if client.IsRunning() != false {
		t.Fatalf("Client not stopped even though we stopped it")
	}
	client.Stop()
	if client.IsRunning() != false {
		t.Fatalf("Client should still be stopped if we stop it again")
	}
}
