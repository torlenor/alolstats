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
		s.backend.StoreChampions(*champions)
		return *champions
	}
	champions, err := s.backend.GetChampions()
	if err != nil {
		s.log.Warnln(err)
		return riotclient.ChampionList{}
	}
	return champions
}

func (s *Storage) championsEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Received Rest API Champions request from", r.RemoteAddr)
	champions := s.getChampions()

	out, err := json.Marshal(champions)
	if err != nil {
		s.log.Warnln(err)
		io.WriteString(w, `{"Error":"Error getting champions data"}`)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
