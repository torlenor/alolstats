package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetSummonerByName returns a Summoner identified by name
func (s *Storage) GetSummonerByName(name string) (riotclient.Summoner, error) {
	duration := time.Since(s.backend.GetSummonerByNameTimeStamp(name))
	if duration.Minutes() > float64(s.config.MaxAgeSummoner) {
		summoner, err := s.riotClient.SummonerByName(name)
		if err != nil {
			s.log.Warnln("Could not get new data from Client, trying to get it from Storage instead", err)
			summoner, err := s.backend.GetSummonerByName(name)
			if err != nil {
				s.log.Warnln("Could not get data from either Storage nor Client:", err)
				return riotclient.Summoner{}, err
			}
			s.log.Debugf("Returned Summoner %s from Storage", name)
			return summoner, nil
		}
		err = s.backend.StoreSummoner(summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner %s from Riot API", name)
		return *summoner, nil
	}
	summoner, err := s.backend.GetSummonerByName(name)
	if err != nil {
		summoner, errClient := s.riotClient.SummonerByName(name)
		if errClient != nil {
			s.log.Warnln("Could not get data from either Storage nor Client:", errClient)
			return riotclient.Summoner{}, errClient
		}
		s.log.Warnln("Could not get Summoner from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummoner(summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner %s from Riot API", name)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner %s from Storage", name)
	return summoner, nil
}

// GetSummonerBySummonerID returns a Summoner identified by Summoner ID
func (s *Storage) GetSummonerBySummonerID(summonerID uint64) (riotclient.Summoner, error) {
	duration := time.Since(s.backend.GetSummonerBySummonerIDTimeStamp(summonerID))
	if duration.Minutes() > float64(s.config.MaxAgeSummoner) {
		summoner, err := s.riotClient.SummonerBySummonerID(summonerID)
		if err != nil {
			s.log.Warnln("Could not get new data from Client, trying to get it from Storage instead", err)
			summoner, err := s.backend.GetSummonerBySummonerID(summonerID)
			if err != nil {
				s.log.Warnln("Could not get data from either Storage nor Client:", err)
				return riotclient.Summoner{}, err
			}
			s.log.Debugf("Returned Summoner with SummonerID %d from Storage", summonerID)
			return summoner, nil
		}
		err = s.backend.StoreSummoner(summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with SummonerID %d from Riot API", summonerID)
		return *summoner, nil
	}
	summoner, err := s.backend.GetSummonerBySummonerID(summonerID)
	if err != nil {
		summoner, errClient := s.riotClient.SummonerBySummonerID(summonerID)
		if errClient != nil {
			s.log.Warnln("Could not get data from either Storage nor Client:", errClient)
			return riotclient.Summoner{}, errClient
		}
		s.log.Warnln("Could not get Summoner from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummoner(summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with SummonerID %d from Riot API", summonerID)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner with SummonerID %d from Storage", summonerID)
	return summoner, nil
}

// GetSummonerByAccountID returns a Summoner identified by Account ID
func (s *Storage) GetSummonerByAccountID(accountID uint64) (riotclient.Summoner, error) {
	duration := time.Since(s.backend.GetSummonerByAccountIDTimeStamp(accountID))
	if duration.Minutes() > float64(s.config.MaxAgeSummoner) {
		summoner, err := s.riotClient.SummonerByAccountID(accountID)
		if err != nil {
			s.log.Warnln("Could not get new data from Client, trying to get it from Storage instead", err)
			summoner, err := s.backend.GetSummonerByAccountID(accountID)
			if err != nil {
				s.log.Warnln("Could not get data from either Storage nor Client:", err)
				return riotclient.Summoner{}, err
			}
			s.log.Debugf("Returned Summoner with AccountID %d from Storage", accountID)
			return summoner, nil
		}
		err = s.backend.StoreSummoner(summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with AccountID %d from Riot API", accountID)
		return *summoner, nil
	}
	summoner, err := s.backend.GetSummonerByAccountID(accountID)
	if err != nil {
		summoner, errClient := s.riotClient.SummonerByAccountID(accountID)
		if errClient != nil {
			s.log.Warnln("Could not get data from either Storage nor Client:", errClient)
			return riotclient.Summoner{}, errClient
		}
		s.log.Warnln("Could not get Summoner from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummoner(summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with AccountID %d from Riot API", accountID)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner with AccountID %d from Storage", accountID)
	return summoner, nil
}

func (s *Storage) summonerByNameEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Received Rest API SummonerByName request from", r.RemoteAddr)

	var summonerName string
	if val, ok := r.URL.Query()["name"]; ok {
		if len(val) == 0 {
			s.log.Warnf("name parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		summonerName = val[0]
	}

	summoner, err := s.GetSummonerByName(summonerName)
	if err != nil {
		s.log.Warnf("Error getting SummonerByName data")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(summoner)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) summonerBySummonerIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API summonerBySummonerID request from", r.RemoteAddr)
	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			s.log.Warnf("Could not convert value %s to SummonerID", val)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		summoner, err := s.GetSummonerBySummonerID(id)
		if err != nil {
			s.log.Warnf("Error getting SummonerBySummonerID data")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		out, err := json.Marshal(summoner)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	}

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) summonerByAccountIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API summonerByAccountID request from", r.RemoteAddr)
	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			s.log.Warnf("Could not convert value %s to AccountID", val)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		summoner, err := s.GetSummonerByAccountID(id)
		if err != nil {
			s.log.Warnf("Error getting SummonerByAccountID data")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		out, err := json.Marshal(summoner)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	}

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
