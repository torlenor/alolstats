package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/utils"
)

func (s *Storage) championsEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Champions request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	gameVersion, err := extractURLStringParameter(r.URL.Query(), "gameversion")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	tier, err := extractURLStringParameter(r.URL.Query(), "tier")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	queue, err := extractURLStringParameter(r.URL.Query(), "queue")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	champions := s.GetChampions(false)

	for key, val := range champions {
		stats, err := s.GetChampionStatsByIDGameVersionTierQueue(val.ID, gameVersion, tier, queue)
		if err == nil {
			val.Roles = stats.Roles
		}
		champions[key] = val
	}

	out, err := json.Marshal(champions)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champions to JSON")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))
}

func (s *Storage) championByKeyEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Champion by Key request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	key, err := extractURLStringParameter(r.URL.Query(), "key")
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

	champion, err := s.GetChampionByKey(key, false)
	if err != nil {
		s.log.Warnf("Could not get Champion with Key %s", key)
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not get Champion with Key %s", key)), http.StatusBadRequest)
		return
	}

	stats, err := s.GetChampionStatsByIDGameVersionTierQueue(champion.ID, gameVersion, tier, queue)
	if err == nil {
		champion.Roles = stats.Roles
	}

	out, err := json.Marshal(champion)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))
}

func (s *Storage) championByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Champion by ID request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	var id string
	var gameVersion string
	var queue string
	var tier string

	id, err := extractURLStringParameter(r.URL.Query(), "id")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	gameVersion, err = extractURLStringParameter(r.URL.Query(), "gameversion")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	queue, err = extractURLStringParameter(r.URL.Query(), "queue")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	tier, err = extractURLStringParameter(r.URL.Query(), "tier")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	champion, err := s.GetChampionByID(id, false)
	if err != nil {
		s.log.Warnf("Could not get Champion with ID %s", id)
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not get Champion with ID %s", id)), http.StatusBadRequest)
		return
	}

	stats, err := s.GetChampionStatsByIDGameVersionTierQueue(champion.ID, gameVersion, tier, queue)
	if err == nil {
		champion.Roles = stats.Roles
	}

	out, err := json.Marshal(champion)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))
}
