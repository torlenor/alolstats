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

type ChampionStats struct {
	ChampionID     uint64 `json:"championid"`
	ChampionRealID string `json:"championrealid"`
	ChampionName   string `json:"championname"`
	GameVersion    string `json:"gameversion"`

	Timestamp time.Time `json:"timestamp"`

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

	WinLossRatio float64 `json:"winlossratio"`
	WinRate      float64 `json:"winrate"`

	Roles []string `json:"roles"`

	LaneRolePercentage []LaneRolePercentage `json:"lanerolepercentage"`

	LaneRolePercentagePlotly []LaneRolePercentagePlotly `json:"lanerolepercentageplotly"`
}

type ChampionStatsStorage struct {
	ChampionStats ChampionStats `json:"championstats"`

	ChampionID   string `json:"championid"`
	ChampionKey  string `json:"championkey"`
	ChampionName string `json:"championname"`
	GameVersion  string `json:"gameversion"`

	SampleSize uint64 `json:"samplesize"`

	TimeStamp time.Time `json:"timestamp"`
}

// GetChampionStatsByIDGameVersion returns the Champion stats for a certain game version
func (s *Storage) GetChampionStatsByIDGameVersion(championID string, gameVersion string) (*ChampionStats, error) {
	s.log.Debugf("GetChampionStatsByIDGameVersion()")

	stats, err := s.backend.GetChampionStatsByChampionIDGameVersion(championID, gameVersion)
	if err != nil {
		s.log.Warnln("Could not get data from Storage Backend:", err)
		return nil, err
	}

	s.log.Debugf("Returned Champion Stats from Storage")
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
