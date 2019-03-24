package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetFreeRotation returns the current free rotation from storage
// When forceUpdate is true, it always fetches it vom Riot API, if
// it is false it depends on how old the free rotation is if it gets fetched
// from Riot API
func (s *Storage) GetFreeRotation(forceUpdate bool) riotclient.FreeRotation {
	duration := time.Since(s.backend.GetFreeRotationTimeStamp())
	if duration.Minutes() > float64(s.config.MaxAgeChampionRotation) || forceUpdate {
		freeRotation, err := s.riotClient.ChampionRotations()
		if err != nil {
			s.log.Warnln(err)
			freeRotation, err := s.backend.GetFreeRotation()
			if err != nil {
				s.log.Warnln(err)
				return riotclient.FreeRotation{}
			}
			return *freeRotation
		}
		s.backend.StoreFreeRotation(freeRotation)
		return *freeRotation
	}

	freeRotation, err := s.backend.GetFreeRotation()
	if err != nil {
		s.log.Warnln(err)
		return riotclient.FreeRotation{}
	}
	return *freeRotation
}

func (s *Storage) freeRotationEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Free Rotation request from", r.RemoteAddr)
	freeRotation := s.GetFreeRotation(false)

	out, err := json.Marshal(freeRotation)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
