package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/torlenor/alolstats/utils"

	"github.com/torlenor/alolstats/riotclient"
)

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

func (s *Storage) championByKeyEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Received Rest API Champion by Key request from", r.RemoteAddr)

	if val, ok := r.URL.Query()["key"]; ok {
		if len(val[0]) == 0 {
			s.log.Warnf("key parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		champions := s.GetChampions()

		champion := riotclient.Champion{}

		for _, champion = range champions.Champions {
			if val[0] == champion.Key {
				break
			}
		}

		out, err := json.Marshal(champion)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	}

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) championByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Received Rest API Champion by ID request from", r.RemoteAddr)

	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "id parameter was empty in request"), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			s.log.Warnf("Could not convert value %s to Champion ID", val)
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not convert value %s to Champion ID", val)), http.StatusBadRequest)
			return
		}

		champion, err := s.GetChampionByID(uint32(id))
		if err != nil {
			s.log.Warnf("Could not get Champion with ID %d", id)
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not get Champion with ID %d", id)), http.StatusBadRequest)
			return
		}

		out, err := json.Marshal(champion)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	}

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
