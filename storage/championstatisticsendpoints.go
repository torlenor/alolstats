package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"git.abyle.org/hps/alolstats/statstypes"
	"git.abyle.org/hps/alolstats/utils"
)

func (s *Storage) championStatsByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API championByID request from", r.RemoteAddr)

	id, err := extractURLStringParameter(r.URL.Query(), "id")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	gameVersion, err := extractURLStringParameter(r.URL.Query(), "gameversion")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	queue, err := extractURLStringParameter(r.URL.Query(), "queue")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	tier, err := extractURLStringParameter(r.URL.Query(), "tier")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	championStats, err := s.GetChampionStatsByIDGameVersionTierQueue(id, gameVersion, tier, queue)
	if err != nil {
		s.log.Errorln(err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, fmt.Sprintf("Could not get stats for champion ID %s", id)), http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(championStats)
	if err != nil {
		s.log.Errorf("Error in championByID with request %s: %s", r.URL.String(), err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) championStats(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API ChampionsStats request from", r.RemoteAddr)
	var gameVersion string
	var tier string
	var queue string

	gameVersion, err := extractURLStringParameter(r.URL.Query(), "gameversion")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	queue, err = extractURLStringParameter(r.URL.Query(), "queue")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	tier, err = extractURLStringParameter(r.URL.Query(), "tier")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	statsSummary, err := s.GetChampionStatsSummaryByGameVersionTierQueue(gameVersion, tier, queue)
	if err != nil {
		s.log.Errorf("Error in ChampionsStats with request %s: %s", r.URL.String(), err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest,
			fmt.Sprintf("Could not get stats for champion stats summary for gameversion %s, tier %, queue %s",
				gameVersion, tier, queue)),
			http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(statsSummary.ChampionsStatsSummary)
	if err != nil {
		s.log.Errorf("Error in ChampionsStats with request %s: %s", r.URL.String(), err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Problem converting Champion to JSON")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}

func (s *Storage) championStatsHistoryByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API championStatsHistoryByID request from", r.RemoteAddr)
	var champID string
	var tier string
	var queue string

	champID, err := extractURLStringParameter(r.URL.Query(), "id")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	queue, err = extractURLStringParameter(r.URL.Query(), "queue")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	tier, err = extractURLStringParameter(r.URL.Query(), "tier")
	if err != nil {
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest, err.Error()), http.StatusBadRequest)
		return
	}

	gameVersions, err := s.GetKnownGameVersions()
	if err != nil {
		s.log.Errorf("Error in Champion Stats History request %s: %s", r.URL.String(), err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error")), http.StatusInternalServerError)
		return
	}

	championStatsHistory := statstypes.ChampionStatsHistory{}
	championStatsHistory.HistoryPeRrole = make(map[string]statstypes.ChampionStatsPerRoleSingleHistory)

	hasHistory := false
	for _, gameVersion := range gameVersions.Versions {
		championStats, err := s.GetChampionStatsByIDGameVersionTierQueue(champID, gameVersion, tier, queue)
		if err != nil {
			continue
		}

		if championStats.SampleSize == 0 {
			s.log.Debugf("ChampionStatsHistoryByID for Champion ID %s, Tier %s and Queue %s: Skipping Game Version %s, no data", champID, tier, queue, gameVersion)
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
				s.log.Debugf("ChampionStatsHistoryByID for Champion ID %s, Tier %s and Queue %s: Skipping Role %s, no data", champID, tier, queue, role)
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
		hasHistory = true
	}

	championStatsHistory.Tier = tier
	championStatsHistory.Timestamp = time.Now()

	if !hasHistory {
		s.log.Errorf("Error in championStatsHistoryByID with request %s: No History Data available", r.URL.String())
		http.Error(w, utils.GenerateStatusResponse(http.StatusBadRequest,
			fmt.Sprintf("No History Data available for gameversion %s, tier %s, queue %s",
				gameVersions, tier, queue)),
			http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(championStatsHistory)
	if err != nil {
		s.log.Errorf("Error in Champion Stats History request %s: %s", r.URL.String(), err)
		http.Error(w, utils.GenerateStatusResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", s.getHTTPGetResponseHeader("Cache-Control"))
	io.WriteString(w, string(out))

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
