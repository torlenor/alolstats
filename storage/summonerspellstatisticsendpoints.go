package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/utils"
)

func (s *Storage) summonerSpellsStatsByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API summonerSpellsStatsByIDEndpoint request from", r.RemoteAddr)

	id, err := extractURLStringParameter(r.URL.Query(), "id")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	gameVersion, err := extractURLStringParameter(r.URL.Query(), "gameversion")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	queue, err := extractURLStringParameter(r.URL.Query(), "queue")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	tier, err := extractURLStringParameter(r.URL.Query(), "tier")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	summonerSpellsStats, err := s.GetSummonerSpellsStatsByIDGameVersionTierQueue(id, gameVersion, tier, queue)
	if err != nil {
		s.log.Errorf("Error in championByID with request %s: %s", r.URL.String(), err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("No data")), http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(summonerSpellsStats)
	if err != nil {
		s.log.Errorf("Error in championByID with request %s: %s", r.URL.String(), err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
		return
	}

	s.log.Info(string(out))

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
