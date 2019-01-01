package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// SummonerResponse contains summary informations of a summoner
type SummonerResponse struct {
	Name           string                           `json:"name"`
	ProfileIcon    int                              `json:"profileIconId"`
	SummonerLevel  int64                            `json:"summonerLevel"`
	RevisionDate   int64                            `json:"revisionDate"`
	Timestamp      time.Time                        `json:"timestamp"`
	LeagueRankings riotclient.LeaguePositionDTOList `json:"leagues"`
}

func (s *Storage) prepareSummonerResponse(summonerName string) (*SummonerResponse, error) {
	summoner, err := s.GetSummonerByName(summonerName)
	if err != nil {
		return nil, fmt.Errorf("Error getting SummonerByName data")
	}

	summonerResponse := SummonerResponse{
		Name:          summoner.Name,
		ProfileIcon:   summoner.ProfileIcon,
		SummonerLevel: summoner.SummonerLevel,
		RevisionDate:  summoner.RevisionDate,
		Timestamp:     summoner.Timestamp,
	}

	leagues, err := s.riotClient.LeaguesForSummoner(summoner.ID)
	if err == nil {
		summonerResponse.LeagueRankings = *leagues
	}

	return &summonerResponse, nil
}

func (s *Storage) summonerByNameEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Received Rest API SummonerByName request from", r.RemoteAddr)

	var summonerName string
	if val, ok := r.URL.Query()["name"]; ok {
		if len(val[0]) == 0 {
			s.log.Warnf("name parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		summonerName = val[0]
	} else {
		s.log.Warnf("There was no name parameter in request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	summonerResponse, err := s.prepareSummonerResponse(summonerName)
	if err != nil {
		s.log.Warnf("Error preparing Summoner Response: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	out, err := json.Marshal(summonerResponse)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
