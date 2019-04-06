package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

func (s *Storage) championStatsByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API championByID request from", r.RemoteAddr)
	var champID string
	var gameVersion string
	var tier string
	var queue string

	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		champID = val[0]
	} else {
		s.log.Warnf("id parameter was missing in request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if val, ok := r.URL.Query()["gameversion"]; ok {
		if len(val) == 0 {
			s.log.Warnf("gameversion parameter was empty in request.")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		gameVersion = val[0]
	} else {
		s.log.Warnf("gameversion parameter was missing in request.")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if val, ok := r.URL.Query()["tier"]; ok {
		if len(val) == 0 {
			s.log.Debugf("tier parameter was empty in request, assuming ALL.")
			tier = "ALL"
		}
		tier = val[0]
	} else {
		tier = "ALL"
	}

	if val, ok := r.URL.Query()["queue"]; ok {
		if len(val) == 0 {
			s.log.Debugf("queue parameter was empty in request, assuming ALL.")
			queue = "RANKED_SOLO"
		}
		queue = val[0]
	} else {
		queue = "RANKED_SOLO"
	}

	championStats, err := s.GetChampionStatsByIDGameVersionTierQueue(champID, gameVersion, tier, queue)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(championStats)
	if err != nil {
		s.log.Errorln(err)
		s.log.Errorf("Error in championByID with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) championStats(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API ChampionsStats request from", r.RemoteAddr)
	var gameVersion string
	var tier string
	var queue string

	if val, ok := r.URL.Query()["gameversion"]; ok {
		if len(val) == 0 {
			s.log.Warnf("gameversion parameter was empty in request.")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		gameVersion = val[0]
	} else {
		s.log.Warnf("gameversion parameter was missing in request.")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if val, ok := r.URL.Query()["tier"]; ok {
		if len(val) == 0 {
			s.log.Debugf("tier parameter was empty in request, assuming ALL.")
			tier = "ALL"
		}
		tier = val[0]
	} else {
		tier = "ALL"
	}
	if val, ok := r.URL.Query()["queue"]; ok {
		if len(val) == 0 {
			s.log.Debugf("queue parameter was empty in request, assuming ALL.")
			queue = "RANKED_SOLO"
		}
		queue = val[0]
	} else {
		queue = "RANKED_SOLO"
	}

	statsSummary, err := s.GetChampionStatsSummaryByGameVersionTierQueue(gameVersion, tier, queue)
	if err != nil {
		s.log.Errorf("Error in ChampionsStats with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(statsSummary.ChampionsStatsSummary)
	if err != nil {
		s.log.Errorf("Error in ChampionsStats with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) championStatsHistoryByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API championStatsHistoryByID request from", r.RemoteAddr)
	var champID string
	var tier string
	var queue string

	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		champID = val[0]
	} else {
		s.log.Warnf("id parameter was missing in request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if val, ok := r.URL.Query()["tier"]; ok {
		if len(val) == 0 {
			s.log.Debugf("tier parameter was empty in request, assuming ALL.")
			tier = "ALL"
		}
		tier = val[0]
	} else {
		tier = "ALL"
	}

	if val, ok := r.URL.Query()["queue"]; ok {
		if len(val) == 0 {
			s.log.Debugf("queue parameter was empty in request, assuming ALL.")
			queue = "RANKED_SOLO"
		}
		queue = val[0]
	} else {
		queue = "RANKED_SOLO"
	}

	gameVersions, err := s.GetKnownGameVersions()
	if err != nil {
		s.log.Errorf("Error in championStatsHistoryByID with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	championStatsHistory := ChampionStatsHistory{}
	championStatsHistory.HistoryPeRrole = make(map[string]ChampionStatsPerRoleSingleHistory)
	for _, gameVersion := range gameVersions.Versions {
		championStats, err := s.GetChampionStatsByIDGameVersionTierQueue(champID, gameVersion, tier, queue)
		if err != nil {
			continue
		}

		if championStats.SampleSize == 0 {
			s.log.Infof("ChampionStatsHistoryByID for Champion ID %s, Tier %s: Skipping Game Version %s, no data", champID, tier, gameVersion)
			continue
		}

		championStatsHistory.Versions = append(championStatsHistory.Versions, gameVersion)
		championStatsHistory.PickRateHistory = append(championStatsHistory.PickRateHistory, championStats.PickRate)
		championStatsHistory.BanRateHistory = append(championStatsHistory.BanRateHistory, championStats.BanRate)
		championStatsHistory.WinRateHistory = append(championStatsHistory.WinRateHistory, championStats.WinRate)
		championStatsHistory.AvgKHistory = append(championStatsHistory.AvgKHistory, championStats.AvgK)
		championStatsHistory.AvgDHistory = append(championStatsHistory.AvgDHistory, championStats.AvgD)
		championStatsHistory.AvgAHistory = append(championStatsHistory.AvgAHistory, championStats.AvgA)
		championStatsHistory.StdDevKHistory = append(championStatsHistory.StdDevKHistory, championStats.StdDevK)
		championStatsHistory.StdDevDHistory = append(championStatsHistory.StdDevDHistory, championStats.StdDevD)
		championStatsHistory.StdDevAHistory = append(championStatsHistory.StdDevAHistory, championStats.StdDevA)

		for role, stats := range championStats.StatsPerRole {
			currentRoleStatsHistory := championStatsHistory.HistoryPeRrole[role]

			if stats.SampleSize == 0 {
				s.log.Infof("ChampionStatsHistoryByID for Champion ID %s, Tier %s: Skipping Role %s, no data", champID, tier, role)
				continue
			}

			currentRoleStatsHistory.Versions = append(currentRoleStatsHistory.Versions, gameVersion)
			currentRoleStatsHistory.WinRateHistory = append(currentRoleStatsHistory.WinRateHistory, stats.WinRate)

			currentRoleStatsHistory.AvgKHistory = append(currentRoleStatsHistory.AvgKHistory, stats.AvgK)
			currentRoleStatsHistory.AvgDHistory = append(currentRoleStatsHistory.AvgDHistory, stats.AvgD)
			currentRoleStatsHistory.AvgAHistory = append(currentRoleStatsHistory.AvgAHistory, stats.AvgA)
			currentRoleStatsHistory.StdDevKHistory = append(currentRoleStatsHistory.StdDevKHistory, stats.StdDevK)
			currentRoleStatsHistory.StdDevDHistory = append(currentRoleStatsHistory.StdDevDHistory, stats.StdDevD)
			currentRoleStatsHistory.StdDevAHistory = append(currentRoleStatsHistory.StdDevAHistory, stats.StdDevA)

			championStatsHistory.HistoryPeRrole[role] = currentRoleStatsHistory
		}

		championStatsHistory.ChampionID = championStats.ChampionID
		championStatsHistory.ChampionName = championStats.ChampionName
		championStatsHistory.ChampionRealID = championStats.ChampionRealID
	}

	championStatsHistory.Tier = tier
	championStatsHistory.Timestamp = time.Now()

	out, err := json.Marshal(championStatsHistory)
	if err != nil {
		s.log.Errorf("Error in championStatsHistoryByID with request %s: %s", r.URL.String(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
