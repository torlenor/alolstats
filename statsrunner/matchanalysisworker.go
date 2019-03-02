package statsrunner

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/torlenor/alolstats/riotclient"
	"github.com/torlenor/alolstats/storage"

	"github.com/torlenor/alolstats/utils"
)

type roleCounters struct {
	Picks uint64
	Wins  uint64

	Kills   uint64
	Deaths  uint64
	Assists uint64

	MatchKills   []float64
	MatchDeaths  []float64
	MatchAssists []float64

	MatchGoldEarned         []float64
	MatchTotalMinionsKilled []float64

	MatchTotalDamageDealt               []float64
	MatchTotalDamageDealtToChampions    []float64
	MatchTotalDamageTaken               []float64
	MatchMagicDamageDealt               []float64
	MatchMagicDamageDealtToChampions    []float64
	MatchPhysicalDamageDealt            []float64
	MatchPhysicalDamageDealtToChampions []float64
	MatchPhysicalDamageTaken            []float64
	MatchTrueDamageDealt                []float64
	MatchTrueDamageDealtToChampions     []float64
	MatchTrueDamageTaken                []float64

	MatchTotalHeal []float64

	MatchDamageDealtToObjectives []float64
	MatchDamageDealtToTurrets    []float64
	MatchTimeCCingOthers         []float64
}

type championCounters struct {
	ChampionID int

	GameVersion string

	TotalPicks uint64
	TotalWins  uint64

	TotalBans uint64

	TotalKills   uint64
	TotalDeaths  uint64
	TotalAssists uint64

	MatchKills   []float64
	MatchDeaths  []float64
	MatchAssists []float64

	MatchGoldEarned         []float64
	MatchTotalMinionsKilled []float64

	MatchTotalDamageDealt               []float64
	MatchTotalDamageDealtToChampions    []float64
	MatchTotalDamageTaken               []float64
	MatchMagicDamageDealt               []float64
	MatchMagicDamageDealtToChampions    []float64
	MatchPhysicalDamageDealt            []float64
	MatchPhysicalDamageDealtToChampions []float64
	MatchPhysicalDamageTaken            []float64
	MatchTrueDamageDealt                []float64
	MatchTrueDamageDealtToChampions     []float64
	MatchTrueDamageTaken                []float64

	MatchTotalHeal []float64

	MatchDamageDealtToObjectives []float64
	MatchDamageDealtToTurrets    []float64
	MatchTimeCCingOthers         []float64

	PerRole map[string]map[string]roleCounters // [lane][role]
}

type championsCounters map[int]championCounters            // [id], e.g., 1, 10, 43, ...
type championsCountersPerTier map[string]championsCounters // [tier], e.g., "GOLD", "SILVER", "UNRANKED"

func determineMatchTier(participants []riotclient.ParticipantDTO) string {
	tierCounts := make(map[string]uint64)
	for _, participant := range participants {
		currentTier := strings.ToUpper(strings.TrimSpace(participant.HighestAchievedSeasonTier))
		tierCounts[currentTier]++
	}

	var maxTier string
	maxCounts := uint64(0)
	for tier, cnt := range tierCounts {
		if cnt > maxCounts {
			maxTier = tier
			maxCounts = cnt
		}
	}

	if maxTier == "" {
		return "UNRANKED"
	}

	return maxTier
}

func (sr *StatsRunner) newChampionsCounters(champions riotclient.ChampionsList, gameVersion string) championsCounters {
	champsCounters := make(championsCounters)
	for _, champ := range champions {
		champID, err := strconv.Atoi(champ.Key)
		if err != nil {
			sr.log.Warnf("Something bad happened: %s", err)
			continue
		}
		champsCounters[champID] = championCounters{
			ChampionID:  champID,
			GameVersion: gameVersion,

			PerRole: make(map[string]map[string]roleCounters),
		}
	}

	return champsCounters
}

func (sr *StatsRunner) matchAnalysisWorker() {
	sr.workersWG.Add(1)
	defer sr.workersWG.Done()

	var nextUpdate time.Duration

	for {
		select {
		case <-sr.stopWorkers:
			sr.log.Printf("Stopping matchAnalysisWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			sr.log.Infof("Performing matchAnalysisWorker run")
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
				sr.log.Debugf("matchAnalysisWorker calculation for Game Version %s started", gameVersion)

				// Prepare championsCountersPerTier
				champsCountersPerTier := make(championsCountersPerTier)
				champsCountersAllTiers := sr.newChampionsCounters(champions, gameVersion)

				totalGamesForGameVersion := uint64(0)
				totalGamesForGameVersionTier := make(map[string]uint64)

				cur, err := sr.storage.GetMatchesCursorByGameVersionMapBetweenQueueIDs(majorMinor, mapID, highQueueID, lowQueueID)
				if err != nil {
					sr.log.Errorf("Error performing matchAnalysisWorker calculation for Game Version %s: %s", gameVersion, err)
					continue
				}
				currentMatch := &riotclient.MatchDTO{}
				cnt := 0
				for cur.Next() {
					err := cur.Decode(currentMatch)
					if err != nil {
						sr.log.Errorf("Error deconding match: %s", err)
						continue
					}

					if currentMatch.MapID != 11 || currentMatch.QueueID < int(lowQueueID) || currentMatch.QueueID > int(highQueueID) {
						sr.log.Warnf("Found match which should not have been returned from storage, skipping...")
						continue
					}

					totalGamesForGameVersion++

					matchTier := determineMatchTier(currentMatch.Participants)
					totalGamesForGameVersionTier[matchTier]++

					// Champion Picks
					for _, participant := range currentMatch.Participants {
						role := participant.Timeline.Role
						lane := participant.Timeline.Lane
						cid := participant.ChampionID

						// Get structs for counting
						if _, ok := champsCountersPerTier[matchTier]; !ok {
							champsCountersPerTier[matchTier] = sr.newChampionsCounters(champions, gameVersion)
						}
						cct := champsCountersPerTier[matchTier]
						cc := cct[cid]
						if _, ok := cc.PerRole[lane]; !ok {
							cc.PerRole[lane] = make(map[string]roleCounters)
						}
						perRole := cc.PerRole[lane][role]

						ccall := champsCountersAllTiers[cid]
						if _, ok := ccall.PerRole[lane]; !ok {
							ccall.PerRole[lane] = make(map[string]roleCounters)
						}
						perRoleAll := ccall.PerRole[lane][role]

						// Do counts
						cc.TotalPicks++
						ccall.TotalPicks++

						cc.TotalKills = cc.TotalKills + uint64(participant.Stats.Kills)
						cc.MatchKills = append(cc.MatchKills, float64(participant.Stats.Kills))
						cc.TotalDeaths = cc.TotalDeaths + uint64(participant.Stats.Deaths)
						cc.MatchDeaths = append(cc.MatchDeaths, float64(participant.Stats.Deaths))
						cc.TotalAssists = cc.TotalAssists + uint64(participant.Stats.Assists)
						cc.MatchAssists = append(cc.MatchAssists, float64(participant.Stats.Assists))

						ccall.TotalKills = ccall.TotalKills + uint64(participant.Stats.Kills)
						ccall.MatchKills = append(ccall.MatchKills, float64(participant.Stats.Kills))
						ccall.TotalDeaths = ccall.TotalDeaths + uint64(participant.Stats.Deaths)
						ccall.MatchDeaths = append(ccall.MatchDeaths, float64(participant.Stats.Deaths))
						ccall.TotalAssists = ccall.TotalAssists + uint64(participant.Stats.Assists)
						ccall.MatchAssists = append(ccall.MatchAssists, float64(participant.Stats.Assists))

						perRole.Picks++
						perRoleAll.Picks++

						perRole.Kills = perRole.Kills + uint64(participant.Stats.Kills)
						perRole.MatchKills = append(perRole.MatchKills, float64(participant.Stats.Kills))
						perRole.Deaths = perRole.Deaths + uint64(participant.Stats.Deaths)
						perRole.MatchDeaths = append(perRole.MatchDeaths, float64(participant.Stats.Deaths))
						perRole.Assists = perRole.Assists + uint64(participant.Stats.Assists)
						perRole.MatchAssists = append(perRole.MatchAssists, float64(participant.Stats.Assists))

						perRoleAll.Kills = perRoleAll.Kills + uint64(participant.Stats.Kills)
						perRoleAll.MatchKills = append(perRoleAll.MatchKills, float64(participant.Stats.Kills))
						perRoleAll.Deaths = perRoleAll.Deaths + uint64(participant.Stats.Deaths)
						perRoleAll.MatchDeaths = append(perRoleAll.MatchDeaths, float64(participant.Stats.Deaths))
						perRoleAll.Assists = perRoleAll.Assists + uint64(participant.Stats.Assists)
						perRoleAll.MatchAssists = append(perRoleAll.MatchAssists, float64(participant.Stats.Assists))

						ccall.MatchGoldEarned = append(ccall.MatchGoldEarned, float64(participant.Stats.GoldEarned))
						ccall.MatchTotalMinionsKilled = append(ccall.MatchTotalMinionsKilled, float64(participant.Stats.TotalMinionsKilled))
						ccall.MatchTotalDamageDealt = append(ccall.MatchTotalDamageDealt, float64(participant.Stats.TotalDamageDealt))
						ccall.MatchTotalDamageDealtToChampions = append(ccall.MatchTotalDamageDealtToChampions, float64(participant.Stats.TotalDamageDealtToChampions))
						ccall.MatchTotalDamageTaken = append(ccall.MatchTotalDamageTaken, float64(participant.Stats.TotalDamageTaken))
						ccall.MatchMagicDamageDealt = append(ccall.MatchMagicDamageDealt, float64(participant.Stats.MagicDamageDealt))
						ccall.MatchMagicDamageDealtToChampions = append(ccall.MatchMagicDamageDealtToChampions, float64(participant.Stats.MagicDamageDealtToChampions))
						ccall.MatchPhysicalDamageDealt = append(ccall.MatchPhysicalDamageDealt, float64(participant.Stats.PhysicalDamageDealt))
						ccall.MatchPhysicalDamageDealtToChampions = append(ccall.MatchPhysicalDamageDealtToChampions, float64(participant.Stats.PhysicalDamageDealtToChampions))
						ccall.MatchPhysicalDamageTaken = append(ccall.MatchPhysicalDamageTaken, float64(participant.Stats.PhysicalDamageTaken))
						ccall.MatchTrueDamageDealt = append(ccall.MatchTrueDamageDealt, float64(participant.Stats.TrueDamageDealt))
						ccall.MatchTrueDamageDealtToChampions = append(ccall.MatchTrueDamageDealtToChampions, float64(participant.Stats.TrueDamageDealtToChampions))
						ccall.MatchTrueDamageTaken = append(ccall.MatchTrueDamageTaken, float64(participant.Stats.TrueDamageTaken))
						ccall.MatchTotalHeal = append(ccall.MatchTotalHeal, float64(participant.Stats.TotalHeal))
						ccall.MatchDamageDealtToObjectives = append(perRole.MatchDamageDealtToObjectives, float64(participant.Stats.DamageDealtToObjectives))
						ccall.MatchDamageDealtToTurrets = append(perRole.MatchDamageDealtToTurrets, float64(participant.Stats.DamageDealtToTurrets))
						ccall.MatchTimeCCingOthers = append(perRole.MatchTimeCCingOthers, float64(participant.Stats.TimeCCingOthers))

						cc.MatchGoldEarned = append(cc.MatchGoldEarned, float64(participant.Stats.GoldEarned))
						cc.MatchTotalMinionsKilled = append(cc.MatchTotalMinionsKilled, float64(participant.Stats.TotalMinionsKilled))
						cc.MatchTotalDamageDealt = append(cc.MatchTotalDamageDealt, float64(participant.Stats.TotalDamageDealt))
						cc.MatchTotalDamageDealtToChampions = append(cc.MatchTotalDamageDealtToChampions, float64(participant.Stats.TotalDamageDealtToChampions))
						cc.MatchTotalDamageTaken = append(cc.MatchTotalDamageTaken, float64(participant.Stats.TotalDamageTaken))
						cc.MatchMagicDamageDealt = append(cc.MatchMagicDamageDealt, float64(participant.Stats.MagicDamageDealt))
						cc.MatchMagicDamageDealtToChampions = append(cc.MatchMagicDamageDealtToChampions, float64(participant.Stats.MagicDamageDealtToChampions))
						cc.MatchPhysicalDamageDealt = append(cc.MatchPhysicalDamageDealt, float64(participant.Stats.PhysicalDamageDealt))
						cc.MatchPhysicalDamageDealtToChampions = append(cc.MatchPhysicalDamageDealtToChampions, float64(participant.Stats.PhysicalDamageDealtToChampions))
						cc.MatchPhysicalDamageTaken = append(cc.MatchPhysicalDamageTaken, float64(participant.Stats.PhysicalDamageTaken))
						cc.MatchTrueDamageDealt = append(cc.MatchTrueDamageDealt, float64(participant.Stats.TrueDamageDealt))
						cc.MatchTrueDamageDealtToChampions = append(cc.MatchTrueDamageDealtToChampions, float64(participant.Stats.TrueDamageDealtToChampions))
						cc.MatchTrueDamageTaken = append(cc.MatchTrueDamageTaken, float64(participant.Stats.TrueDamageTaken))
						cc.MatchTotalHeal = append(cc.MatchTotalHeal, float64(participant.Stats.TotalHeal))
						cc.MatchDamageDealtToObjectives = append(cc.MatchDamageDealtToObjectives, float64(participant.Stats.DamageDealtToObjectives))
						cc.MatchDamageDealtToTurrets = append(cc.MatchDamageDealtToTurrets, float64(participant.Stats.DamageDealtToTurrets))
						cc.MatchTimeCCingOthers = append(cc.MatchTimeCCingOthers, float64(participant.Stats.TimeCCingOthers))

						perRole.MatchGoldEarned = append(perRole.MatchGoldEarned, float64(participant.Stats.GoldEarned))
						perRole.MatchTotalMinionsKilled = append(perRole.MatchTotalMinionsKilled, float64(participant.Stats.TotalMinionsKilled))
						perRole.MatchTotalDamageDealt = append(perRole.MatchTotalDamageDealt, float64(participant.Stats.TotalDamageDealt))
						perRole.MatchTotalDamageDealtToChampions = append(perRole.MatchTotalDamageDealtToChampions, float64(participant.Stats.TotalDamageDealtToChampions))
						perRole.MatchTotalDamageTaken = append(perRole.MatchTotalDamageTaken, float64(participant.Stats.TotalDamageTaken))
						perRole.MatchMagicDamageDealt = append(perRole.MatchMagicDamageDealt, float64(participant.Stats.MagicDamageDealt))
						perRole.MatchMagicDamageDealtToChampions = append(perRole.MatchMagicDamageDealtToChampions, float64(participant.Stats.MagicDamageDealtToChampions))
						perRole.MatchPhysicalDamageDealt = append(perRole.MatchPhysicalDamageDealt, float64(participant.Stats.PhysicalDamageDealt))
						perRole.MatchPhysicalDamageDealtToChampions = append(perRole.MatchPhysicalDamageDealtToChampions, float64(participant.Stats.PhysicalDamageDealtToChampions))
						perRole.MatchPhysicalDamageTaken = append(perRole.MatchPhysicalDamageTaken, float64(participant.Stats.PhysicalDamageTaken))
						perRole.MatchTrueDamageDealt = append(perRole.MatchTrueDamageDealt, float64(participant.Stats.TrueDamageDealt))
						perRole.MatchTrueDamageDealtToChampions = append(perRole.MatchTrueDamageDealtToChampions, float64(participant.Stats.TrueDamageDealtToChampions))
						perRole.MatchTrueDamageTaken = append(perRole.MatchTrueDamageTaken, float64(participant.Stats.TrueDamageTaken))
						perRole.MatchTotalHeal = append(perRole.MatchTotalHeal, float64(participant.Stats.TotalHeal))
						perRole.MatchDamageDealtToObjectives = append(perRole.MatchDamageDealtToObjectives, float64(participant.Stats.DamageDealtToObjectives))
						perRole.MatchDamageDealtToTurrets = append(perRole.MatchDamageDealtToTurrets, float64(participant.Stats.DamageDealtToTurrets))
						perRole.MatchTimeCCingOthers = append(perRole.MatchTimeCCingOthers, float64(participant.Stats.TimeCCingOthers))

						perRoleAll.MatchGoldEarned = append(perRoleAll.MatchGoldEarned, float64(participant.Stats.GoldEarned))
						perRoleAll.MatchTotalMinionsKilled = append(perRoleAll.MatchTotalMinionsKilled, float64(participant.Stats.TotalMinionsKilled))
						perRoleAll.MatchTotalDamageDealt = append(perRoleAll.MatchTotalDamageDealt, float64(participant.Stats.TotalDamageDealt))
						perRoleAll.MatchTotalDamageDealtToChampions = append(perRoleAll.MatchTotalDamageDealtToChampions, float64(participant.Stats.TotalDamageDealtToChampions))
						perRoleAll.MatchTotalDamageTaken = append(perRoleAll.MatchTotalDamageTaken, float64(participant.Stats.TotalDamageTaken))
						perRoleAll.MatchMagicDamageDealt = append(perRoleAll.MatchMagicDamageDealt, float64(participant.Stats.MagicDamageDealt))
						perRoleAll.MatchMagicDamageDealtToChampions = append(perRoleAll.MatchMagicDamageDealtToChampions, float64(participant.Stats.MagicDamageDealtToChampions))
						perRoleAll.MatchPhysicalDamageDealt = append(perRoleAll.MatchPhysicalDamageDealt, float64(participant.Stats.PhysicalDamageDealt))
						perRoleAll.MatchPhysicalDamageDealtToChampions = append(perRoleAll.MatchPhysicalDamageDealtToChampions, float64(participant.Stats.PhysicalDamageDealtToChampions))
						perRoleAll.MatchPhysicalDamageTaken = append(perRoleAll.MatchPhysicalDamageTaken, float64(participant.Stats.PhysicalDamageTaken))
						perRoleAll.MatchTrueDamageDealt = append(perRoleAll.MatchTrueDamageDealt, float64(participant.Stats.TrueDamageDealt))
						perRoleAll.MatchTrueDamageDealtToChampions = append(perRoleAll.MatchTrueDamageDealtToChampions, float64(participant.Stats.TrueDamageDealtToChampions))
						perRoleAll.MatchTrueDamageTaken = append(perRoleAll.MatchTrueDamageTaken, float64(participant.Stats.TrueDamageTaken))
						perRoleAll.MatchTotalHeal = append(perRoleAll.MatchTotalHeal, float64(participant.Stats.TotalHeal))
						perRoleAll.MatchDamageDealtToObjectives = append(perRoleAll.MatchDamageDealtToObjectives, float64(participant.Stats.DamageDealtToObjectives))
						perRoleAll.MatchDamageDealtToTurrets = append(perRoleAll.MatchDamageDealtToTurrets, float64(participant.Stats.DamageDealtToTurrets))
						perRoleAll.MatchTimeCCingOthers = append(perRoleAll.MatchTimeCCingOthers, float64(participant.Stats.TimeCCingOthers))

						if participant.Stats.Win {
							perRole.Wins++
							perRoleAll.Wins++
							cc.TotalWins++
							ccall.TotalWins++
						}

						// Backassign structs
						cc.PerRole[lane][role] = perRole
						cct[cid] = cc
						champsCountersPerTier[matchTier] = cct

						ccall.PerRole[lane][role] = perRoleAll
						champsCountersAllTiers[cid] = ccall
					}

					// Champion Bans
					bannedIDs := make(map[int]bool)
					for _, team := range currentMatch.Teams {
						for _, ban := range team.Bans {
							cid := ban.ChampionID
							bannedIDs[cid] = true
						}
					}
					for cid := range bannedIDs {
						// Get structs for counting
						cc := champsCountersPerTier[matchTier][cid]
						ccall := champsCountersAllTiers[cid]

						// Do counts
						cc.TotalBans++
						ccall.TotalBans++

						// Backassign structs
						champsCountersPerTier[matchTier][cid] = cc
						champsCountersAllTiers[cid] = ccall
					}

					cnt++
				}

				// Prepare results for ChampionsStats (ALL tiers)
				for cid, champCounters := range champsCountersAllTiers {
					stats, err := sr.prepareChampionStats(uint64(cid), version[0], version[1], totalGamesForGameVersion, &champCounters)
					stats.Tier = "ALL"
					if err == nil {
						err = sr.storage.StoreChampionStats(stats)
						if err != nil {
							sr.log.Warnf("Something went wrong storing the Champion Stats: %s", err)
						}
					}
				}

				// Prepare results for ChampionsStats (per tier)
				for tier, champsCounters := range champsCountersPerTier {
					for cid, champCounters := range champsCounters {
						stats, err := sr.prepareChampionStats(uint64(cid), version[0], version[1], totalGamesForGameVersionTier[tier], &champCounters)
						stats.Tier = tier
						if err == nil {
							err = sr.storage.StoreChampionStats(stats)
							if err != nil {
								sr.log.Warnf("Something went wrong storing the Champion Stats: %s", err)
							}
						}
					}
				}

				cur.Close()
				sr.log.Debugf("matchAnalysisWorker calculation for Game Version %s done. Analyzed %d matches", gameVersion, cnt)
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

			nextUpdate = time.Minute * time.Duration(sr.config.RScriptsUpdateInterval)
			elapsed := time.Since(start)
			sr.log.Infof("Finished matchAnalysisWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
		}
	}
}

func (sr *StatsRunner) prepareChampionStats(champID uint64, majorVersion uint32, minorVersion uint32, totalGamesForGameVersion uint64, champCounters *championCounters) (*storage.ChampionStats, error) {

	gameVersion := fmt.Sprintf("%d.%d", majorVersion, minorVersion)

	championStats := storage.ChampionStats{}
	championStats.ChampionID = champID
	championStats.GameVersion = gameVersion
	championStats.SampleSize = champCounters.TotalPicks

	championStats.AvgK, championStats.StdDevK = calcMeanStdDev(champCounters.MatchKills, nil)
	if math.IsNaN(championStats.StdDevK) {
		championStats.StdDevK = 0
	}
	championStats.AvgD, championStats.StdDevD = calcMeanStdDev(champCounters.MatchDeaths, nil)
	if math.IsNaN(championStats.StdDevD) {
		championStats.StdDevD = 0
	}
	championStats.AvgA, championStats.StdDevA = calcMeanStdDev(champCounters.MatchAssists, nil)
	if math.IsNaN(championStats.StdDevA) {
		championStats.StdDevA = 0
	}

	championStats.AvgGoldEarned, championStats.StdDevGoldEarned = calcMeanStdDev(champCounters.MatchGoldEarned, nil)
	if math.IsNaN(championStats.StdDevGoldEarned) {
		championStats.StdDevGoldEarned = 0
	}
	championStats.AvgTotalMinionsKilled, championStats.StdDevTotalMinionsKilled = calcMeanStdDev(champCounters.MatchTotalMinionsKilled, nil)
	if math.IsNaN(championStats.StdDevTotalMinionsKilled) {
		championStats.StdDevTotalMinionsKilled = 0
	}
	championStats.AvgTotalDamageDealt, championStats.StdDevTotalDamageDealt = calcMeanStdDev(champCounters.MatchTotalDamageDealt, nil)
	if math.IsNaN(championStats.StdDevTotalDamageDealt) {
		championStats.StdDevTotalDamageDealt = 0
	}
	championStats.AvgTotalDamageDealtToChampions, championStats.StdDevTotalDamageDealtToChampions = calcMeanStdDev(champCounters.MatchTotalDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevTotalDamageDealtToChampions) {
		championStats.StdDevTotalDamageDealtToChampions = 0
	}
	championStats.AvgTotalDamageTaken, championStats.StdDevTotalDamageTaken = calcMeanStdDev(champCounters.MatchTotalDamageTaken, nil)
	if math.IsNaN(championStats.StdDevTotalDamageTaken) {
		championStats.StdDevTotalDamageTaken = 0
	}
	championStats.AvgMagicDamageDealt, championStats.StdDevMagicDamageDealt = calcMeanStdDev(champCounters.MatchMagicDamageDealt, nil)
	if math.IsNaN(championStats.StdDevMagicDamageDealt) {
		championStats.StdDevMagicDamageDealt = 0
	}
	championStats.AvgMagicDamageDealtToChampions, championStats.StdDevMagicDamageDealtToChampions = calcMeanStdDev(champCounters.MatchMagicDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevMagicDamageDealtToChampions) {
		championStats.StdDevMagicDamageDealtToChampions = 0
	}
	championStats.AvgPhysicalDamageDealt, championStats.StdDevPhysicalDamageDealt = calcMeanStdDev(champCounters.MatchPhysicalDamageDealt, nil)
	if math.IsNaN(championStats.StdDevPhysicalDamageDealt) {
		championStats.StdDevPhysicalDamageDealt = 0
	}
	championStats.AvgPhysicalDamageDealtToChampions, championStats.StdDevPhysicalDamageDealtToChampions = calcMeanStdDev(champCounters.MatchPhysicalDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevPhysicalDamageDealtToChampions) {
		championStats.StdDevPhysicalDamageDealtToChampions = 0
	}
	championStats.AvgPhysicalDamageTaken, championStats.StdDevPhysicalDamageTaken = calcMeanStdDev(champCounters.MatchPhysicalDamageTaken, nil)
	if math.IsNaN(championStats.StdDevPhysicalDamageTaken) {
		championStats.StdDevPhysicalDamageTaken = 0
	}
	championStats.AvgTrueDamageDealt, championStats.StdDevTrueDamageDealt = calcMeanStdDev(champCounters.MatchTrueDamageDealt, nil)
	if math.IsNaN(championStats.StdDevTrueDamageDealt) {
		championStats.StdDevTrueDamageDealt = 0
	}
	championStats.AvgTrueDamageDealtToChampions, championStats.StdDevTrueDamageDealtToChampions = calcMeanStdDev(champCounters.MatchTrueDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevTrueDamageDealtToChampions) {
		championStats.StdDevTrueDamageDealtToChampions = 0
	}
	championStats.AvgTrueDamageTaken, championStats.StdDevTrueDamageTaken = calcMeanStdDev(champCounters.MatchTrueDamageTaken, nil)
	if math.IsNaN(championStats.StdDevTrueDamageTaken) {
		championStats.StdDevTrueDamageTaken = 0
	}
	championStats.AvgTotalHeal, championStats.StdDevTotalHeal = calcMeanStdDev(champCounters.MatchTotalHeal, nil)
	if math.IsNaN(championStats.StdDevTotalHeal) {
		championStats.StdDevTotalHeal = 0
	}
	championStats.AvgDamageDealtToObjectives, championStats.StdDevDamageDealtToObjectives = calcMeanStdDev(champCounters.MatchDamageDealtToObjectives, nil)
	if math.IsNaN(championStats.StdDevDamageDealtToObjectives) {
		championStats.StdDevDamageDealtToObjectives = 0
	}
	championStats.AvgDamageDealtToTurrets, championStats.StdDevDamageDealtToTurrets = calcMeanStdDev(champCounters.MatchDamageDealtToTurrets, nil)
	if math.IsNaN(championStats.StdDevDamageDealtToTurrets) {
		championStats.StdDevDamageDealtToTurrets = 0
	}
	championStats.AvgTimeCCingOthers, championStats.StdDevTimeCCingOthers = calcMeanStdDev(champCounters.MatchTimeCCingOthers, nil)
	if math.IsNaN(championStats.StdDevTimeCCingOthers) {
		championStats.StdDevTimeCCingOthers = 0
	}

	var err error
	championStats.MedianK, err = calcMedian(champCounters.MatchKills, nil)
	if err != nil {
		sr.log.Debugf("Error calculating Median for MatchKills: %s", err)
	}
	championStats.MedianD, err = calcMedian(champCounters.MatchDeaths, nil)
	if err != nil {
		sr.log.Debugf("Error calculating Median for MatchDeaths: %s", err)
	}
	championStats.MedianA, err = calcMedian(champCounters.MatchAssists, nil)
	if err != nil {
		sr.log.Debugf("Error calculating Median for MatchAssists: %s", err)
	}

	losses := champCounters.TotalPicks - champCounters.TotalWins
	wins := champCounters.TotalWins

	if losses > 0 {
		championStats.WinLossRatio = float64(wins) / float64(losses)
	} else if losses == 0 && wins > 0 {
		championStats.WinLossRatio = 1.0
	} else {
		championStats.WinLossRatio = 0
	}

	if champCounters.TotalPicks > 0 {
		championStats.WinRate = float64(wins) / float64(champCounters.TotalPicks)
	} else {
		championStats.WinRate = 0
	}

	if totalGamesForGameVersion > 0 {
		championStats.BanRate = float64(champCounters.TotalBans) / float64(totalGamesForGameVersion)
	} else {
		championStats.BanRate = 0
	}

	if totalGamesForGameVersion > 0 {
		championStats.PickRate = float64(champCounters.TotalPicks) / float64(totalGamesForGameVersion)
	} else {
		championStats.PickRate = 0
	}

	topWins := uint64(0)
	topLosses := uint64(0)
	midWins := uint64(0)
	midLosses := uint64(0)
	jungleWins := uint64(0)
	jungleLosses := uint64(0)
	botCarryWins := uint64(0)
	botCarryLosses := uint64(0)
	botSupportWins := uint64(0)
	botSupportLosses := uint64(0)
	botUnknownWins := uint64(0)
	botUnknownLosses := uint64(0)
	unknownWins := uint64(0)
	unknownLosses := uint64(0)

	for key, lane := range champCounters.PerRole {
		if key == "TOP" {
			for _, role := range lane {
				topWins = topWins + role.Wins
				topLosses = topLosses + (role.Picks - role.Wins)
			}
		} else if key == "MIDDLE" {
			for _, role := range lane {
				midWins = midWins + role.Wins
				midLosses = midLosses + (role.Picks - role.Wins)
			}
		} else if key == "JUNGLE" {
			for _, role := range lane {
				jungleWins = jungleWins + role.Wins
				jungleLosses = jungleLosses + (role.Picks - role.Wins)
			}
		} else if key == "BOTTOM" {
			for key, role := range lane {
				if key == "DUO_CARRY" {
					botCarryWins = botCarryWins + role.Wins
					botCarryLosses = botCarryLosses + (role.Picks - role.Wins)
				} else if key == "DUO_SUPPORT" {
					botSupportWins = botSupportWins + role.Wins
					botSupportLosses = botSupportLosses + (role.Picks - role.Wins)
				} else {
					botUnknownWins = botUnknownWins + role.Wins
					botUnknownLosses = botUnknownLosses + (role.Picks - role.Wins)
				}
			}
		} else {
			for _, role := range lane {
				unknownWins = unknownWins + role.Wins
				unknownLosses = unknownLosses + (role.Picks - role.Wins)
			}
		}
	}

	topPercentage := float64(topWins+topLosses) / float64(champCounters.TotalPicks) * 100.0
	midPercentage := float64(midWins+midLosses) / float64(champCounters.TotalPicks) * 100.0
	junglePercentage := float64(jungleWins+jungleLosses) / float64(champCounters.TotalPicks) * 100.0
	botCarryPercentage := float64(botCarryWins+botCarryLosses) / float64(champCounters.TotalPicks) * 100.0
	botSupportPercentage := float64(botSupportWins+botSupportLosses) / float64(champCounters.TotalPicks) * 100.0

	totWithoutUnknownPercentage := topPercentage + midPercentage + junglePercentage + botCarryPercentage + botSupportPercentage

	renormTopPercentage := topPercentage / totWithoutUnknownPercentage * 100.0
	renormMidPercentage := midPercentage / totWithoutUnknownPercentage * 100.0
	renormJunglePercentage := junglePercentage / totWithoutUnknownPercentage * 100.0
	renormBotCarryPercentage := botCarryPercentage / totWithoutUnknownPercentage * 100.0
	renormBotSupportPercentage := botSupportPercentage / totWithoutUnknownPercentage * 100.0

	var roles []string
	if renormTopPercentage > 33 {
		roles = append(roles, "Top")
	}
	if renormMidPercentage > 33 {
		roles = append(roles, "Mid")
	}
	if renormJunglePercentage > 33 {
		roles = append(roles, "Jungle")
	}
	if renormBotCarryPercentage > 33 {
		roles = append(roles, "Carry")
	}
	if renormBotSupportPercentage > 33 {
		roles = append(roles, "Support")
	}

	championStats.Roles = roles

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "TOP",
			Role: "Solo",

			Percentage: float64(topWins+topLosses) / float64(champCounters.TotalPicks) * 100.0,
			Wins:       uint32(topWins),
			NGames:     uint32(topWins + topLosses),
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "MIDDLE",
			Role: "Solo",

			Percentage: float64(midWins+midLosses) / float64(champCounters.TotalPicks) * 100.0,
			Wins:       uint32(midWins),
			NGames:     uint32(midWins + midLosses),
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "JUNGLE",
			Role: "Solo",

			Percentage: float64(jungleWins+jungleLosses) / float64(champCounters.TotalPicks) * 100.0,
			Wins:       uint32(jungleWins),
			NGames:     uint32(jungleWins + jungleLosses),
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "BOT",
			Role: "Carry",

			Percentage: float64(botCarryWins+botCarryLosses) / float64(champCounters.TotalPicks) * 100.0,
			Wins:       uint32(botCarryWins),
			NGames:     uint32(botCarryWins + botCarryLosses),
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "BOT",
			Role: "Support",

			Percentage: float64(botSupportWins+botSupportLosses) / float64(champCounters.TotalPicks) * 100.0,
			Wins:       uint32(botSupportWins),
			NGames:     uint32(botSupportWins + botSupportLosses),
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "BOT",
			Role: "Unknown",

			Percentage: float64(botUnknownWins+botUnknownLosses) / float64(champCounters.TotalPicks) * 100.0,
			Wins:       uint32(botUnknownWins),
			NGames:     uint32(botUnknownWins + botUnknownLosses),
		},
	)

	championStats.LaneRolePercentage = append(championStats.LaneRolePercentage,
		storage.LaneRolePercentage{
			Lane: "UNKNOWN",
			Role: "Unknown",

			Percentage: float64(unknownWins+unknownLosses) / float64(champCounters.TotalPicks) * 100.0,
			Wins:       uint32(unknownWins),
			NGames:     uint32(unknownWins + unknownLosses),
		},
	)

	champions := sr.storage.GetChampions(false)
	for _, val := range champions {
		if val.Key == strconv.FormatUint(champID, 10) {
			championStats.ChampionName = val.Name
			championStats.ChampionRealID = val.ID
			break
		}
	}

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Uknown",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				0.0,
				0.0,
				0.0,
				float64(botUnknownWins+botUnknownLosses) / float64(champCounters.TotalPicks) * 100.0,
				float64(unknownWins+unknownLosses) / float64(champCounters.TotalPicks) * 100.0,
			},
		},
	)

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Support",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				0.0,
				0.0,
				0.0,
				float64(botSupportWins+botSupportLosses) / float64(champCounters.TotalPicks) * 100.0,
				0.0,
			},
		},
	)

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Carry",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				0.0,
				0.0,
				0.0,
				float64(botCarryWins+botCarryLosses) / float64(champCounters.TotalPicks) * 100.0,
				0.0,
			},
		},
	)

	championStats.LaneRolePercentagePlotly = append(championStats.LaneRolePercentagePlotly,
		storage.LaneRolePercentagePlotly{
			Name: "Solo",
			Type: "bar",

			X: []string{
				"TOP", "MIDDLE", "JUNGLE", "BOT", "UNKNOWN",
			},
			Y: []float64{
				float64(topWins+topLosses) / float64(champCounters.TotalPicks) * 100.0,
				float64(midWins+midLosses) / float64(champCounters.TotalPicks) * 100.0,
				float64(jungleWins+jungleLosses) / float64(champCounters.TotalPicks) * 100.0,
				0.0,
				0.0,
			},
		},
	)

	championStats.Timestamp = time.Now()

	return &championStats, nil
}
