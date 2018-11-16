package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

func (s *Storage) getChampions() riotclient.ChampionList {
	duration := time.Since(s.backend.GetChampionsTimeStamp())
	if duration.Minutes() > float64(s.config.MaxAgeChampion) {
		champions, err := s.riotClient.Champions()
		if err != nil {
			s.log.Warnln(err)
			champions, err := s.backend.GetChampions()
			if err != nil {
				s.log.Warnln(err)
				return riotclient.ChampionList{}
			}
			return champions
		}
		err = s.backend.StoreChampions(*champions)
		if err != nil {
			s.log.Warnln("Could not store Champions in storage backend:", err)
		}
		return *champions
	}
	champions, err := s.backend.GetChampions()
	if err != nil {
		champions, errClient := s.riotClient.Champions()
		if errClient != nil {
			s.log.Warnln(err)
			return riotclient.ChampionList{}
		}
		s.log.Warnln("Could not get Champions from storage backend, returning from Client instead:", err)
		err = s.backend.StoreChampions(*champions)
		if err != nil {
			s.log.Warnln("Could not store Champions in storage backend:", err)
		}
		return *champions
	}
	return champions
}

func (s *Storage) championsEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Received Rest API Champions request from", r.RemoteAddr)
	champions := s.getChampions()

	out, err := json.Marshal(champions)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
