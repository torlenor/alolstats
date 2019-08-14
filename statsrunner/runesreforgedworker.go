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

type runesReforgedSinglePickWinCounter struct {
	Picks         uint64
	Wins          uint64
	RunesReforged storage.RunesReforgedPicks
}
type runesReforgedSinglePickWinCounters map[string]runesReforgedSinglePickWinCounter // [runesReforgedHash]

type runesReforgedPickWinCounter struct {
	ChampionID int

	GameVersion string

	TotalPicks uint64

	TotalCounters runesReforgedSinglePickWinCounters
	PerRole       map[string]map[string]runesReforgedSinglePickWinCounters // [lane][role]
}
type runesReforgedPickWinCounters map[int]runesReforgedPickWinCounter            // [champID], e.g., 1, 10, 43, ...
type runesReforgedPickWinCountersPerTier map[string]runesReforgedPickWinCounters // [tier], e.g., "GOLD", "SILVER", "UNRANKED"

func (sr *StatsRunner) newPickWinCounters(champions riotclient.ChampionsList, gameVersion string) runesReforgedPickWinCounters {
	runesReforgedPickWinCounters := make(map[int]runesReforgedPickWinCounter)
	for _, champ := range champions {
		champID, err := strconv.Atoi(champ.Key)
		if err != nil {
			sr.log.Warnf("Something bad happened: %s", err)
			continue
		}
		runesReforgedPickWinCounters[champID] = runesReforgedPickWinCounter{
			ChampionID:  champID,
			GameVersion: gameVersion,

			TotalCounters: make(runesReforgedSinglePickWinCounters),
			PerRole:       make(map[string]map[string]runesReforgedSinglePickWinCounters),
		}
	}

	return runesReforgedPickWinCounters
}

func hashRunesReforged(summonerSpells []int) (string, []int) {
	sort.Ints(summonerSpells)
	var s string
	for _, item := range summonerSpells {
		s = s + strconv.Itoa(item) + "_"
	}

	return strings.TrimSuffix(s, "_"), summonerSpells
}

func (sr *StatsRunner) fillRunesReforgedPicks(runesReforgedDesc *riotclient.RunesReforgedList, stats *riotclient.ParticipantStatsDTO) storage.RunesReforgedPicks {
	runesReforgedPicks := storage.RunesReforgedPicks{}
	runesReforgedPicks.SlotPrimary.ID = stats.PerkPrimaryStyle
	runesReforgedPicks.SlotPrimary.Rune0.ID = stats.Perk0
	runesReforgedPicks.SlotPrimary.Rune1.ID = stats.Perk1
	runesReforgedPicks.SlotPrimary.Rune2.ID = stats.Perk2
	runesReforgedPicks.SlotPrimary.Rune3.ID = stats.Perk3

	for _, primary := range *runesReforgedDesc {
		if primary.ID == runesReforgedPicks.SlotPrimary.ID {
			runesReforgedPicks.SlotPrimary.Key = primary.Key
			runesReforgedPicks.SlotPrimary.Icon = primary.Icon
			runesReforgedPicks.SlotPrimary.Name = primary.Name
			for _, slot := range primary.Slots {
				for _, rune := range slot.Runes {
					if rune.ID == runesReforgedPicks.SlotPrimary.Rune0.ID {
						runesReforgedPicks.SlotPrimary.Rune0.Icon = rune.Icon
						runesReforgedPicks.SlotPrimary.Rune0.Key = rune.Key
						runesReforgedPicks.SlotPrimary.Rune0.Name = rune.Name
						continue
					}
					if rune.ID == runesReforgedPicks.SlotPrimary.Rune1.ID {
						runesReforgedPicks.SlotPrimary.Rune1.Icon = rune.Icon
						runesReforgedPicks.SlotPrimary.Rune1.Key = rune.Key
						runesReforgedPicks.SlotPrimary.Rune1.Name = rune.Name
						continue
					}
					if rune.ID == runesReforgedPicks.SlotPrimary.Rune2.ID {
						runesReforgedPicks.SlotPrimary.Rune2.Icon = rune.Icon
						runesReforgedPicks.SlotPrimary.Rune2.Key = rune.Key
						runesReforgedPicks.SlotPrimary.Rune2.Name = rune.Name
						continue
					}
					if rune.ID == runesReforgedPicks.SlotPrimary.Rune3.ID {
						runesReforgedPicks.SlotPrimary.Rune3.Icon = rune.Icon
						runesReforgedPicks.SlotPrimary.Rune3.Key = rune.Key
						runesReforgedPicks.SlotPrimary.Rune3.Name = rune.Name
						continue
					}
				}
			}
		}
	}

	runesReforgedPicks.SlotSecondary.ID = stats.PerkSubStyle
	runesReforgedPicks.SlotSecondary.Rune0.ID = stats.Perk4
	runesReforgedPicks.SlotSecondary.Rune1.ID = stats.Perk5

	for _, secondary := range *runesReforgedDesc {
		if secondary.ID == runesReforgedPicks.SlotSecondary.ID {
			runesReforgedPicks.SlotSecondary.Key = secondary.Key
			runesReforgedPicks.SlotSecondary.Icon = secondary.Icon
			runesReforgedPicks.SlotSecondary.Name = secondary.Name
			for _, slot := range secondary.Slots {
				for _, rune := range slot.Runes {
					if rune.ID == runesReforgedPicks.SlotSecondary.Rune0.ID {
						runesReforgedPicks.SlotSecondary.Rune0.Icon = rune.Icon
						runesReforgedPicks.SlotSecondary.Rune0.Key = rune.Key
						runesReforgedPicks.SlotSecondary.Rune0.Name = rune.Name
						continue
					}
					if rune.ID == runesReforgedPicks.SlotSecondary.Rune1.ID {
						runesReforgedPicks.SlotSecondary.Rune1.Icon = rune.Icon
						runesReforgedPicks.SlotSecondary.Rune1.Key = rune.Key
						runesReforgedPicks.SlotSecondary.Rune1.Name = rune.Name
						continue
					}
				}
			}
		}
	}

	runesReforgedPicks.StatPerks.Perk0.ID = stats.StatPerk0
	runesReforgedPicks.StatPerks.Perk1.ID = stats.StatPerk1
	runesReforgedPicks.StatPerks.Perk2.ID = stats.StatPerk2

	return runesReforgedPicks
}

func (sr *StatsRunner) runesReforgedWorker() {
	sr.workersWG.Add(1)
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

			champions := sr.storage.GetChampions(false)
			runesReforgedDesc := sr.storage.GetRunesReforged(false)
			if runesReforgedDesc == nil {
				sr.log.Errorf("Could not get Runes Reforged from Storage, canceling runesReforgedWorker run")
				nextUpdate = time.Minute * time.Duration(sr.config.RunesReforgedStats.UpdateInverval)
				continue
			}

			mapID := uint64(11)

			for queueID, queue := range queueIDtoQueue {
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
					sr.log.Infof("RunesReforgedWorker calculation for Game Version %s and Queue %s started", gameVersion, queue)

					runesReforgedPickWinCountersPerTier := make(runesReforgedPickWinCountersPerTier)
					runesReforgedPickWinCountersAllTiers := sr.newPickWinCounters(champions, gameVersion)

					totalGamesForGameVersion := uint64(0)
					totalGamesForGameVersionTier := make(map[string]uint64)

					cur, err := sr.storage.GetMatchesCursorByGameVersionMapQueueID(majorMinor, mapID, queueID)
					if err != nil {
						sr.log.Errorf("Error performing runesReforgedWorker calculation for Game Version %s: %s", gameVersion, err)
						continue
					}
					currentMatch := &riotclient.MatchDTO{}
					cnt := 0

					for cur.Next() {
						if sr.shouldWorkersStop {
							return
						}
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
							runesReforged := []int{
								participant.Stats.PerkPrimaryStyle,
								participant.Stats.Perk0,
								participant.Stats.Perk1,
								participant.Stats.Perk2,
								participant.Stats.Perk3,
								participant.Stats.PerkSubStyle,
								participant.Stats.Perk4,
								participant.Stats.Perk5,
								participant.Stats.StatPerk0,
								participant.Stats.StatPerk1,
								participant.Stats.StatPerk2,
							}
							runesReforgedHash, _ := hashRunesReforged(runesReforged)

							runesReforgedPicks := sr.fillRunesReforgedPicks(runesReforgedDesc, &participant.Stats)

							role := participant.Timeline.Role
							lane := participant.Timeline.Lane
							cid := participant.ChampionID

							// Get structs for counting
							if _, ok := runesReforgedPickWinCountersPerTier[matchTier]; !ok {
								runesReforgedPickWinCountersPerTier[matchTier] = sr.newPickWinCounters(champions, gameVersion)
							}
							cct := runesReforgedPickWinCountersPerTier[matchTier]
							cc := cct[cid]
							if _, ok := cc.PerRole[lane]; !ok {
								cc.PerRole[lane] = make(map[string]runesReforgedSinglePickWinCounters)
							}
							perRole := cc.PerRole[lane][role]

							ccSpellsCTier := cc.TotalCounters[runesReforgedHash]
							ccSpellsTierPerRole := perRole[runesReforgedHash]

							ccall := runesReforgedPickWinCountersAllTiers[cid]
							if _, ok := ccall.PerRole[lane]; !ok {
								ccall.PerRole[lane] = make(map[string]runesReforgedSinglePickWinCounters)
							}
							if _, ok := ccall.PerRole[lane][role]; !ok {
								ccall.PerRole[lane][role] = make(runesReforgedSinglePickWinCounters)
							}
							perRoleAll := ccall.PerRole[lane][role]

							ccSpellsCAll := ccall.TotalCounters[runesReforgedHash]
							ccSpellsAllPerRole := perRoleAll[runesReforgedHash]

							ccall.TotalPicks++
							cc.TotalPicks++

							ccSpellsCTier.Picks++
							ccSpellsCTier.RunesReforged = runesReforgedPicks
							ccSpellsTierPerRole.Picks++
							ccSpellsTierPerRole.RunesReforged = runesReforgedPicks

							ccSpellsCAll.Picks++
							ccSpellsCAll.RunesReforged = runesReforgedPicks
							ccSpellsAllPerRole.Picks++
							ccSpellsAllPerRole.RunesReforged = runesReforgedPicks

							if participant.Stats.Win {
								ccSpellsCTier.Wins++
								ccSpellsTierPerRole.Wins++

								ccSpellsCAll.Wins++
								ccSpellsAllPerRole.Wins++
							}

							// Backassign structs
							cc.PerRole[lane][role] = perRole
							cct[cid] = cc
							runesReforgedPickWinCountersPerTier[matchTier] = cct

							perRoleAll[runesReforgedHash] = ccSpellsAllPerRole
							ccall.PerRole[lane][role] = perRoleAll
							ccall.TotalCounters[runesReforgedHash] = ccSpellsCAll
							runesReforgedPickWinCountersAllTiers[cid] = ccall
						}

						cnt++
					}

					// Prepare results for Summoner Spells Stats (ALL tiers)
					for cid, runesReforgedPickWinCounter := range runesReforgedPickWinCountersAllTiers {
						stats, err := sr.prepareRunesReforgedStats(uint64(cid), version[0], version[1], runesReforgedPickWinCounter.TotalPicks, &runesReforgedPickWinCounter)
						stats.Tier = "ALL"
						stats.Queue = queue
						if err == nil {
							err = sr.storage.StoreRunesReforgedStats(stats)
							if err != nil {
								sr.log.Errorf("Something went wrong storing the Summoner Spells Stats: %s", err)
							}
						} else {
							sr.log.Errorf("Something went wrong calculating the Summoner Spells Stats: %s", err)
						}
					}

					// Prepare results for Summoner Spells Stats (per tier)
					for tier, runesReforgedPickWinCounters := range runesReforgedPickWinCountersPerTier {
						for cid, runesReforgedPickWinCounter := range runesReforgedPickWinCounters {
							stats, err := sr.prepareRunesReforgedStats(uint64(cid), version[0], version[1], runesReforgedPickWinCounter.TotalPicks, &runesReforgedPickWinCounter)
							stats.Tier = tier
							stats.Queue = queue
							if err == nil {
								err = sr.storage.StoreRunesReforgedStats(stats)
								if err != nil {
									sr.log.Warnf("Something went wrong storing the Champion Item Stats: %s", err)
								}
							}
						}
					}

					cur.Close()
					sr.log.Infof("RunesReforgedWorker calculation for Game Version %s and Queue %s done. Analyzed %d matches", gameVersion, queue, cnt)
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

			nextUpdate = time.Minute * time.Duration(sr.config.RunesReforgedStats.UpdateInverval)
			elapsed := time.Since(start)
			sr.log.Infof("Finished runesReforgedWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
			sr.calculationMutex.Unlock()
		}
	}
}

func (sr *StatsRunner) prepareRunesReforgedStats(champID uint64, majorVersion uint32, minorVersion uint32, totalPicks uint64, cc *runesReforgedPickWinCounter) (*storage.RunesReforgedStats, error) {

	gameVersion := fmt.Sprintf("%d.%d", majorVersion, minorVersion)

	runesReforgedStats := storage.RunesReforgedStats{}
	runesReforgedStats.ChampionID = champID
	runesReforgedStats.GameVersion = gameVersion
	runesReforgedStats.SampleSize = cc.TotalPicks

	runesReforgedStats.Stats = make(storage.RunesReforgedStatsValues)

	for runesReforgedHash, counters := range cc.TotalCounters {
		is := runesReforgedStats.Stats[runesReforgedHash]
		is.WinRate = float64(counters.Wins) / float64(counters.Picks)
		is.PickRate = float64(counters.Picks) / float64(totalPicks)
		is.SampleSize = counters.Picks
		is.RunesReforged = counters.RunesReforged
		runesReforgedStats.Stats[runesReforgedHash] = is
	}

	if sr.config.RunesReforgedStats.KeepOnlyHighestPickRate {
		var highestPickRateHash string
		highestSampleSize := uint64(0)
		for runesReforgedHash, stats := range runesReforgedStats.Stats {
			if stats.SampleSize > highestSampleSize {
				highestPickRateHash = runesReforgedHash
				highestSampleSize = stats.SampleSize
			}
		}
		runesReforgedStats.Stats = storage.RunesReforgedStatsValues{
			highestPickRateHash: runesReforgedStats.Stats[highestPickRateHash],
		}
	}

	// Calculation of stats per Role
	runesReforgedStats.StatsPerRole = make(map[string]storage.RunesReforgedStatsValues)
	for _, role := range []string{"Top", "Mid", "Jungle", "Carry", "Support"} {
		statsValues := sr.calculateRunesReforgedRoleStats(cc, role)

		if sr.config.RunesReforgedStats.KeepOnlyHighestPickRate {
			var highestPickRateHash string
			highestSampleSize := uint64(0)
			for runesReforgedHash, stats := range statsValues {
				if stats.SampleSize > highestSampleSize {
					highestPickRateHash = runesReforgedHash
					highestSampleSize = stats.SampleSize
				}
			}
			if highestPickRateHash != "" {
				runesReforgedStats.StatsPerRole[role] = storage.RunesReforgedStatsValues{
					highestPickRateHash: statsValues[highestPickRateHash],
				}
			} else {
				runesReforgedStats.StatsPerRole[role] = storage.RunesReforgedStatsValues{}
			}
		} else {
			runesReforgedStats.StatsPerRole[role] = statsValues
		}

	}

	champions := sr.storage.GetChampions(false)
	for _, val := range champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			runesReforgedStats.ChampionName = val.Name
			runesReforgedStats.ChampionRealID = val.ID
			break
		}
	}

	runesReforgedStats.Timestamp = time.Now()

	return &runesReforgedStats, nil
}

func (sr *StatsRunner) calculateRunesReforgedRoleStats(cc *runesReforgedPickWinCounter, role string) storage.RunesReforgedStatsValues {
	summedCounters := runesReforgedSinglePickWinCounters{}

	switch role {
	case "Top":
		for _, cnters := range cc.PerRole["TOP"] {
			sumRunesReforgedSinglePickWinCounters(&summedCounters, cnters)
		}
	case "Mid":
		for _, cnters := range cc.PerRole["MIDDLE"] {
			sumRunesReforgedSinglePickWinCounters(&summedCounters, cnters)
		}
	case "Jungle":
		for _, cnters := range cc.PerRole["JUNGLE"] {
			sumRunesReforgedSinglePickWinCounters(&summedCounters, cnters)
		}
	case "Carry":
		if lane, ok := cc.PerRole["BOTTOM"]; ok {
			if cnters, ok := lane["DUO_CARRY"]; ok {
				sumRunesReforgedSinglePickWinCounters(&summedCounters, cnters)
			}
		}
	case "Support":
		if lane, ok := cc.PerRole["BOTTOM"]; ok {
			if cnters, ok := lane["DUO_SUPPORT"]; ok {
				sumRunesReforgedSinglePickWinCounters(&summedCounters, cnters)
			}
		}
	}

	return sr.calcRunesReforgedStatsFromCounters(&summedCounters)
}

func sumRunesReforgedSinglePickWinCounters(summedCounters *runesReforgedSinglePickWinCounters, countersToAdd runesReforgedSinglePickWinCounters) {
	for key, counter := range countersToAdd {
		count := (*summedCounters)[key]
		count.Picks += counter.Picks
		count.Wins += counter.Wins
		count.RunesReforged = counter.RunesReforged
		(*summedCounters)[key] = count
	}
}

func (sr *StatsRunner) calcRunesReforgedStatsFromCounters(counters *runesReforgedSinglePickWinCounters) storage.RunesReforgedStatsValues {
	statsValues := storage.RunesReforgedStatsValues{}

	var totalCount uint64
	for _, count := range *counters {
		totalCount += count.Picks
	}

	for runesReforgedHash, count := range *counters {
		stats := statsValues[runesReforgedHash]
		if count.Picks > 0 {
			stats.SampleSize = count.Picks
			stats.RunesReforged = count.RunesReforged
			stats.WinRate = float64(count.Wins) / float64(count.Picks)
			stats.PickRate = float64(count.Picks) / float64(totalCount)
		}
		statsValues[runesReforgedHash] = stats
	}

	return statsValues
}
