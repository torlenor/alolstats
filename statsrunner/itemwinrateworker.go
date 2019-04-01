package statsrunner

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/torlenor/alolstats/riotclient"
	"github.com/torlenor/alolstats/storage"
	"github.com/torlenor/alolstats/utils"
)

type singleItemCounter struct {
	Picks uint64
	Wins  uint64
}

type singleItemCounters map[string]singleItemCounter // [itemHash]

type itemCounter struct {
	ChampionID int

	GameVersion string

	TotalPicks uint64

	SingleItemCounters singleItemCounters
	PerRole            map[string]map[string]singleItemCounters // [lane][role]
}

type itemCounters map[int]itemCounter // [id], e.g., 1, 10, 43, ...

func (sr *StatsRunner) newItemsCounters(champions riotclient.ChampionsList, gameVersion string) itemCounters {
	itemCounters := make(itemCounters)
	for _, champ := range champions {
		champID, err := strconv.Atoi(champ.Key)
		if err != nil {
			sr.log.Warnf("Something bad happened: %s", err)
			continue
		}
		itemCounters[champID] = itemCounter{
			ChampionID:  champID,
			GameVersion: gameVersion,

			SingleItemCounters: make(singleItemCounters),
			PerRole:            make(map[string]map[string]singleItemCounters),
		}
	}

	return itemCounters
}

func hashItems(items []int) string {
	sort.Ints(items)
	var s string
	for _, item := range items {
		s = s + strconv.Itoa(item) + "_"
	}

	return strings.TrimSuffix(s, "_")
}

func (sr *StatsRunner) itemWinRateWorker() {
	sr.workersWG.Add(1)
	defer sr.workersWG.Done()

	var nextUpdate time.Duration

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

			champions := sr.storage.GetChampions(false)

			mapID := uint64(11)
			highQueueID := uint64(440)
			lowQueueID := uint64(400)

			for _, versionStr := range sr.config.GameVersion {
				if sr.shouldWorkersStop {
					return
				}

				version, err := utils.SplitNumericVersion(versionStr)
				if err != nil {
					sr.log.Warnf("Something bad happened: %s", err)
					continue
				}
				gameVersion := fmt.Sprintf("%d.%d", version[0], version[1])
				majorMinor := fmt.Sprintf("%d\\.%d\\.", version[0], version[1])
				sr.log.Debugf("itemWinRateWorker calculation for Game Version %s started", gameVersion)

				itemCountersAllTiers := sr.newItemsCounters(champions, gameVersion)

				totalGamesForGameVersion := uint64(0)

				cur, err := sr.storage.GetMatchesCursorByGameVersionMapBetweenQueueIDs(majorMinor, mapID, highQueueID, lowQueueID)
				if err != nil {
					sr.log.Errorf("Error performing itemWinRateWorker calculation for Game Version %s: %s", gameVersion, err)
					continue
				}
				currentMatch := &riotclient.MatchDTO{}
				cnt := 0

				for cur.Next() {
					err := cur.Decode(currentMatch)
					if err != nil {
						sr.log.Errorf("Error decoding match: %s", err)
						continue
					}

					if currentMatch.MapID != 11 || currentMatch.QueueID < int(lowQueueID) || currentMatch.QueueID > int(highQueueID) {
						sr.log.Warnf("Found match which should not have been returned from storage, skipping...")
						continue
					}

					totalGamesForGameVersion++

					// Champion Picks
				ParticipantsLoop:
					for _, participant := range currentMatch.Participants {
						items := []int{
							participant.Stats.Item0, participant.Stats.Item1, participant.Stats.Item2,
							participant.Stats.Item3, participant.Stats.Item4, participant.Stats.Item5,
						}
						for _, item := range items {
							if item == 0 {
								// We only want complete builts
								continue ParticipantsLoop
							}
						}
						itemsHash := hashItems(items)

						role := participant.Timeline.Role
						lane := participant.Timeline.Lane
						cid := participant.ChampionID

						ccall := itemCountersAllTiers[cid]
						if _, ok := ccall.PerRole[lane]; !ok {
							ccall.PerRole[lane] = make(map[string]singleItemCounters)
						}
						perRoleAll := ccall.PerRole[lane][role]

						ccitem := ccall.SingleItemCounters[itemsHash]
						ccitemPerRole := perRoleAll[itemsHash]

						ccall.TotalPicks++
						ccitem.Picks++
						ccitemPerRole.Picks++

						if participant.Stats.Win {
							ccitem.Wins++
							ccitemPerRole.Wins++
						}

						ccall.PerRole[lane][role] = perRoleAll
						ccall.SingleItemCounters[itemsHash] = ccitem
						itemCountersAllTiers[cid] = ccall
					}

					cnt++
				}

				// Prepare results for ItemStats (ALL tiers)
				for cid, itemCounter := range itemCountersAllTiers {
					stats, err := sr.prepareItemStats(uint64(cid), version[0], version[1], itemCounter.TotalPicks, &itemCounter)
					if err == nil {
						err = sr.storage.StoreItemStats(stats)
						if err != nil {
							sr.log.Warnf("Something went wrong storing the Champion Item Stats: %s", err)
						}
					}
				}

				cur.Close()
				sr.log.Debugf("itemWinRateWorker calculation for Game Version %s done. Analyzed %d matches", gameVersion, cnt)
			}

			nextUpdate = time.Minute * time.Duration(sr.config.RScriptsUpdateInterval)
			elapsed := time.Since(start)
			sr.log.Infof("Finished matchAnalysisWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
			sr.calculationMutex.Unlock()
		}
	}
}

func (sr *StatsRunner) prepareItemStats(champID uint64, majorVersion uint32, minorVersion uint32, totalPicks uint64, itemCounter *itemCounter) (*storage.ItemStats, error) {

	gameVersion := fmt.Sprintf("%d.%d", majorVersion, minorVersion)

	itemStats := storage.ItemStats{}
	itemStats.ChampionID = champID
	itemStats.GameVersion = gameVersion

	itemStats.ItemStatsValues = make(storage.ItemStatsValues)

	itemStatsBySampleSize := make(map[float64]storage.SingleItemStatsValues)

	for itemCombination, itemCounts := range itemCounter.SingleItemCounters {
		if itemCounts.Picks > 1 {
			is := itemStats.ItemStatsValues[itemCombination]
			is.WinRate = float64(itemCounts.Wins) / float64(itemCounts.Picks)
			sr.log.Info(itemCounts.Wins, itemCounts.Picks, is.WinRate)
			is.PickRate = float64(itemCounts.Picks) / float64(totalPicks)
			is.SampleSize = itemCounts.Picks
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

	sr.log.Infof("Highest pick rate (%f) for Champ %d was item combination %s with a Win Rate of %f percent", itemStatsBySampleSize[keys[len(keys)-1]].PickRate, champID, itemStatsBySampleSize[keys[len(keys)-1]].ItemHash, 100*itemStatsBySampleSize[keys[len(keys)-1]].WinRate)

	champions := sr.storage.GetChampions(false)
	for _, val := range champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			itemStats.ChampionName = val.Name
			itemStats.ChampionRealID = val.ID
			break
		}
	}

	itemStats.Timestamp = time.Now()

	return &itemStats, nil
}
