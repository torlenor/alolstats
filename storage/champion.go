package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetChampions returns a list of all currently known champions
func (s *Storage) GetChampions() riotclient.ChampionList {
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
			s.log.Debugf("Could not get Champions from Client, returning from Storage Backend instead")
			return champions
		}
		err = s.backend.StoreChampions(*champions)
		if err != nil {
			s.log.Warnln("Could not store Champions in storage backend:", err)
		}
		s.log.Debugf("Returned Champions from Client")
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
	s.log.Debugf("Returned Champions from Storage Backend")
	return champions
}

func (s *Storage) championsEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Received Rest API Champions request from", r.RemoteAddr)
	champions := s.GetChampions()

	out, err := json.Marshal(champions)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
