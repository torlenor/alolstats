package statstypes

import "time"

type ChampionStatsSingleHistory struct {
	Versions        []string  `json:"versions"`
	WinRateHistory  []float64 `json:"winRateHistory"`
	PickRateHistory []float64 `json:"pickRateHistory"`
	BanRateHistory  []float64 `json:"banRateHistory"`

	AvgKHistory    []float64 `json:"averagekillsHistory"`
	StdDevKHistory []float64 `json:"stddevkillsHistory"`
	AvgDHistory    []float64 `json:"averagedeathsHistory"`
	StdDevDHistory []float64 `json:"stddevdeathsHistory"`
	AvgAHistory    []float64 `json:"averageassistsHistory"`
	StdDevAHistory []float64 `json:"stddevassistsHistory"`
}

type ChampionStatsPerRoleSingleHistory struct {
	Versions       []string  `json:"versions"`
	WinRateHistory []float64 `json:"winRateHistory"`

	AvgKHistory    []float64 `json:"averagekillsHistory"`
	StdDevKHistory []float64 `json:"stddevkillsHistory"`
	AvgDHistory    []float64 `json:"averagedeathsHistory"`
	StdDevDHistory []float64 `json:"stddevdeathsHistory"`
	AvgAHistory    []float64 `json:"averageassistsHistory"`
	StdDevAHistory []float64 `json:"stddevassistsHistory"`
}

type ChampionStatsHistory struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`

	Tier string `json:"tier"`
	// Queue is the Queue the analysis takes into account, e.g., ALL, NORMAL_DRAFT, NORMAL_BLIND, RANKED_SOLO, RANKED_FLEX, ARAM
	Queue string `json:"queue"`

	Timestamp time.Time `json:"timestamp"`

	ChampionStatsSingleHistory

	HistoryPeRrole map[string]ChampionStatsPerRoleSingleHistory `json:"historyperrole"`
}
