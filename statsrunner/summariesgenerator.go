package statsrunner

import (
	"time"

	"github.com/torlenor/alolstats/storage"
)

func (sr *StatsRunner) generateChampionsSummary(gameVersion string, league string) (*storage.ChampionStatsSummaryStorage, error) {
	var championsStatsSummary storage.ChampionStatsSummaryStorage

	champions := sr.storage.GetChampions(false)
	for _, champ := range champions {
		championStats, err := sr.storage.GetChampionStatsByIDGameVersionTier(champ.ID, gameVersion, league)
		if err != nil {
			continue
		}
		if championStats.SampleSize > 0 {
			summary := storage.ChampionStatsSummary{}

			summary.ChampionID = championStats.ChampionID
			summary.ChampionName = championStats.ChampionName
			summary.ChampionRealID = championStats.ChampionRealID
			summary.AvgK = championStats.AvgK
			summary.AvgA = championStats.AvgA
			summary.AvgD = championStats.AvgD
			summary.WinRate = championStats.WinRate
			summary.PickRate = championStats.PickRate
			summary.BanRate = championStats.BanRate
			summary.Roles = championStats.Roles
			summary.SampleSize = championStats.SampleSize

			summary.GameVersion = gameVersion
			summary.Tier = league
			summary.Timestamp = time.Now()

			championsStatsSummary.ChampionsStatsSummary = append(championsStatsSummary.ChampionsStatsSummary, summary)
		}
	}

	championsStatsSummary.GameVersion = gameVersion
	championsStatsSummary.Tier = league

	return &championsStatsSummary, nil
}

func (sr *StatsRunner) generatePatchHistories() {

}
