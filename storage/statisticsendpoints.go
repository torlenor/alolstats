package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
)

func (s *Storage) championStatsByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API championByID request from", r.RemoteAddr)
	var champID string
	var gameVersion string
	var tier string

	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		champID = val[0]
	} else {
		s.log.Warnf("id parameter was missing in request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if val, ok := r.URL.Query()["gameversion"]; ok {
		if len(val) == 0 {
			s.log.Warnf("gameversion parameter was empty in request.")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		gameVersion = val[0]
	} else {
		s.log.Warnf("gameversion parameter was missing in request.")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if val, ok := r.URL.Query()["tier"]; ok {
		if len(val) == 0 {
			s.log.Debugf("tier parameter was empty in request, assuming ALL.")
			tier = "ALL"
		}
		tier = val[0]
	} else {
		tier = "ALL"
	}

	championStats, err := s.GetChampionStatsByIDGameVersionTier(champID, gameVersion, tier)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(championStats)
	if err != nil {
		s.log.Errorln(err)
		s.log.Errorf("Error in championByID with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) championStats(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API ChampionsStats request from", r.RemoteAddr)
	var gameVersion string
	var tier string

	if val, ok := r.URL.Query()["gameversion"]; ok {
		if len(val) == 0 {
			s.log.Warnf("gameversion parameter was empty in request.")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		gameVersion = val[0]
	} else {
		s.log.Warnf("gameversion parameter was missing in request.")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if val, ok := r.URL.Query()["tier"]; ok {
		if len(val) == 0 {
			s.log.Debugf("tier parameter was empty in request, assuming ALL.")
			tier = "ALL"
		}
		tier = val[0]
	} else {
		tier = "ALL"
	}

	statsSummary, err := s.GetChampionStatsSummaryByGameVersionTier(gameVersion, tier)
	if err != nil {
		s.log.Errorf("Error in ChampionsStats with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(statsSummary.ChampionsStatsSummary)
	if err != nil {
		s.log.Errorf("Error in ChampionsStats with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
