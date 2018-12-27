package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/torlenor/alolstats/riotclient"
	"github.com/torlenor/alolstats/utils"
)

// GetActiveGameBySummonerID returns the active game (live game) for the given Summoner ID
func (s *Storage) GetActiveGameBySummonerID(summonerID string) (*riotclient.CurrentGameInfoDTO, error) {
	s.log.Debugf("GetActiveGameBySummonerID(%s)", summonerID)

	currentGame, err := s.riotClient.ActiveGameBySummonerID(summonerID)
	if err != nil {
		s.log.Warnln("Could not get data from Riot API:", err)
		return nil, err
	}

	s.log.Debugf("Returned Active Game for SummonerID %s from Riot API", summonerID)
	return currentGame, nil
}

// GetFeaturedGames returns the currently features games from Riot
func (s *Storage) GetFeaturedGames() (*riotclient.FeaturedGamesDTO, error) {
	s.log.Debugf("GetFeaturedGames()")

	featuredGames, err := s.riotClient.FeaturedGames()
	if err != nil {
		s.log.Warnln("Could not get data from Riot API:", err)
		return nil, err
	}

	s.log.Debugf("Returned Featured Games from Riot API")
	return featuredGames, nil
}

func (s *Storage) getActiveGameBySummonerNameEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Active Game By Summoner Name request from", r.RemoteAddr)

	var summonerName string
	if val, ok := r.URL.Query()["name"]; ok {
		if len(val[0]) == 0 {
			s.log.Warnf("name parameter was empty in request")
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "name parameter was empty in request"), http.StatusBadRequest)
			return
		}
		summonerName = val[0]
	} else {
		s.log.Warnf("There was no name parameter in request")
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "No name parameter in request"), http.StatusBadRequest)
		return
	}

	summoner, err := s.GetSummonerByName(summonerName)
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
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) getFeaturedGamesEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Features Games request from", r.RemoteAddr)

	match, err := s.GetFeaturedGames()
	if err != nil {
		s.log.Errorf("getFeaturedGamesEndpoint error %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(match)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
