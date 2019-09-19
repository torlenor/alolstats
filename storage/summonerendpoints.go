package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"git.abyle.org/hps/alolstats/riotclient"
	"git.abyle.org/hps/alolstats/utils"
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

func (s *Storage) prepareSummonerResponse(summonerName string, forceUpdate bool) (*SummonerResponse, error) {
	summoner, err := s.GetSummonerByName(summonerName, forceUpdate)
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

	leagues, err := s.GetLeaguesForSummonerBySummonerID(summoner.ID, forceUpdate)
	if err == nil {
		summonerResponse.LeagueRankings = leagues
	} else {
		s.log.Warnf("Unable to get League Data for Summoner %s: %s", summoner.Name, err)
		summonerResponse.LeagueRankings = riotclient.LeaguePositionDTOList{}
	}

	return &summonerResponse, nil
}

func (s *Storage) summonerByNameEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API SummonerByName request from", r.RemoteAddr)

	summonerName, err := extractURLStringParameter(r.URL.Query(), "name")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	summonerResponse, err := s.prepareSummonerResponse(summonerName, false)
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
	}

	out, err := json.Marshal(summonerResponse)
	if err != nil {
		s.log.Errorf("Could not marshal Summoner data to JSON: %s", err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
