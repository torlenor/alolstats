package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/torlenor/alolstats/riotclient"
)

func (s *Storage) getMatch(id uint64) (riotclient.Match, error) {
	match, err := s.backend.GetMatch(id)
	if err != nil {
		s.log.Warnln(err)
		match, err := s.riotClient.MatchByID(id)
		if err != nil {
			s.log.Warnln(err)
			return riotclient.Match{}, err
		}
		s.log.Debugf("Returned Match %d from Riot API", id)
		s.backend.StoreMatch(match)
		return *match, nil
	}
	s.log.Debugf("Returned Match %d from Storage", id)
	return match, nil
}

// GetStoredMatchesByGameVersion gets all matches for a specific game version
func (s *Storage) GetStoredMatchesByGameVersion(gameVersion string) (riotclient.Matches, error) {
	return s.backend.GetMatchesByGameVersion(gameVersion)
}

func (s *Storage) getMatchEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Match request from", r.RemoteAddr)
	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			s.log.Warnf("Could not convert value %s to GameID", val)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		match, err := s.getMatch(id)
		if err != nil {
			s.log.Warnf("Could not get match for id=%d", id)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		out, err := json.Marshal(match)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	}

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) getMatchesEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Matches request from", r.RemoteAddr)
	if val, ok := r.URL.Query()["gameversion"]; ok {
		if len(val) == 0 {
			s.log.Warnf("gameversion parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		matches, err := s.GetStoredMatchesByGameVersion(val[0])
		if err != nil {
			s.log.Warnf("Could not get matches for game version = %s", val[0])
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		out, err := json.Marshal(matches)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	}

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
