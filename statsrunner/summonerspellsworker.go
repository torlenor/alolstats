package statsrunner

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"git.abyle.org/hps/alolstats/riotclient"
	"git.abyle.org/hps/alolstats/storage"
	"git.abyle.org/hps/alolstats/utils"
)

type summonerSpellsSinglePickWinCounter struct {
	Picks          uint64
	Wins           uint64
	SummonerSpells []riotclient.SummonerSpell
}
type summonerSpellsSinglePickWinCounters map[string]summonerSpellsSinglePickWinCounter // [summonerSpellsHash]

type summonerSpellsCounter struct {
	ChampionID int

	GameVersion string

	TotalPicks uint64

	TotalCounters summonerSpellsSinglePickWinCounters
	PerRole       map[string]map[string]summonerSpellsSinglePickWinCounters // [lane][role]
}
type summonerSpellsCounters map[int]summonerSpellsCounter            // [champID], e.g., 1, 10, 43, ...
type summonerSpellsCountersPerTier map[string]summonerSpellsCounters // [tier], e.g., "GOLD", "SILVER", "UNRANKED"

func (sr *StatsRunner) newSummonerSpellsCounters(champions riotclient.ChampionsList, gameVersion string) summonerSpellsCounters {
	summonerSpellsCounters := make(map[int]summonerSpellsCounter)
	for _, champ := range champions {
		champID, err := strconv.Atoi(champ.Key)
		if err != nil {
			sr.log.Warnf("Something bad happened: %s", err)
			continue
		}
		summonerSpellsCounters[champID] = summonerSpellsCounter{
			ChampionID:  champID,
			GameVersion: gameVersion,

			TotalCounters: make(summonerSpellsSinglePickWinCounters),
			PerRole:       make(map[string]map[string]summonerSpellsSinglePickWinCounters),
		}
	}

	return summonerSpellsCounters
}

func hashSummonerSpells(summonerSpells []int) (string, []int) {
	sort.Ints(summonerSpells)
	var s string
	for _, item := range summonerSpells {
		s = s + strconv.Itoa(item) + "_"
	}

	return strings.TrimSuffix(s, "_"), summonerSpells
}

func (sr *StatsRunner) summonerSpellsWorker() {
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
			sr.log.Printf("Stopping summonerSpellsWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			sr.calculationMutex.Lock()
			sr.log.Infof("Performing summonerSpellsWorker run")
			start := time.Now()

			champions := sr.storage.GetChampions(false)
			summonerSpellsDesc := sr.storage.GetSummonerSpells(false)
			if summonerSpellsDesc == nil {
				sr.log.Errorf("Could not get Summoner Spells from Storage, canceling summonerSpellsWorker run")
				nextUpdate = time.Minute * time.Duration(sr.config.SummonerSpellsStats.UpdateInverval)
				continue
			}

			mapID := uint64(11)

			for queueID, queue := range queueIDtoQueue {
				for _, versionStr := range sr.config.GameVersion {
					sr.shouldWorkersSopMutex.RLock()
					if sr.shouldWorkersStop {
						return
					}
					sr.shouldWorkersSopMutex.RUnlock()
					version, err := utils.SplitNumericVersion(versionStr)
					if err != nil {
						sr.log.Warnf("Something bad happened: %s", err)
						continue
					}
					gameVersion := fmt.Sprintf("%d.%d", version[0], version[1])
					majorMinor := fmt.Sprintf("%d\\.%d\\.", version[0], version[1])
					sr.log.Infof("SummonerSpellsWorker calculation for Game Version %s and Queue %s started", gameVersion, queue)

					summonerSpellsCountersPerTier := make(summonerSpellsCountersPerTier)
					summonerSpellsCountersAllTiers := sr.newSummonerSpellsCounters(champions, gameVersion)

					totalGamesForGameVersion := uint64(0)
					totalGamesForGameVersionTier := make(map[string]uint64)

					cur, err := sr.storage.GetMatchesCursorByGameVersionMapQueueID(majorMinor, mapID, queueID)
					if err != nil {
						sr.log.Errorf("Error performing summonerSpellsWorker calculation for Game Version %s: %s", gameVersion, err)
						continue
					}
					currentMatch := &riotclient.MatchDTO{}
					cnt := 0

					for cur.Next() {
						sr.shouldWorkersSopMutex.RLock()
						if sr.shouldWorkersStop {
							return
						}
						sr.shouldWorkersSopMutex.RUnlock()
						err := cur.Decode(currentMatch)
						if err != nil {
							sr.log.Errorf("Error decoding match: %s", err)
							continue
						}

						if currentMatch.MapID != int(mapID) || currentMatch.QueueID != int(queueID) {
							sr.log.Warnf("Found match which should not have been returned from storage, skipping...")
							continue
						}

						totalGamesForGameVersion++

						matchTier := determineMatchTier(currentMatch.Participants)
						totalGamesForGameVersionTier[matchTier]++

						for _, participant := range currentMatch.Participants {
							summonerSpells := []int{
								participant.Spell1ID, participant.Spell2ID,
							}
							summonerSpellsHash, sortedSummonerSpells := hashSummonerSpells(summonerSpells)

							role := participant.Timeline.Role
							lane := participant.Timeline.Lane
							cid := participant.ChampionID

							summonerSpellsList := []riotclient.SummonerSpell{}
							for _, spell := range sortedSummonerSpells {
								for _, val := range *summonerSpellsDesc {
									if strconv.Itoa(spell) == val.Key {
										summonerSpellsList = append(summonerSpellsList, val)
									}
								}
							}

							// Get structs for counting
							if _, ok := summonerSpellsCountersPerTier[matchTier]; !ok {
								summonerSpellsCountersPerTier[matchTier] = sr.newSummonerSpellsCounters(champions, gameVersion)
							}
							cct := summonerSpellsCountersPerTier[matchTier]
							cc := cct[cid]
							if _, ok := cc.PerRole[lane]; !ok {
								cc.PerRole[lane] = make(map[string]summonerSpellsSinglePickWinCounters)
							}
							perRole := cc.PerRole[lane][role]

							ccSpellsCTier := cc.TotalCounters[summonerSpellsHash]
							ccSpellsTierPerRole := perRole[summonerSpellsHash]

							ccall := summonerSpellsCountersAllTiers[cid]
							if _, ok := ccall.PerRole[lane]; !ok {
								ccall.PerRole[lane] = make(map[string]summonerSpellsSinglePickWinCounters)
							}
							if _, ok := ccall.PerRole[lane][role]; !ok {
								ccall.PerRole[lane][role] = make(summonerSpellsSinglePickWinCounters)
							}
							perRoleAll := ccall.PerRole[lane][role]

							ccSpellsCAll := ccall.TotalCounters[summonerSpellsHash]
							ccSpellsAllPerRole := perRoleAll[summonerSpellsHash]

							ccall.TotalPicks++
							cc.TotalPicks++

							ccSpellsCTier.Picks++
							ccSpellsCTier.SummonerSpells = summonerSpellsList
							ccSpellsTierPerRole.Picks++
							ccSpellsTierPerRole.SummonerSpells = summonerSpellsList

							ccSpellsCAll.Picks++
							ccSpellsCAll.SummonerSpells = summonerSpellsList
							ccSpellsAllPerRole.Picks++
							ccSpellsAllPerRole.SummonerSpells = summonerSpellsList

							if participant.Stats.Win {
								ccSpellsCTier.Wins++
								ccSpellsTierPerRole.Wins++

								ccSpellsCAll.Wins++
								ccSpellsAllPerRole.Wins++
							}

							// Backassign structs
							cc.PerRole[lane][role] = perRole
							cct[cid] = cc
							summonerSpellsCountersPerTier[matchTier] = cct

							perRoleAll[summonerSpellsHash] = ccSpellsAllPerRole
							ccall.PerRole[lane][role] = perRoleAll
							ccall.TotalCounters[summonerSpellsHash] = ccSpellsCAll
							summonerSpellsCountersAllTiers[cid] = ccall
						}

						cnt++
					}

					// Prepare results for Summoner Spells Stats (ALL tiers)
					for cid, summonerSpellsCounter := range summonerSpellsCountersAllTiers {
						stats, err := sr.prepareSummonerSpellsStats(uint64(cid), version[0], version[1], summonerSpellsCounter.TotalPicks, &summonerSpellsCounter)
						stats.Tier = "ALL"
						stats.Queue = queue
						if err == nil {
							err = sr.storage.StoreSummonerSpellsStats(stats)
							if err != nil {
								sr.log.Errorf("Something went wrong storing the Summoner Spells Stats: %s", err)
							}
						} else {
							sr.log.Errorf("Something went wrong calculating the Summoner Spells Stats: %s", err)
						}
					}

					// Prepare results for Summoner Spells Stats (per tier)
					for tier, summonerSpellsCounters := range summonerSpellsCountersPerTier {
						for cid, summonerSpellsCounter := range summonerSpellsCounters {
							stats, err := sr.prepareSummonerSpellsStats(uint64(cid), version[0], version[1], summonerSpellsCounter.TotalPicks, &summonerSpellsCounter)
							stats.Tier = tier
							stats.Queue = queue
							if err == nil {
								err = sr.storage.StoreSummonerSpellsStats(stats)
								if err != nil {
									sr.log.Warnf("Something went wrong storing the Champion Item Stats: %s", err)
								}
							}
						}
					}

					cur.Close()
					sr.log.Infof("SummonerSpellsWorker calculation for Game Version %s and Queue %s done. Analyzed %d matches", gameVersion, queue, cnt)
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

			nextUpdate = time.Minute * time.Duration(sr.config.SummonerSpellsStats.UpdateInverval)
			elapsed := time.Since(start)
			sr.log.Infof("Finished summonerSpellsWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
			sr.calculationMutex.Unlock()
		}
	}
}

func (sr *StatsRunner) prepareSummonerSpellsStats(champID uint64, majorVersion uint32, minorVersion uint32, totalPicks uint64, cc *summonerSpellsCounter) (*storage.SummonerSpellsStats, error) {

	gameVersion := fmt.Sprintf("%d.%d", majorVersion, minorVersion)

	summonerSpellsStats := storage.SummonerSpellsStats{}
	summonerSpellsStats.ChampionID = champID
	summonerSpellsStats.GameVersion = gameVersion
	summonerSpellsStats.SampleSize = cc.TotalPicks

	summonerSpellsStats.Stats = make(storage.SummonerSpellsStatsValues)

	for summonerSpellsHash, counters := range cc.TotalCounters {
		is := summonerSpellsStats.Stats[summonerSpellsHash]
		is.WinRate = float64(counters.Wins) / float64(counters.Picks)
		is.PickRate = float64(counters.Picks) / float64(totalPicks)
		is.SampleSize = counters.Picks
		is.SummonerSpells = counters.SummonerSpells
		summonerSpellsStats.Stats[summonerSpellsHash] = is
	}

	if sr.config.SummonerSpellsStats.KeepOnlyHighestPickRate {
		var highestPickRateHash string
		highestSampleSize := uint64(0)
		for summonerSpellsHash, stats := range summonerSpellsStats.Stats {
			if stats.SampleSize > highestSampleSize {
				highestPickRateHash = summonerSpellsHash
				highestSampleSize = stats.SampleSize
			}
		}
		summonerSpellsStats.Stats = storage.SummonerSpellsStatsValues{
			highestPickRateHash: summonerSpellsStats.Stats[highestPickRateHash],
		}
	}

	// Calculation of stats per Role
	summonerSpellsStats.StatsPerRole = make(map[string]storage.SummonerSpellsStatsValues)
	for _, role := range []string{"Top", "Mid", "Jungle", "Carry", "Support"} {
		statsValues := sr.calculateSummonerSpellsRoleStats(cc, role)

		if sr.config.SummonerSpellsStats.KeepOnlyHighestPickRate {
			var highestPickRateHash string
			highestSampleSize := uint64(0)
			for summonerSpellsHash, stats := range statsValues {
				if stats.SampleSize > highestSampleSize {
					highestPickRateHash = summonerSpellsHash
					highestSampleSize = stats.SampleSize
				}
			}
			if highestPickRateHash != "" {
				summonerSpellsStats.StatsPerRole[role] = storage.SummonerSpellsStatsValues{
					highestPickRateHash: statsValues[highestPickRateHash],
				}
			} else {
				summonerSpellsStats.StatsPerRole[role] = storage.SummonerSpellsStatsValues{}
			}
		} else {
			summonerSpellsStats.StatsPerRole[role] = statsValues
		}

	}

	champions := sr.storage.GetChampions(false)
	for _, val := range champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			summonerSpellsStats.ChampionName = val.Name
			summonerSpellsStats.ChampionRealID = val.ID
			break
		}
	}

	summonerSpellsStats.Timestamp = time.Now()

	return &summonerSpellsStats, nil
}

func (sr *StatsRunner) calculateSummonerSpellsRoleStats(cc *summonerSpellsCounter, role string) storage.SummonerSpellsStatsValues {
	summedCounters := summonerSpellsSinglePickWinCounters{}

	switch role {
	case "Top":
		for _, cnters := range cc.PerRole["TOP"] {
			sumSinglePickWinCounters(&summedCounters, cnters)
		}
	case "Mid":
		for _, cnters := range cc.PerRole["MIDDLE"] {
			sumSinglePickWinCounters(&summedCounters, cnters)
		}
	case "Jungle":
		for _, cnters := range cc.PerRole["JUNGLE"] {
			sumSinglePickWinCounters(&summedCounters, cnters)
		}
	case "Carry":
		if lane, ok := cc.PerRole["BOTTOM"]; ok {
			if cnters, ok := lane["DUO_CARRY"]; ok {
				sumSinglePickWinCounters(&summedCounters, cnters)
			}
		}
	case "Support":
		if lane, ok := cc.PerRole["BOTTOM"]; ok {
			if cnters, ok := lane["DUO_SUPPORT"]; ok {
				sumSinglePickWinCounters(&summedCounters, cnters)
			}
		}
	}

	return sr.calcSummonerSpellsStatsFromCounters(&summedCounters)
}

func sumSinglePickWinCounters(summedCounters *summonerSpellsSinglePickWinCounters, countersToAdd summonerSpellsSinglePickWinCounters) {
	for key, counter := range countersToAdd {
		count := (*summedCounters)[key]
		count.Picks += counter.Picks
		count.Wins += counter.Wins
		count.SummonerSpells = counter.SummonerSpells
		(*summedCounters)[key] = count
	}
}

func (sr *StatsRunner) calcSummonerSpellsStatsFromCounters(counters *summonerSpellsSinglePickWinCounters) storage.SummonerSpellsStatsValues {
	statsValues := storage.SummonerSpellsStatsValues{}

	var totalCount uint64
	for _, count := range *counters {
		totalCount += count.Picks
	}

	for summonerSpellsHash, count := range *counters {
		stats := statsValues[summonerSpellsHash]
		if count.Picks > 0 {
			stats.SampleSize = count.Picks
			stats.SummonerSpells = count.SummonerSpells
			stats.WinRate = float64(count.Wins) / float64(count.Picks)
			stats.PickRate = float64(count.Picks) / float64(totalCount)
		}
		statsValues[summonerSpellsHash] = stats
	}

	return statsValues
}
