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

func (sr *StatsRunner) itemWinRateWorker() {
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
			sr.log.Printf("Stopping itemWinRateWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			sr.calculationMutex.Lock()
			sr.log.Infof("Performing itemWinRateWorker run")
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

					analyzer := analyzer.NewItemAnalyzer(int(version[0]), int(version[1]))

					sr.log.Infof("Calculation of itemWinRateWorker for Game Version %s and Queue %s started", gameVersion, queue)

					cur, err := sr.storage.GetMatchesCursorByGameVersionMapQueueID(majorMinor, mapID, queueID)
					if err != nil {
						sr.log.Errorf("Error performing itemWinRateWorker calculation for Game Version %s: %s", gameVersion, err)
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
					for _, itemCombiStats := range result {
						stats, err := sr.prepareItemStats(itemCombiStats, queue, "ALL")
						if err == nil {
							err = sr.storage.StoreItemStats(stats)
							if err != nil {
								sr.log.Warnf("Something went wrong storing the Champion Item Stats: %s", err)
							}
						}
					}

					sr.log.Infof("Calculation of itemWinRateWorker for Game Version %s and Queue %s done. Analyzed %d matches", gameVersion, queue, cnt)
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
			sr.log.Infof("Finished itemWinRateWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
			sr.calculationMutex.Unlock()
		}
	}
}

func (sr *StatsRunner) prepareItemStats(stats *analyzer.ChampionItemCombiStatistics, queue string, tier string) (*storage.ItemStats, error) {
	if stats.TotalSampleSize == 0 {
		return nil, fmt.Errorf("No data")
	}

	gameVersion := fmt.Sprintf("%d.%d", stats.GameVersionMajor, stats.GameVersionMinor)

	itemStats := storage.ItemStats{}
	itemStats.ChampionID = uint64(stats.ChampionID)
	itemStats.GameVersion = gameVersion

	itemStats.ItemStatsValues = make(storage.ItemStatsValues)

	itemStatsBySampleSize := make(map[float64]storage.SingleItemStatsValues)

	for itemCombination, itemCounts := range stats.Total {
		if itemCounts.Picks > 0 {
			is := itemStats.ItemStatsValues[itemCombination]
			is.WinRate = float64(itemCounts.Wins) / float64(itemCounts.Picks)
			is.PickRate = float64(itemCounts.Picks) / float64(stats.TotalSampleSize)
			is.SampleSize = uint64(itemCounts.Picks)
			is.ItemHash = itemCombination
			itemStats.ItemStatsValues[itemCombination] = is
			itemStatsBySampleSize[is.PickRate] = is
		}
	}

	var keys []float64
	for k := range itemStatsBySampleSize {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	// sr.log.Infof("Highest pick rate (%f) for Champ %d was item combination %s with a Win Rate of %f percent", itemStatsBySampleSize[keys[len(keys)-1]].PickRate, stats.ChampionID, itemStatsBySampleSize[keys[len(keys)-1]].ItemHash, 100*itemStatsBySampleSize[keys[len(keys)-1]].WinRate)

	champions := sr.storage.GetChampions(false)
	for _, val := range champions {
		if val.Key == strconv.FormatUint(uint64(stats.ChampionID), 10) {
			itemStats.ChampionName = val.Name
			itemStats.ChampionRealID = val.ID
			break
		}
	}

	itemStats.Queue = queue
	itemStats.Tier = tier

	itemStats.SampleSize = uint64(stats.TotalSampleSize)

	itemStats.Timestamp = time.Now()

	return &itemStats, nil
}
