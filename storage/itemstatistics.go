package storage

import (
	"fmt"
	"time"
)

type SingleItemStatsValues struct {
	SampleSize uint64 `json:"samplesize"`
	ItemHash   string `json:"itemHash"`

	PickRate float64 `json:"pickrate"`
	WinRate  float64 `json:"winrate"`
}

type ItemStatsValues map[string]SingleItemStatsValues

type ItemStats struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`
	GameVersion    string `json:"gameversion"`

	Timestamp time.Time `json:"timestamp"`

	ItemStatsValues

	StatsPerRole map[string]ItemStatsValues `json:"statsperrole"`
}

type ItemStatsStorage struct {
	ItemStats ItemStats `json:"itemstats"`

	ChampionID   string `json:"championid"`
	ChampionKey  string `json:"championkey"`
	ChampionName string `json:"championname"`
	GameVersion  string `json:"gameversion"`

	SampleSize uint64 `json:"samplesize"`

	TimeStamp time.Time `json:"timestamp"`
}

// GetItemStatsByIDGameVersion returns the Champion Item stats for a certain game version
func (s *Storage) GetItemStatsByIDGameVersion(championID string, gameVersion string) (*ItemStats, error) {
	stats, err := s.backend.GetItemStatsByChampionIDGameVersion(championID, gameVersion)
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

		TimeStamp: time.Now(),
	}

	s.backend.StoreItemStats(&statsStorage)
	return nil
}
