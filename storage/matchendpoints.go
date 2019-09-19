package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/utils"
)

func (s *Storage) getMatchEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Match request from", r.RemoteAddr)

	id, err := extractURLStringParameter(r.URL.Query(), "id")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}
	idNum, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		s.log.Warnf("Could not convert value %s to GameID", id)
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Provided game ID %s not a positive number", id)), http.StatusBadRequest)
		return
	}

	match, err := s.GetMatch(idNum)
	if err != nil {
		s.log.Warnf("Could not get match for id %d", idNum)
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not get match for id %d", idNum)), http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(match)
	if err != nil {
		s.log.Errorf("Could not marshal Match data to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
