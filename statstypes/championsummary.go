package statstypes

import "time"

type ChampionStatsSummary struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`

	GameVersion string `json:"gameversion"`
	Tier        string `json:"tier"`
	// Queue is the Queue the analysis takes into account, e.g., ALL, NORMAL_DRAFT, NORMAL_BLIND, RANKED_SOLO, RANKED_FLEX, ARAM
	Queue string `json:"queue"`

	Timestamp time.Time `json:"timestamp"`

	SampleSize uint64 `json:"samplesize"`

	WinRate float64 `json:"winrate"`

	AvgK float64 `json:"averagekills"`
	AvgD float64 `json:"averagedeaths"`
	AvgA float64 `json:"averageassists"`

	BanRate  float64 `json:"banrate"`
	PickRate float64 `json:"pickrate"`

	Roles []string `json:"roles"`
}
