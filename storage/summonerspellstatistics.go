package storage

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

type SingleSummonerSpellsStatsValues struct {
	SampleSize uint64 `json:"samplesize"`

	SummonerSpells []riotclient.SummonerSpell `json:"summonerspells"`

	PickRate float64 `json:"pickrate"`
	WinRate  float64 `json:"winrate"`
}

type SummonerSpellsStatsValues map[string]SingleSummonerSpellsStatsValues

type SummonerSpellsStats struct {
	ChampionID               uint64 `json:"championid"`
	ChampionRealID           string `json:"championrealid"`
	ChampionName             string `json:"championname"`
	GameVersion              string `json:"gameversion"`
	TotalGamesForGameVersion uint64 `json:"totalgamesforgameversion"`

	Tier string `json:"tier"`
	// Queue is the Queue the analysis takes into account, e.g., ALL, NORMAL_DRAFT, NORMAL_BLIND, RANKED_SOLO, RANKED_FLEX, ARAM
	Queue string `json:"queue"`

	SampleSize uint64 `json:"samplesize"`

	Timestamp time.Time `json:"timestamp"`

	Stats SummonerSpellsStatsValues `json:"stats"`

	StatsPerRole map[string]SummonerSpellsStatsValues `json:"statsperrole"`
}

type SummonerSpellsStatsStorage struct {
	SummonerSpellsStats SummonerSpellsStats `json:"summonerspellsstats"`

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

// GetSummonerSpellsStatsByIDGameVersionTierQueue returns the Champion Summoner Spells stats for a certain game version
func (s *Storage) GetSummonerSpellsStatsByIDGameVersionTierQueue(championID, gameVersion, tier, queue string) (*SummonerSpellsStats, error) {
	returnStats, err := s.backend.GetSummonerSpellsStatsByChampionIDGameVersionTierQueue(championID, gameVersion, tier, queue)
	if err != nil {
		s.log.Warnln("Could not get SummonerSpellsStats data from Storage Backend:", err)
		return nil, err
	}

	return &returnStats.SummonerSpellsStats, nil
}

// StoreSummonerSpellsStats stores the Champion Summoner Spells stats for a certain game version, tier and queue
func (s *Storage) StoreSummonerSpellsStats(stats *SummonerSpellsStats) error {
	key := fmt.Sprintf("%d", stats.ChampionID)

	statsStorage := SummonerSpellsStatsStorage{
		SummonerSpellsStats: *stats,

		ChampionID:   stats.ChampionRealID,
		ChampionKey:  key,
		ChampionName: stats.ChampionName,
		GameVersion:  stats.GameVersion,

		Tier:  stats.Tier,
		Queue: stats.Queue,

		SampleSize: stats.SampleSize,

		TimeStamp: time.Now(),
	}

	return s.backend.StoreSummonerSpellsStats(&statsStorage)
}
