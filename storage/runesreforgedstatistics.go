package storage

import (
	"fmt"
	"time"
)

type Rune struct {
	ID int
}

type RunesReforgedPicks struct {
	SlotPrimary struct {
		Rune

		Rune0 Rune
		Rune1 Rune
		Rune2 Rune
		Rune3 Rune
	} `json:"slotprimary"`

	SlotSecondary struct {
		Rune

		Rune0 Rune
		Rune1 Rune
	} `json:"slotsecondary"`

	StatPerks struct {
		Perk0 Rune
		Perk1 Rune
		Perk2 Rune
	} `json:"statperks"`
}

// SingleRunesReforgedStatsValues is one stat entry for a given
// Runes Reforged combination identified by Hash.
type SingleRunesReforgedStatsValues struct {
	SampleSize uint64 `json:"samplesize"`
	Hash       string `json:"hash"`

	RunesReforged RunesReforgedPicks `json:"runesreforged"`

	PickRate float64 `json:"pickrate"`
	WinRate  float64 `json:"winrate"`
}

// RunesReforgedStatsValues holds a set of different Runes Reforged
// combinations.
type RunesReforgedStatsValues []SingleRunesReforgedStatsValues

// RunesReforgedStats holds one particular stat entry for
// the given champion, game version, tier and queue.
type RunesReforgedStats struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`
	GameVersion    string `json:"gameversion"`

	Tier string `json:"tier"`
	// Queue is the Queue the analysis takes into account, e.g., ALL, NORMAL_DRAFT, NORMAL_BLIND, RANKED_SOLO, RANKED_FLEX, ARAM
	Queue string `json:"queue"`

	SampleSize uint64 `json:"samplesize"`

	Timestamp time.Time `json:"timestamp"`

	RunesReforgedStatsValues

	StatsPerRole map[string]RunesReforgedStatsValues `json:"statsperrole"`
}

// RunesReforgedStatsStorage is the struct which is used to store/retreive
// data to/from the storage backend.
//
// It contains all the necessary fields to distinguishe this special
// stat from the other ones.
type RunesReforgedStatsStorage struct {
	RunesReforgedStats RunesReforgedStats `json:"runesreforgedstats"`

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

// GetRunesReforgedStatsByIDGameVersionTierQueue returns the Champion Runes Reforged stats for a certain game version
func (s *Storage) GetRunesReforgedStatsByIDGameVersionTierQueue(championID, gameVersion, tier, queue string) (*RunesReforgedStats, error) {
	returnStats, err := s.backend.GetRunesReforgedStatsByChampionIDGameVersionTierQueue(championID, gameVersion, tier, queue)
	if err != nil {
		s.log.Warnln("Could not get RunesReforgedStats data from Storage Backend:", err)
		return nil, err
	}

	return &returnStats.RunesReforgedStats, nil
}

// StoreRunesReforgedStats stores the Champion Runes Reforged stats for a certain game version, tier and queue
func (s *Storage) StoreRunesReforgedStats(stats *RunesReforgedStats) error {
	key := fmt.Sprintf("%d", stats.ChampionID)

	statsStorage := RunesReforgedStatsStorage{
		RunesReforgedStats: *stats,

		ChampionID:   stats.ChampionRealID,
		ChampionKey:  key,
		ChampionName: stats.ChampionName,
		GameVersion:  stats.GameVersion,

		Tier:  stats.Tier,
		Queue: stats.Queue,

		SampleSize: stats.SampleSize,

		TimeStamp: time.Now(),
	}

	return s.backend.StoreRunesReforgedStats(&statsStorage)
}
