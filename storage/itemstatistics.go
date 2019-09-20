package storage

import (
	"fmt"
	"time"
)

type SingleItemStatsValues struct {
	SampleSize uint64 `json:"samplesize"`
	ItemHash   string `json:"itemHash"`

	Items []int `json:"items"`

	PickRate float64 `json:"pickrate"`
	WinRate  float64 `json:"winrate"`
}

type ItemStatsValues []SingleItemStatsValues

type ItemStats struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`
	GameVersion    string `json:"gameversion"`

	Tier string `json:"tier"`
	// Queue is the Queue the analysis takes into account, e.g., ALL, NORMAL_DRAFT, NORMAL_BLIND, RANKED_SOLO, RANKED_FLEX, ARAM
	Queue string `json:"queue"`

	Timestamp time.Time `json:"timestamp"`

	ItemStatsValues

	StatsPerRole map[string]ItemStatsValues `json:"statsperrole"`

	SampleSize uint64 `json:"samplesize"`
}

type ItemStatsStorage struct {
	ItemStats ItemStats `json:"itemstats"`

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

// GetItemStatsByIDGameVersion returns the Champion Item stats for a certain game version
func (s *Storage) GetItemStatsByIDGameVersion(championID string, gameVersion string, tier string, queue string) (*ItemStats, error) {
	stats, err := s.backend.GetItemStatsByChampionIDGameVersionTierQueue(championID, gameVersion, tier, queue)
	if err != nil {
		s.log.Warnln("Could not get data from Storage Backend:", err)
		return nil, err
	}

	returnStats := &stats.ItemStats
	return returnStats, nil
}

// StoreItemStats stores the Champion Item stats for a certain game version
func (s *Storage) StoreItemStats(stats *ItemStats) error {
	key := fmt.Sprintf("%d", stats.ChampionID)

	statsStorage := ItemStatsStorage{
		ItemStats: *stats,

		ChampionID:   stats.ChampionRealID,
		ChampionKey:  key,
		ChampionName: stats.ChampionName,
		GameVersion:  stats.GameVersion,

		Tier:  stats.Tier,
		Queue: stats.Queue,

		SampleSize: stats.SampleSize,

		TimeStamp: time.Now(),
	}

	return s.backend.StoreItemStats(&statsStorage)
}
