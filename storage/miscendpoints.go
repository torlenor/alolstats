package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/utils"
)

func (s *Storage) storageSummaryEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API StorageSummary request from", r.RemoteAddr)

	storageSummary, err := s.backend.GetStorageSummary()
	if err != nil {
		s.log.Errorf("Could not get Storage Summary from backend: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(storageSummary)
	if err != nil {
		s.log.Errorf("Could not marshal Storage Summary to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) getKnownVersionsEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Known Versions request from", r.RemoteAddr)

	ver, err := s.backend.GetKnownGameVersions()
	if err != nil {
		s.log.Errorf("Could not get Known Game Versions from backend: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(ver)
	if err != nil {
		s.log.Errorf("Could not marshal Known Versions to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) getStatLeaguesEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API StatLeagues request from", r.RemoteAddr)

	type leagues struct {
		Leagues []string `json:"leagues"`
	}

	lea := leagues{Leagues: []string{"All", "Master", "Diamond", "Platinum", "Gold", "Silver", "Bronze"}}

	out, err := json.Marshal(lea)
	if err != nil {
		s.log.Errorf("Could not marshal Stat Leagues to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) getStatQueuesEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API StatQueues request from", r.RemoteAddr)

	type queues struct {
		Queues []string `json:"queues"`
	}

	que := queues{Queues: []string{"RANKED_SOLO", "RANKED_FLEX", "NORMAL_BLIND", "NORMAL_DRAFT"}}

	out, err := json.Marshal(que)
	if err != nil {
		s.log.Errorf("Could not marshal Stat Queues to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
