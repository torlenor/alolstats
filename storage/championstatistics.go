package storage

import (
	"fmt"
	"time"

	"git.abyle.org/hps/alolstats/statstypes"
)

type ChampionStatsStorage struct {
	ChampionStats statstypes.ChampionStats `json:"championstats"`

	ChampionID   string `json:"championid"`
	ChampionKey  string `json:"championkey"`
	ChampionName string `json:"championname"`
	GameVersion  string `json:"gameversion"`

	Tier string `json:"tier"`
	// Queue is the Queue the analysis takes into account, e.g., ALL, NORMAL_DRAFT, NORMAL_BLIND, RANKED_SOLO, RANKED_FLEX, ARAM
	Queue string `json:"queue"`

	SampleSize uint64 `json:"samplesize"`

	TimeStamp time.Time `json:"timestamp"`
}

type ChampionStatsSummaryStorage struct {
	ChampionsStatsSummary []statstypes.ChampionStatsSummary `json:"championstatssummary"`

	GameVersion string `json:"gameversion"`
	Tier        string `json:"tier"`
	// Queue is the Queue the analysis takes into account, e.g., ALL, NORMAL_DRAFT, NORMAL_BLIND, RANKED_SOLO, RANKED_FLEX, ARAM
	Queue string `json:"queue"`
}

// GetChampionStatsByIDGameVersionTierQueue returns the Champion stats for a certain game version, tier and queue
func (s *Storage) GetChampionStatsByIDGameVersionTierQueue(championID string, gameVersion string, tier string, queue string) (*statstypes.ChampionStats, error) {
	stats, err := s.backend.GetChampionStatsByChampionIDGameVersionTierQueue(championID, gameVersion, tier, queue)
	if err != nil {
		s.log.Warnln("Could not get data from Storage Backend:", err)
		return nil, err
	}

	returnStats := &stats.ChampionStats
	return returnStats, nil
}

// StoreChampionStats stores the Champion stats for a certain game version
func (s *Storage) StoreChampionStats(stats *statstypes.ChampionStats) error {
	key := fmt.Sprintf("%d", stats.ChampionID)

	championstatsStorage := ChampionStatsStorage{
		ChampionStats: *stats,

		ChampionID:   stats.ChampionRealID,
		ChampionKey:  key,
		ChampionName: stats.ChampionName,
		GameVersion:  stats.GameVersion,

		Tier:  stats.Tier,
		Queue: stats.Queue,

		SampleSize: stats.SampleSize,

		TimeStamp: time.Now(),
	}

	s.backend.StoreChampionStats(&championstatsStorage)
	return nil
}

// StoreKnownGameVersions stores a new list of known game versions
func (s *Storage) StoreKnownGameVersions(gameVersions *GameVersions) error {
	return s.backend.StoreKnownGameVersions(gameVersions)
}

// GetKnownGameVersions retrieves a list of known game versions
func (s *Storage) GetKnownGameVersions() (*GameVersions, error) {
	return s.backend.GetKnownGameVersions()
}

// GetChampionStatsSummaryByGameVersionTierQueue returns the Champion stats for a certain game version, tier and queue
func (s *Storage) GetChampionStatsSummaryByGameVersionTierQueue(gameVersion, tier, queue string) (*ChampionStatsSummaryStorage, error) {
	stats, err := s.backend.GetChampionStatsSummaryByGameVersionTierQueue(gameVersion, tier, queue)
	if err != nil {
		s.log.Warnln("Could not get data from Storage Backend:", err)
		return nil, err
	}

	return stats, nil
}

// StoreChampionStatsSummary stores the Champion stats Summary
func (s *Storage) StoreChampionStatsSummary(statsSummary *ChampionStatsSummaryStorage) error {

	s.backend.StoreChampionStatsSummary(statsSummary)
	return nil
}
