package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/utils"
)

func (s *Storage) getRunesReforgedEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Summoner Spells request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	gameVersion, err := extractURLStringParameter(r.URL.Query(), "gameversion")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	language, err := extractURLStringParameter(r.URL.Query(), "language")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	items, err := s.GetRunesReforged(gameVersion, language)
	if err != nil {
		s.log.Errorf("Could not get RunesReforged for gameversion %s, language %s: %s", gameVersion, language, err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(items)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))
}
