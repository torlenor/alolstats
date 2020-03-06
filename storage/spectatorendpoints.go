package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/utils"
)

// ActiveGameResponse contains the prepare ActiveGame information for a API response
type ActiveGameResponse struct {
}

func (s *Storage) getActiveGameBySummonerNameEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Active Game By Summoner Name request from", r.RemoteAddr)

	summonerName, err := extractURLStringParameter(r.URL.Query(), "name")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	summoner, err := s.GetSummonerByName(summonerName, false)
	if err != nil {
		s.log.Warnf("Error getting SummonerByName data")
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "Summoner "+summonerName+" not found"), http.StatusBadRequest)
		return
	}

	activeGame, err := s.GetActiveGameBySummonerID(summoner.ID)
	if err != nil {
		s.log.Errorf("getActiveGameBySummonerNameEndpoint error %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusNotFound, "Active Game for Summoner "+summonerName+" not found"), http.StatusNotFound)
		return
	}

	out, err := json.Marshal(activeGame)
	if err != nil {
		s.log.Errorf("Could not marshal Active Game to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) getFeaturedGamesEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Features Games request from", r.RemoteAddr)

	featuredGames, err := s.GetFeaturedGames()
	if err != nil {
		s.log.Errorf("Could not get featured games: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(featuredGames)
	if err != nil {
		s.log.Errorf("Could not marshal features games to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
