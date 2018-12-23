package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

func (s *Storage) getFreeRotation() riotclient.FreeRotation {
	duration := time.Since(s.backend.GetFreeRotationTimeStamp())
	if duration.Minutes() > float64(s.config.MaxAgeChampionRotation) {
		freeRotation, err := s.riotClient.ChampionRotations()
		if err != nil {
			s.log.Warnln(err)
			freeRotation, err := s.backend.GetFreeRotation()
			if err != nil {
				s.log.Warnln(err)
				return riotclient.FreeRotation{}
			}
			return freeRotation
		}
		s.backend.StoreFreeRotation(*freeRotation)
		return *freeRotation
	}

	freeRotation, err := s.backend.GetFreeRotation()
	if err != nil {
		s.log.Warnln(err)
		return riotclient.FreeRotation{}
	}
	return freeRotation
}

func (s *Storage) freeRotationEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Free Rotation request from", r.RemoteAddr)
	freeRotation := s.getFreeRotation()

	out, err := json.Marshal(freeRotation)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
