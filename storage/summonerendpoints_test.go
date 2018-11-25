package storage

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/riotclient"
)

func TestSummonerByNameEndpoint(t *testing.T) {

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

	summonerBackend := riotclient.Summoner{
		AccountID:    112345,
		ID:           212345,
		Name:         "Backend Summoner",
		Level:        10,
		RevisionDate: 312345,
		Timestamp:    time.Now().Add(-time.Minute * 10),
	}

	backend.setSummoner(summonerBackend)
	riotClient.setFailSummoner(true)

	req := httptest.NewRequest("GET", "http://example.com/endpoint?name=Backend%20Summoner", nil)
	w := httptest.NewRecorder()
	storage.summonerByNameEndpoint(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	receivedSummoner := riotclient.Summoner{}
	err = json.Unmarshal(body, &receivedSummoner)
	if err != nil {
		t.Fatalf("Decoding the json into the correct data struct was not possible: %s", err)
	}

	if !cmp.Equal(summonerBackend, receivedSummoner) {
		t.Error("Data not equal")
	}

	// name not given
	req = httptest.NewRequest("GET", "http://example.com/endpoint", nil)
	w = httptest.NewRecorder()
	storage.summonerByNameEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// name empty
	req = httptest.NewRequest("GET", "http://example.com/endpoint?name=", nil)
	w = httptest.NewRecorder()
	storage.summonerByNameEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// name does not give a summoner
	req = httptest.NewRequest("GET", "http://example.com/endpoint?name=blub", nil)
	w = httptest.NewRecorder()
	storage.summonerByNameEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}
}

func TestSummonerBySummonerIDEndpoint(t *testing.T) {

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

	summonerBackend := riotclient.Summoner{
		AccountID:    112345,
		ID:           212345,
		Name:         "Backend Summoner",
		Level:        10,
		RevisionDate: 312345,
		Timestamp:    time.Now().Add(-time.Minute * 10),
	}

	backend.setSummoner(summonerBackend)
	riotClient.setFailSummoner(true)

	req := httptest.NewRequest("GET", "http://example.com/endpoint?id=212345", nil)
	w := httptest.NewRecorder()
	storage.summonerBySummonerIDEndpoint(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	receivedSummoner := riotclient.Summoner{}
	err = json.Unmarshal(body, &receivedSummoner)
	if err != nil {
		t.Fatalf("Decoding the json into the correct data struct was not possible: %s", err)
	}

	if !cmp.Equal(summonerBackend, receivedSummoner) {
		t.Error("Data not equal")
	}

	// id not given
	req = httptest.NewRequest("GET", "http://example.com/endpoint", nil)
	w = httptest.NewRecorder()
	storage.summonerBySummonerIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// id empty
	req = httptest.NewRequest("GET", "http://example.com/endpoint?id=", nil)
	w = httptest.NewRecorder()
	storage.summonerBySummonerIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// id not a number
	req = httptest.NewRequest("GET", "http://example.com/endpoint?id=blub", nil)
	w = httptest.NewRecorder()
	storage.summonerBySummonerIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// id does not give a summoner
	req = httptest.NewRequest("GET", "http://example.com/endpoint?id=987654", nil)
	w = httptest.NewRecorder()
	storage.summonerBySummonerIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}
}

func TestSummonerByAccountIDEndpoint(t *testing.T) {

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

	summonerBackend := riotclient.Summoner{
		AccountID:    112345,
		ID:           212345,
		Name:         "Backend Summoner",
		Level:        10,
		RevisionDate: 312345,
		Timestamp:    time.Now().Add(-time.Minute * 10),
	}

	backend.setSummoner(summonerBackend)
	riotClient.setFailSummoner(true)

	req := httptest.NewRequest("GET", "http://example.com/endpoint?id=112345", nil)
	w := httptest.NewRecorder()
	storage.summonerByAccountIDEndpoint(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	receivedSummoner := riotclient.Summoner{}
	err = json.Unmarshal(body, &receivedSummoner)
	if err != nil {
		t.Fatalf("Decoding the json into the correct data struct was not possible: %s", err)
	}

	if !cmp.Equal(summonerBackend, receivedSummoner) {
		t.Error("Data not equal")
	}

	// id not given
	req = httptest.NewRequest("GET", "http://example.com/endpoint", nil)
	w = httptest.NewRecorder()
	storage.summonerByAccountIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get correct status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// id empty
	req = httptest.NewRequest("GET", "http://example.com/endpoint?id=", nil)
	w = httptest.NewRecorder()
	storage.summonerByAccountIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// id not a number
	req = httptest.NewRequest("GET", "http://example.com/endpoint?id=blub", nil)
	w = httptest.NewRecorder()
	storage.summonerByAccountIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}

	// id does not give a summoner
	req = httptest.NewRequest("GET", "http://example.com/endpoint?id=987654", nil)
	w = httptest.NewRecorder()
	storage.summonerByAccountIDEndpoint(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Did not get a valid status, status received was: %d", resp.StatusCode)
	}

	if strings.Compare(string(body), http.StatusText(http.StatusBadRequest)+"\n") != 0 {
		t.Fatalf("Did not get a valid text, status text received was: %s, should have been: %s", string(body), http.StatusText(http.StatusBadRequest))
	}
}
