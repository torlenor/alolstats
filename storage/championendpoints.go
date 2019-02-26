package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/torlenor/alolstats/utils"
)

var fallbackGameVersion = "9.4"

func checkParamterForceUpdate(values url.Values) bool {
	if val, ok := values["forceupdate"]; ok {
		if len(val[0]) == 0 {
			return false
		}
		if strings.ToLower(val[0]) == "true" {
			return true
		}
	}
	return false
}

func (s *Storage) championsEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Champions request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	champions := s.GetChampions(checkParamterForceUpdate(r.URL.Query()))

	for key, val := range champions {
		gameVersion := fallbackGameVersion
		if val, ok := r.URL.Query()["gameversion"]; ok {
			if len(val[0]) == 0 {
				s.log.Warnf("gameversion parameter was empty in request, using default %s", gameVersion)
			}
			gameVersion = val[0]
		}

		stats, err := s.GetChampionStatsByIDGameVersion(val.ID, gameVersion)
		if err == nil {
			val.Roles = stats.Roles
		}
		champions[key] = val
	}

	out, err := json.Marshal(champions)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champions to JSON")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))
}

func (s *Storage) championByKeyEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Champion by Key request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	if val, ok := r.URL.Query()["key"]; ok {
		if len(val[0]) == 0 {
			s.log.Warnf("key parameter was empty in request")
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "key parameter was empty in request"), http.StatusBadRequest)
			return
		}
		key := val[0]

		champion, err := s.GetChampionByKey(key, checkParamterForceUpdate(r.URL.Query()))
		if err != nil {
			s.log.Warnf("Could not get Champion with Key %s", key)
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not get Champion with Key %s", key)), http.StatusBadRequest)
			return
		}

		gameVersion := fallbackGameVersion
		if val, ok := r.URL.Query()["gameversion"]; ok {
			if len(val[0]) == 0 {
				s.log.Warnf("gameversion parameter was empty in request, using default %s", gameVersion)
			}
			gameVersion = val[0]
		}

		stats, err := s.GetChampionStatsByIDGameVersion(champion.ID, gameVersion)
		if err == nil {
			champion.Roles = stats.Roles
		}

		out, err := json.Marshal(champion)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	} else {
		s.log.Warnf("key parameter was missing in request")
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "key parameter was missing in request"), http.StatusBadRequest)
		return
	}
}

func (s *Storage) championByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Champion by ID request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	if val, ok := r.URL.Query()["id"]; ok {
		if len(val[0]) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "id parameter was empty in request"), http.StatusBadRequest)
			return
		}
		id := val[0]

		champion, err := s.GetChampionByID(id, checkParamterForceUpdate(r.URL.Query()))
		if err != nil {
			s.log.Warnf("Could not get Champion with ID %s", id)
			http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not get Champion with ID %s", id)), http.StatusBadRequest)
			return
		}

		gameVersion := fallbackGameVersion
		if val, ok := r.URL.Query()["gameversion"]; ok {
			if len(val[0]) == 0 {
				s.log.Warnf("gameversion parameter was empty in request, using default %s", gameVersion)
			}
			gameVersion = val[0]
		}

		stats, err := s.GetChampionStatsByIDGameVersion(champion.ID, gameVersion)
		if err == nil {
			champion.Roles = stats.Roles
		}

		out, err := json.Marshal(champion)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	} else {
		s.log.Warnf("id parameter was missing in request")
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, "id parameter was missing in request"), http.StatusBadRequest)
		return
	}
}
