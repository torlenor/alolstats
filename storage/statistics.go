package storage

import (
	"fmt"
	"time"
)

type LaneRolePercentage struct {
	Lane string `json:"lane"`
	Role string `json:"role"`

	Percentage float64 `json:"percentage"`
	Wins       uint32  `json:"wins"`
	NGames     uint32  `json:"ngames"`
}

type LaneRolePercentagePlotly struct {
	X []string  `json:"x"` // ['TOP', 'MIDDLE', 'JUNGLE', 'BOT', 'UNKNOWN'],
	Y []float64 `json:"y"` // [2.2058823529411766, 2.941176470588235, 0.7352941176470588, 0, 0],

	Name string `json:"name"` // 'Solo',
	Type string `json:"type"` // 'bar'
}

type StatsValues struct {
	SampleSize uint64 `json:"samplesize"`

	AvgK    float64 `json:"averagekills"`
	StdDevK float64 `json:"stddevkills"`
	MedianK float64 `json:"mediankills"`

	AvgD    float64 `json:"averagedeaths"`
	StdDevD float64 `json:"stddevdeaths"`
	MedianD float64 `json:"mediandeaths"`

	AvgA    float64 `json:"averageassists"`
	StdDevA float64 `json:"stddevassists"`
	MedianA float64 `json:"medianassists"`

	AvgGoldEarned                     float64 `json:"average_goldearned"`
	AvgTotalMinionsKilled             float64 `json:"average_totalminionskilled"`
	AvgTotalHeal                      float64 `json:"average_totalheal"`
	AvgTotalDamageDealt               float64 `json:"average_totaldamagedealt"`
	AvgTotalDamageDealtToChampions    float64 `json:"average_totaldamagedealttochampions"`
	AvgTotalDamageTaken               float64 `json:"average_totaldamagetaken"`
	AvgMagicDamageDealt               float64 `json:"average_magicdamagedealt"`
	AvgMagicDamageDealtToChampions    float64 `json:"average_magicdamagedealttochampions"`
	AvgPhysicalDamageDealt            float64 `json:"average_physicaldamagedealt"`
	AvgPhysicalDamageDealtToChampions float64 `json:"average_physicaldamagedealttochampions"`
	AvgPhysicalDamageTaken            float64 `json:"average_physicaldamagetaken"`
	AvgTrueDamageDealt                float64 `json:"average_truedamagedealt"`
	AvgTrueDamageDealtToChampions     float64 `json:"average_truedamagedealttochampions"`
	AvgTrueDamageTaken                float64 `json:"average_truedamagetaken"`

	StdDevGoldEarned                     float64 `json:"stddev_goldearned"`
	StdDevTotalMinionsKilled             float64 `json:"stddev_totalminionskilled"`
	StdDevTotalHeal                      float64 `json:"stddev_totalheal"`
	StdDevTotalDamageDealt               float64 `json:"stddev_totaldamagedealt"`
	StdDevTotalDamageDealtToChampions    float64 `json:"stddev_totaldamagedealttochampions"`
	StdDevTotalDamageTaken               float64 `json:"stddev_totaldamagetaken"`
	StdDevMagicDamageDealt               float64 `json:"stddev_magicdamagedealt"`
	StdDevMagicDamageDealtToChampions    float64 `json:"stddev_magicdamagedealttochampions"`
	StdDevPhysicalDamageDealt            float64 `json:"stddev_physicaldamagedealt"`
	StdDevPhysicalDamageDealtToChampions float64 `json:"stddev_physicaldamagedealttochampions"`
	StdDevPhysicalDamageTaken            float64 `json:"stddev_physicaldamagetaken"`
	StdDevTrueDamageDealt                float64 `json:"stddev_truedamagedealt"`
	StdDevTrueDamageDealtToChampions     float64 `json:"stddev_truedamagedealttochampions"`
	StdDevTrueDamageTaken                float64 `json:"stddev_truedamagetaken"`

	AvgDamageDealtToObjectives float64 `json:"average_damagedealttoobjectives"`
	AvgDamageDealtToTurrets    float64 `json:"average_damagedealttoturrets"`
	AvgTimeCCingOthers         float64 `json:"average_timeccingothers"`

	StdDevDamageDealtToObjectives float64 `json:"stddev_damagedealttoobjectives"`
	StdDevDamageDealtToTurrets    float64 `json:"stddev_damagedealttoturrets"`
	StdDevTimeCCingOthers         float64 `json:"stddev_timeccingothers"`

	WinLossRatio float64 `json:"winlossratio"`
	WinRate      float64 `json:"winrate"`
}

type ChampionStats struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`
	GameVersion    string `json:"gameversion"`

	Tier string `json:"tier"`

	Timestamp time.Time `json:"timestamp"`

	StatsValues

	BanRate float64 `json:"banrate"`

	PickRate float64 `json:"pickrate"`

	Roles []string `json:"roles"`

	StatsPerRole map[string]StatsValues `json:"statsperrole"`

	LaneRolePercentage []LaneRolePercentage `json:"lanerolepercentage"`

	LaneRolePercentagePlotly []LaneRolePercentagePlotly `json:"lanerolepercentageplotly"`
}

type ChampionStatsStorage struct {
	ChampionStats ChampionStats `json:"championstats"`

	ChampionID   string `json:"championid"`
	ChampionKey  string `json:"championkey"`
	ChampionName string `json:"championname"`
	GameVersion  string `json:"gameversion"`

	Tier string `json:"tier"`

	SampleSize uint64 `json:"samplesize"`

	TimeStamp time.Time `json:"timestamp"`
}

type ChampionStatsSummary struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`

	GameVersion string `json:"gameversion"`
	Tier        string `json:"tier"`

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

type ChampionStatsSummaryStorage struct {
	ChampionsStatsSummary []ChampionStatsSummary `json:"championstatssummary"`

	GameVersion string `json:"gameversion"`
	Tier        string `json:"tier"`
}

// GetChampionStatsByIDGameVersionTier returns the Champion stats for a certain game version
func (s *Storage) GetChampionStatsByIDGameVersionTier(championID string, gameVersion string, tier string) (*ChampionStats, error) {
	stats, err := s.backend.GetChampionStatsByChampionIDGameVersionTier(championID, gameVersion, tier)
	if err != nil {
		s.log.Warnln("Could not get data from Storage Backend:", err)
		return nil, err
	}

	returnStats := &stats.ChampionStats
	return returnStats, nil
}

// StoreChampionStats stores the Champion stats for a certain game version
func (s *Storage) StoreChampionStats(stats *ChampionStats) error {
	key := fmt.Sprintf("%d", stats.ChampionID)

	championstatsStorage := ChampionStatsStorage{
		ChampionStats: *stats,

		ChampionID:   stats.ChampionRealID,
		ChampionKey:  key,
		ChampionName: stats.ChampionName,
		GameVersion:  stats.GameVersion,

		Tier: stats.Tier,

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

// GetChampionStatsSummaryByGameVersionTier returns the Champion stats for a certain game version
func (s *Storage) GetChampionStatsSummaryByGameVersionTier(gameVersion string, tier string) (*ChampionStatsSummaryStorage, error) {
	stats, err := s.backend.GetChampionStatsSummaryByGameVersionTier(gameVersion, tier)
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
