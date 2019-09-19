package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
)

func (s *Storage) freeRotationEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Free Rotation request from", r.RemoteAddr)
	freeRotation := s.GetFreeRotation(false)

	out, err := json.Marshal(freeRotation)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
