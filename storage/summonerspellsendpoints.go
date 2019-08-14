package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"git.abyle.org/hps/alolstats/utils"
)

func (s *Storage) summonerSpellsEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Summoner Spells request from", r.RemoteAddr)
	atomic.AddUint64(&s.stats.handledRequests, 1)

	summonerSpells := s.GetSummonerSpells(checkParamterForceUpdate(r.URL.Query()))

	out, err := json.Marshal(summonerSpells)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Summoner Spells to JSON")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))
}
