package statsrunner

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"git.abyle.org/hps/alolstats/riotclient"
	"git.abyle.org/hps/alolstats/statsrunner/analyzer"
	"git.abyle.org/hps/alolstats/storage"
	"git.abyle.org/hps/alolstats/utils"
)

func (sr *StatsRunner) runesReforgedWorker() {
	defer sr.workersWG.Done()

	var nextUpdate time.Duration

	queueIDtoQueue := map[uint64]string{
		400: "NORMAL_DRAFT",
		420: "RANKED_SOLO",
		430: "NORMAL_BLIND",
		440: "RANKED_FLEX",
	}

	for {
		select {
		case <-sr.stopWorkers:
			sr.log.Printf("Stopping runesReforgedWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			sr.calculationMutex.Lock()
			sr.log.Infof("Performing runesReforgedWorker run")
			start := time.Now()

			mapID := uint64(11)

			for queueID, queue := range queueIDtoQueue {
				for _, versionStr := range sr.config.GameVersion {
					sr.shouldWorkersStopMutex.RLock()
					if sr.shouldWorkersStop {
						return
					}
					sr.shouldWorkersStopMutex.RUnlock()
					version, err := utils.SplitNumericVersion(versionStr)
					if err != nil {
						sr.log.Errorf("Could not get game determine requested game version: %s", err)
						continue
					}
					gameVersion := fmt.Sprintf("%d.%d", version[0], version[1])
					majorMinor := fmt.Sprintf("%d\\.%d\\.", version[0], version[1])

					analyzer := analyzer.NewRunesReforgedAnalyzer(int(version[0]), int(version[1]))

					sr.log.Infof("Calculation of runesReforgedWorker for Game Version %s and Queue %s started", gameVersion, queue)

					cur, err := sr.storage.GetMatchesCursorByGameVersionMapQueueID(majorMinor, mapID, queueID)
					if err != nil {
						sr.log.Errorf("Error performing runesReforgedWorker calculation for Game Version %s: %s", gameVersion, err)
						continue
					}
					currentMatch := &riotclient.MatchDTO{}
					cnt := 0
					for cur.Next() {
						sr.shouldWorkersStopMutex.RLock()
						if sr.shouldWorkersStop {
							return
						}
						sr.shouldWorkersStopMutex.RUnlock()
						err := cur.Decode(currentMatch)
						if err != nil {
							sr.log.Errorf("Error decoding match: %s", err)
							continue
						}

						analyzer.FeedMatch(currentMatch)
						cnt++
					}
					cur.Close()

					result := analyzer.Analyze()

					// Prepare results for ItemStats (ALL tiers)
					for _, runesReforgedCombiStats := range result {
						stats, err := sr.prepareRunesReforgedStats(runesReforgedCombiStats, queue, "ALL")
						if err == nil {
							err = sr.storage.StoreRunesReforgedStats(stats)
							if err != nil {
								sr.log.Warnf("Something went wrong storing the Champion Runes Reforged Stats: %s", err)
							}
						}
					}

					sr.log.Infof("Calculation of runesReforgedWorker for Game Version %s and Queue %s done. Analyzed %d matches", gameVersion, queue, cnt)
				}
			}

			gameVersions := storage.GameVersions{}
			for _, val := range sr.config.GameVersion {
				ver, err := utils.SplitNumericVersion(val)
				if err != nil {
					continue
				}

				verStr := fmt.Sprintf("%d.%d", ver[0], ver[1])
				gameVersions.Versions = append(gameVersions.Versions, verStr)
			}
			sr.storage.StoreKnownGameVersions(&gameVersions)

			nextUpdate = time.Minute * time.Duration(sr.config.ChampionsStats.UpdateInverval)
			elapsed := time.Since(start)
			sr.log.Infof("Finished runesReforgedWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
			sr.calculationMutex.Unlock()
		}
	}
}

func (sr *StatsRunner) prepareRunesReforgedStatsValues(runesReforgedCombiStats analyzer.RunesReforgedCombiStatistics, totalSampleSize uint32) storage.RunesReforgedStatsValues {
	runesReforgedStatsBySampleSize := make(map[float64]storage.SingleRunesReforgedStatsValues)

	var runesReforgedStatsValues storage.RunesReforgedStatsValues
	for runesReforgedCombination, itemCounts := range runesReforgedCombiStats {
		if itemCounts.Picks > 0 {
			is := storage.SingleRunesReforgedStatsValues{}
			is.WinRate = float64(itemCounts.Wins) / float64(itemCounts.Picks)
			is.PickRate = float64(itemCounts.Picks) / float64(totalSampleSize)
			is.SampleSize = uint64(itemCounts.Picks)
			is.Hash = runesReforgedCombination
			is.RunesReforged = itemCounts.RunesReforged
			runesReforgedStatsValues = append(runesReforgedStatsValues, is)
			runesReforgedStatsBySampleSize[is.PickRate] = is
		}
	}

	if sr.config.ItemsStats.KeepOnlyHighestPickRate {
		var keys []float64
		for k := range runesReforgedStatsBySampleSize {
			keys = append(keys, k)
		}
		sort.Float64s(keys)

		var highestRunesReforgedStatsValue storage.RunesReforgedStatsValues

		cnt := uint32(0)
		for i := len(keys) - 1; i >= 0; i-- {
			highestRunesReforgedStatsValue = append(highestRunesReforgedStatsValue, runesReforgedStatsBySampleSize[keys[i]])
			cnt++
			if cnt > (sr.config.ItemsStats.KeepOnlyNHighest - 1) {
				break
			}
		}
		return highestRunesReforgedStatsValue
	}

	return runesReforgedStatsValues
}

func (sr *StatsRunner) prepareRunesReforgedStats(stats *analyzer.ChampionRunesReforgedCombiStatistics, queue string, tier string) (*storage.RunesReforgedStats, error) {
	if stats.TotalSampleSize == 0 {
		return nil, fmt.Errorf("No data")
	}

	gameVersion := fmt.Sprintf("%d.%d", stats.GameVersionMajor, stats.GameVersionMinor)

	runesReforgedStats := storage.RunesReforgedStats{}
	runesReforgedStats.ChampionID = uint64(stats.ChampionID)
	runesReforgedStats.GameVersion = gameVersion

	runesReforgedStats.RunesReforgedStatsValues = sr.prepareRunesReforgedStatsValues(stats.Total, stats.TotalSampleSize)

	runesReforgedStats.StatsPerRole = make(map[string]storage.RunesReforgedStatsValues)
	for role, statValues := range stats.PerRole {
		if roleSampleSize, ok := stats.PerRoleSampleSize[role]; ok {
			runesReforgedStats.StatsPerRole[role] = sr.prepareRunesReforgedStatsValues(statValues, roleSampleSize)
		} else {
			sr.log.Warnf("Bug! There is no PerRoleSampleSize for role %s", role)
		}
	}

	champions := sr.storage.GetChampions(false)
	for _, val := range champions {
		if val.Key == strconv.FormatUint(uint64(stats.ChampionID), 10) {
			runesReforgedStats.ChampionName = val.Name
			runesReforgedStats.ChampionRealID = val.ID
			break
		}
	}

	runesReforgedStats.Queue = queue
	runesReforgedStats.Tier = tier

	runesReforgedStats.SampleSize = uint64(stats.TotalSampleSize)

	runesReforgedStats.Timestamp = time.Now()

	return &runesReforgedStats, nil
}
