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

type matchCounters struct {
	MatchKills   []uint16
	MatchDeaths  []uint16
	MatchAssists []uint16

	MatchGoldEarned         []uint32
	MatchTotalMinionsKilled []uint32

	MatchTotalDamageDealt               []uint32
	MatchTotalDamageDealtToChampions    []uint32
	MatchTotalDamageTaken               []uint32
	MatchMagicDamageDealt               []uint32
	MatchMagicDamageDealtToChampions    []uint32
	MatchPhysicalDamageDealt            []uint32
	MatchPhysicalDamageDealtToChampions []uint32
	MatchPhysicalDamageTaken            []uint32
	MatchTrueDamageDealt                []uint32
	MatchTrueDamageDealtToChampions     []uint32
	MatchTrueDamageTaken                []uint32

	MatchTotalHeal []uint32

	MatchDamageDealtToObjectives []uint32
	MatchDamageDealtToTurrets    []uint32
	MatchTimeCCingOthers         []uint32
}

type roleCounters struct {
	Picks uint64
	Wins  uint64

	Kills   uint64
	Deaths  uint64
	Assists uint64

	matchCounters
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

	matchCounters

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
						cc.MatchKills = append(cc.MatchKills, uint16(participant.Stats.Kills))
						cc.TotalDeaths = cc.TotalDeaths + uint64(participant.Stats.Deaths)
						cc.MatchDeaths = append(cc.MatchDeaths, uint16(participant.Stats.Deaths))
						cc.TotalAssists = cc.TotalAssists + uint64(participant.Stats.Assists)
						cc.MatchAssists = append(cc.MatchAssists, uint16(participant.Stats.Assists))

						ccall.TotalKills = ccall.TotalKills + uint64(participant.Stats.Kills)
						ccall.MatchKills = append(ccall.MatchKills, uint16(participant.Stats.Kills))
						ccall.TotalDeaths = ccall.TotalDeaths + uint64(participant.Stats.Deaths)
						ccall.MatchDeaths = append(ccall.MatchDeaths, uint16(participant.Stats.Deaths))
						ccall.TotalAssists = ccall.TotalAssists + uint64(participant.Stats.Assists)
						ccall.MatchAssists = append(ccall.MatchAssists, uint16(participant.Stats.Assists))

						perRole.Picks++
						perRoleAll.Picks++

						perRole.Kills = perRole.Kills + uint64(participant.Stats.Kills)
						perRole.MatchKills = append(perRole.MatchKills, uint16(participant.Stats.Kills))
						perRole.Deaths = perRole.Deaths + uint64(participant.Stats.Deaths)
						perRole.MatchDeaths = append(perRole.MatchDeaths, uint16(participant.Stats.Deaths))
						perRole.Assists = perRole.Assists + uint64(participant.Stats.Assists)
						perRole.MatchAssists = append(perRole.MatchAssists, uint16(participant.Stats.Assists))

						perRoleAll.Kills = perRoleAll.Kills + uint64(participant.Stats.Kills)
						perRoleAll.MatchKills = append(perRoleAll.MatchKills, uint16(participant.Stats.Kills))
						perRoleAll.Deaths = perRoleAll.Deaths + uint64(participant.Stats.Deaths)
						perRoleAll.MatchDeaths = append(perRoleAll.MatchDeaths, uint16(participant.Stats.Deaths))
						perRoleAll.Assists = perRoleAll.Assists + uint64(participant.Stats.Assists)
						perRoleAll.MatchAssists = append(perRoleAll.MatchAssists, uint16(participant.Stats.Assists))

						ccall.MatchGoldEarned = append(ccall.MatchGoldEarned, uint32(participant.Stats.GoldEarned))
						ccall.MatchTotalMinionsKilled = append(ccall.MatchTotalMinionsKilled, uint32(participant.Stats.TotalMinionsKilled))
						ccall.MatchTotalDamageDealt = append(ccall.MatchTotalDamageDealt, uint32(participant.Stats.TotalDamageDealt))
						ccall.MatchTotalDamageDealtToChampions = append(ccall.MatchTotalDamageDealtToChampions, uint32(participant.Stats.TotalDamageDealtToChampions))
						ccall.MatchTotalDamageTaken = append(ccall.MatchTotalDamageTaken, uint32(participant.Stats.TotalDamageTaken))
						ccall.MatchMagicDamageDealt = append(ccall.MatchMagicDamageDealt, uint32(participant.Stats.MagicDamageDealt))
						ccall.MatchMagicDamageDealtToChampions = append(ccall.MatchMagicDamageDealtToChampions, uint32(participant.Stats.MagicDamageDealtToChampions))
						ccall.MatchPhysicalDamageDealt = append(ccall.MatchPhysicalDamageDealt, uint32(participant.Stats.PhysicalDamageDealt))
						ccall.MatchPhysicalDamageDealtToChampions = append(ccall.MatchPhysicalDamageDealtToChampions, uint32(participant.Stats.PhysicalDamageDealtToChampions))
						ccall.MatchPhysicalDamageTaken = append(ccall.MatchPhysicalDamageTaken, uint32(participant.Stats.PhysicalDamageTaken))
						ccall.MatchTrueDamageDealt = append(ccall.MatchTrueDamageDealt, uint32(participant.Stats.TrueDamageDealt))
						ccall.MatchTrueDamageDealtToChampions = append(ccall.MatchTrueDamageDealtToChampions, uint32(participant.Stats.TrueDamageDealtToChampions))
						ccall.MatchTrueDamageTaken = append(ccall.MatchTrueDamageTaken, uint32(participant.Stats.TrueDamageTaken))
						ccall.MatchTotalHeal = append(ccall.MatchTotalHeal, uint32(participant.Stats.TotalHeal))
						ccall.MatchDamageDealtToObjectives = append(perRole.MatchDamageDealtToObjectives, uint32(participant.Stats.DamageDealtToObjectives))
						ccall.MatchDamageDealtToTurrets = append(perRole.MatchDamageDealtToTurrets, uint32(participant.Stats.DamageDealtToTurrets))
						ccall.MatchTimeCCingOthers = append(perRole.MatchTimeCCingOthers, uint32(participant.Stats.TimeCCingOthers))

						cc.MatchGoldEarned = append(cc.MatchGoldEarned, uint32(participant.Stats.GoldEarned))
						cc.MatchTotalMinionsKilled = append(cc.MatchTotalMinionsKilled, uint32(participant.Stats.TotalMinionsKilled))
						cc.MatchTotalDamageDealt = append(cc.MatchTotalDamageDealt, uint32(participant.Stats.TotalDamageDealt))
						cc.MatchTotalDamageDealtToChampions = append(cc.MatchTotalDamageDealtToChampions, uint32(participant.Stats.TotalDamageDealtToChampions))
						cc.MatchTotalDamageTaken = append(cc.MatchTotalDamageTaken, uint32(participant.Stats.TotalDamageTaken))
						cc.MatchMagicDamageDealt = append(cc.MatchMagicDamageDealt, uint32(participant.Stats.MagicDamageDealt))
						cc.MatchMagicDamageDealtToChampions = append(cc.MatchMagicDamageDealtToChampions, uint32(participant.Stats.MagicDamageDealtToChampions))
						cc.MatchPhysicalDamageDealt = append(cc.MatchPhysicalDamageDealt, uint32(participant.Stats.PhysicalDamageDealt))
						cc.MatchPhysicalDamageDealtToChampions = append(cc.MatchPhysicalDamageDealtToChampions, uint32(participant.Stats.PhysicalDamageDealtToChampions))
						cc.MatchPhysicalDamageTaken = append(cc.MatchPhysicalDamageTaken, uint32(participant.Stats.PhysicalDamageTaken))
						cc.MatchTrueDamageDealt = append(cc.MatchTrueDamageDealt, uint32(participant.Stats.TrueDamageDealt))
						cc.MatchTrueDamageDealtToChampions = append(cc.MatchTrueDamageDealtToChampions, uint32(participant.Stats.TrueDamageDealtToChampions))
						cc.MatchTrueDamageTaken = append(cc.MatchTrueDamageTaken, uint32(participant.Stats.TrueDamageTaken))
						cc.MatchTotalHeal = append(cc.MatchTotalHeal, uint32(participant.Stats.TotalHeal))
						cc.MatchDamageDealtToObjectives = append(cc.MatchDamageDealtToObjectives, uint32(participant.Stats.DamageDealtToObjectives))
						cc.MatchDamageDealtToTurrets = append(cc.MatchDamageDealtToTurrets, uint32(participant.Stats.DamageDealtToTurrets))
						cc.MatchTimeCCingOthers = append(cc.MatchTimeCCingOthers, uint32(participant.Stats.TimeCCingOthers))

						perRole.MatchGoldEarned = append(perRole.MatchGoldEarned, uint32(participant.Stats.GoldEarned))
						perRole.MatchTotalMinionsKilled = append(perRole.MatchTotalMinionsKilled, uint32(participant.Stats.TotalMinionsKilled))
						perRole.MatchTotalDamageDealt = append(perRole.MatchTotalDamageDealt, uint32(participant.Stats.TotalDamageDealt))
						perRole.MatchTotalDamageDealtToChampions = append(perRole.MatchTotalDamageDealtToChampions, uint32(participant.Stats.TotalDamageDealtToChampions))
						perRole.MatchTotalDamageTaken = append(perRole.MatchTotalDamageTaken, uint32(participant.Stats.TotalDamageTaken))
						perRole.MatchMagicDamageDealt = append(perRole.MatchMagicDamageDealt, uint32(participant.Stats.MagicDamageDealt))
						perRole.MatchMagicDamageDealtToChampions = append(perRole.MatchMagicDamageDealtToChampions, uint32(participant.Stats.MagicDamageDealtToChampions))
						perRole.MatchPhysicalDamageDealt = append(perRole.MatchPhysicalDamageDealt, uint32(participant.Stats.PhysicalDamageDealt))
						perRole.MatchPhysicalDamageDealtToChampions = append(perRole.MatchPhysicalDamageDealtToChampions, uint32(participant.Stats.PhysicalDamageDealtToChampions))
						perRole.MatchPhysicalDamageTaken = append(perRole.MatchPhysicalDamageTaken, uint32(participant.Stats.PhysicalDamageTaken))
						perRole.MatchTrueDamageDealt = append(perRole.MatchTrueDamageDealt, uint32(participant.Stats.TrueDamageDealt))
						perRole.MatchTrueDamageDealtToChampions = append(perRole.MatchTrueDamageDealtToChampions, uint32(participant.Stats.TrueDamageDealtToChampions))
						perRole.MatchTrueDamageTaken = append(perRole.MatchTrueDamageTaken, uint32(participant.Stats.TrueDamageTaken))
						perRole.MatchTotalHeal = append(perRole.MatchTotalHeal, uint32(participant.Stats.TotalHeal))
						perRole.MatchDamageDealtToObjectives = append(perRole.MatchDamageDealtToObjectives, uint32(participant.Stats.DamageDealtToObjectives))
						perRole.MatchDamageDealtToTurrets = append(perRole.MatchDamageDealtToTurrets, uint32(participant.Stats.DamageDealtToTurrets))
						perRole.MatchTimeCCingOthers = append(perRole.MatchTimeCCingOthers, uint32(participant.Stats.TimeCCingOthers))

						perRoleAll.MatchGoldEarned = append(perRoleAll.MatchGoldEarned, uint32(participant.Stats.GoldEarned))
						perRoleAll.MatchTotalMinionsKilled = append(perRoleAll.MatchTotalMinionsKilled, uint32(participant.Stats.TotalMinionsKilled))
						perRoleAll.MatchTotalDamageDealt = append(perRoleAll.MatchTotalDamageDealt, uint32(participant.Stats.TotalDamageDealt))
						perRoleAll.MatchTotalDamageDealtToChampions = append(perRoleAll.MatchTotalDamageDealtToChampions, uint32(participant.Stats.TotalDamageDealtToChampions))
						perRoleAll.MatchTotalDamageTaken = append(perRoleAll.MatchTotalDamageTaken, uint32(participant.Stats.TotalDamageTaken))
						perRoleAll.MatchMagicDamageDealt = append(perRoleAll.MatchMagicDamageDealt, uint32(participant.Stats.MagicDamageDealt))
						perRoleAll.MatchMagicDamageDealtToChampions = append(perRoleAll.MatchMagicDamageDealtToChampions, uint32(participant.Stats.MagicDamageDealtToChampions))
						perRoleAll.MatchPhysicalDamageDealt = append(perRoleAll.MatchPhysicalDamageDealt, uint32(participant.Stats.PhysicalDamageDealt))
						perRoleAll.MatchPhysicalDamageDealtToChampions = append(perRoleAll.MatchPhysicalDamageDealtToChampions, uint32(participant.Stats.PhysicalDamageDealtToChampions))
						perRoleAll.MatchPhysicalDamageTaken = append(perRoleAll.MatchPhysicalDamageTaken, uint32(participant.Stats.PhysicalDamageTaken))
						perRoleAll.MatchTrueDamageDealt = append(perRoleAll.MatchTrueDamageDealt, uint32(participant.Stats.TrueDamageDealt))
						perRoleAll.MatchTrueDamageDealtToChampions = append(perRoleAll.MatchTrueDamageDealtToChampions, uint32(participant.Stats.TrueDamageDealtToChampions))
						perRoleAll.MatchTrueDamageTaken = append(perRoleAll.MatchTrueDamageTaken, uint32(participant.Stats.TrueDamageTaken))
						perRoleAll.MatchTotalHeal = append(perRoleAll.MatchTotalHeal, uint32(participant.Stats.TotalHeal))
						perRoleAll.MatchDamageDealtToObjectives = append(perRoleAll.MatchDamageDealtToObjectives, uint32(participant.Stats.DamageDealtToObjectives))
						perRoleAll.MatchDamageDealtToTurrets = append(perRoleAll.MatchDamageDealtToTurrets, uint32(participant.Stats.DamageDealtToTurrets))
						perRoleAll.MatchTimeCCingOthers = append(perRoleAll.MatchTimeCCingOthers, uint32(participant.Stats.TimeCCingOthers))

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

			type leagues struct {
				Leagues []string `json:"leagues"`
			}
			leas := leagues{Leagues: []string{"All", "Master", "Diamond", "Platinum", "Gold", "Silver", "Bronze"}}
			for _, gameVersion := range gameVersions.Versions {
				for _, tier := range leas.Leagues {
					statsSummary, err := sr.generateChampionsSummary(gameVersion, strings.ToUpper(tier))
					if err != nil {
						sr.log.Errorf("Error generating statistics summary: %s", err)
						continue
					}
					sr.storage.StoreChampionStatsSummary(statsSummary)
				}
			}

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

	championStats.AvgK, championStats.StdDevK = calcMeanStdDevUint16(champCounters.MatchKills, nil)
	if math.IsNaN(championStats.StdDevK) {
		championStats.StdDevK = 0
	}
	championStats.AvgD, championStats.StdDevD = calcMeanStdDevUint16(champCounters.MatchDeaths, nil)
	if math.IsNaN(championStats.StdDevD) {
		championStats.StdDevD = 0
	}
	championStats.AvgA, championStats.StdDevA = calcMeanStdDevUint16(champCounters.MatchAssists, nil)
	if math.IsNaN(championStats.StdDevA) {
		championStats.StdDevA = 0
	}

	championStats.AvgGoldEarned, championStats.StdDevGoldEarned = calcMeanStdDevUint32(champCounters.MatchGoldEarned, nil)
	if math.IsNaN(championStats.StdDevGoldEarned) {
		championStats.StdDevGoldEarned = 0
	}
	championStats.AvgTotalMinionsKilled, championStats.StdDevTotalMinionsKilled = calcMeanStdDevUint32(champCounters.MatchTotalMinionsKilled, nil)
	if math.IsNaN(championStats.StdDevTotalMinionsKilled) {
		championStats.StdDevTotalMinionsKilled = 0
	}
	championStats.AvgTotalDamageDealt, championStats.StdDevTotalDamageDealt = calcMeanStdDevUint32(champCounters.MatchTotalDamageDealt, nil)
	if math.IsNaN(championStats.StdDevTotalDamageDealt) {
		championStats.StdDevTotalDamageDealt = 0
	}
	championStats.AvgTotalDamageDealtToChampions, championStats.StdDevTotalDamageDealtToChampions = calcMeanStdDevUint32(champCounters.MatchTotalDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevTotalDamageDealtToChampions) {
		championStats.StdDevTotalDamageDealtToChampions = 0
	}
	championStats.AvgTotalDamageTaken, championStats.StdDevTotalDamageTaken = calcMeanStdDevUint32(champCounters.MatchTotalDamageTaken, nil)
	if math.IsNaN(championStats.StdDevTotalDamageTaken) {
		championStats.StdDevTotalDamageTaken = 0
	}
	championStats.AvgMagicDamageDealt, championStats.StdDevMagicDamageDealt = calcMeanStdDevUint32(champCounters.MatchMagicDamageDealt, nil)
	if math.IsNaN(championStats.StdDevMagicDamageDealt) {
		championStats.StdDevMagicDamageDealt = 0
	}
	championStats.AvgMagicDamageDealtToChampions, championStats.StdDevMagicDamageDealtToChampions = calcMeanStdDevUint32(champCounters.MatchMagicDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevMagicDamageDealtToChampions) {
		championStats.StdDevMagicDamageDealtToChampions = 0
	}
	championStats.AvgPhysicalDamageDealt, championStats.StdDevPhysicalDamageDealt = calcMeanStdDevUint32(champCounters.MatchPhysicalDamageDealt, nil)
	if math.IsNaN(championStats.StdDevPhysicalDamageDealt) {
		championStats.StdDevPhysicalDamageDealt = 0
	}
	championStats.AvgPhysicalDamageDealtToChampions, championStats.StdDevPhysicalDamageDealtToChampions = calcMeanStdDevUint32(champCounters.MatchPhysicalDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevPhysicalDamageDealtToChampions) {
		championStats.StdDevPhysicalDamageDealtToChampions = 0
	}
	championStats.AvgPhysicalDamageTaken, championStats.StdDevPhysicalDamageTaken = calcMeanStdDevUint32(champCounters.MatchPhysicalDamageTaken, nil)
	if math.IsNaN(championStats.StdDevPhysicalDamageTaken) {
		championStats.StdDevPhysicalDamageTaken = 0
	}
	championStats.AvgTrueDamageDealt, championStats.StdDevTrueDamageDealt = calcMeanStdDevUint32(champCounters.MatchTrueDamageDealt, nil)
	if math.IsNaN(championStats.StdDevTrueDamageDealt) {
		championStats.StdDevTrueDamageDealt = 0
	}
	championStats.AvgTrueDamageDealtToChampions, championStats.StdDevTrueDamageDealtToChampions = calcMeanStdDevUint32(champCounters.MatchTrueDamageDealtToChampions, nil)
	if math.IsNaN(championStats.StdDevTrueDamageDealtToChampions) {
		championStats.StdDevTrueDamageDealtToChampions = 0
	}
	championStats.AvgTrueDamageTaken, championStats.StdDevTrueDamageTaken = calcMeanStdDevUint32(champCounters.MatchTrueDamageTaken, nil)
	if math.IsNaN(championStats.StdDevTrueDamageTaken) {
		championStats.StdDevTrueDamageTaken = 0
	}
	championStats.AvgTotalHeal, championStats.StdDevTotalHeal = calcMeanStdDevUint32(champCounters.MatchTotalHeal, nil)
	if math.IsNaN(championStats.StdDevTotalHeal) {
		championStats.StdDevTotalHeal = 0
	}
	championStats.AvgDamageDealtToObjectives, championStats.StdDevDamageDealtToObjectives = calcMeanStdDevUint32(champCounters.MatchDamageDealtToObjectives, nil)
	if math.IsNaN(championStats.StdDevDamageDealtToObjectives) {
		championStats.StdDevDamageDealtToObjectives = 0
	}
	championStats.AvgDamageDealtToTurrets, championStats.StdDevDamageDealtToTurrets = calcMeanStdDevUint32(champCounters.MatchDamageDealtToTurrets, nil)
	if math.IsNaN(championStats.StdDevDamageDealtToTurrets) {
		championStats.StdDevDamageDealtToTurrets = 0
	}
	championStats.AvgTimeCCingOthers, championStats.StdDevTimeCCingOthers = calcMeanStdDevUint32(champCounters.MatchTimeCCingOthers, nil)
	if math.IsNaN(championStats.StdDevTimeCCingOthers) {
		championStats.StdDevTimeCCingOthers = 0
	}

	championStats.MedianK, _ = calcMedianUint16(champCounters.MatchKills, nil)
	championStats.MedianD, _ = calcMedianUint16(champCounters.MatchDeaths, nil)
	championStats.MedianA, _ = calcMedianUint16(champCounters.MatchAssists, nil)

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

	// Role determination
	var roles []string
	if renormTopPercentage > sr.config.RoleThreshold {
		roles = append(roles, "Top")
	}
	if renormMidPercentage > sr.config.RoleThreshold {
		roles = append(roles, "Mid")
	}
	if renormJunglePercentage > sr.config.RoleThreshold {
		roles = append(roles, "Jungle")
	}
	if renormBotCarryPercentage > sr.config.RoleThreshold {
		roles = append(roles, "Carry")
	}
	if renormBotSupportPercentage > sr.config.RoleThreshold {
		roles = append(roles, "Support")
	}
	championStats.Roles = roles

	// Calculation of stats per Role
	championStats.StatsPerRole = make(map[string]storage.StatsValues)
	for _, role := range []string{"Top", "Mid", "Jungle", "Carry", "Support"} {
		statsValues := sr.calculateRoleStats(champCounters, role)
		championStats.StatsPerRole[role] = statsValues
	}

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

func (sr *StatsRunner) calculateRoleStats(champCounters *championCounters, role string) storage.StatsValues {
	summedCounters := roleCounters{}

	switch role {
	case "Top":
		for _, cnters := range champCounters.PerRole["TOP"] {
			sumCounters(&summedCounters, cnters)
		}
	case "Mid":
		for _, cnters := range champCounters.PerRole["MIDDLE"] {
			sumCounters(&summedCounters, cnters)
		}
	case "Jungle":
		for _, cnters := range champCounters.PerRole["JUNGLE"] {
			sumCounters(&summedCounters, cnters)
		}
	case "Carry":
		if lane, ok := champCounters.PerRole["BOTTOM"]; ok {
			if cnters, ok := lane["DUO_CARRY"]; ok {
				sumCounters(&summedCounters, cnters)
			}
		}
	case "Support":
		if lane, ok := champCounters.PerRole["BOTTOM"]; ok {
			if cnters, ok := lane["DUO_SUPPORT"]; ok {
				sumCounters(&summedCounters, cnters)
			}
		}
	}

	return sr.calcStatsFromCounters(&summedCounters)
}
func (sr *StatsRunner) calcStatsFromCounters(counters *roleCounters) storage.StatsValues {
	statsValues := storage.StatsValues{}

	statsValues.SampleSize = counters.Picks

	if (counters.Picks) > 1 {
		statsValues.AvgK, statsValues.StdDevK = calcMeanStdDevUint16(counters.MatchKills, nil)
		if math.IsNaN(statsValues.StdDevK) {
			statsValues.StdDevK = 0
		}
		statsValues.AvgD, statsValues.StdDevD = calcMeanStdDevUint16(counters.MatchDeaths, nil)
		if math.IsNaN(statsValues.StdDevD) {
			statsValues.StdDevD = 0
		}
		statsValues.AvgA, statsValues.StdDevA = calcMeanStdDevUint16(counters.MatchAssists, nil)
		if math.IsNaN(statsValues.StdDevA) {
			statsValues.StdDevA = 0
		}

		statsValues.AvgGoldEarned, statsValues.StdDevGoldEarned = calcMeanStdDevUint32(counters.MatchGoldEarned, nil)
		if math.IsNaN(statsValues.StdDevGoldEarned) {
			statsValues.StdDevGoldEarned = 0
		}
		statsValues.AvgTotalMinionsKilled, statsValues.StdDevTotalMinionsKilled = calcMeanStdDevUint32(counters.MatchTotalMinionsKilled, nil)
		if math.IsNaN(statsValues.StdDevTotalMinionsKilled) {
			statsValues.StdDevTotalMinionsKilled = 0
		}
		statsValues.AvgTotalDamageDealt, statsValues.StdDevTotalDamageDealt = calcMeanStdDevUint32(counters.MatchTotalDamageDealt, nil)
		if math.IsNaN(statsValues.StdDevTotalDamageDealt) {
			statsValues.StdDevTotalDamageDealt = 0
		}
		statsValues.AvgTotalDamageDealtToChampions, statsValues.StdDevTotalDamageDealtToChampions = calcMeanStdDevUint32(counters.MatchTotalDamageDealtToChampions, nil)
		if math.IsNaN(statsValues.StdDevTotalDamageDealtToChampions) {
			statsValues.StdDevTotalDamageDealtToChampions = 0
		}
		statsValues.AvgTotalDamageTaken, statsValues.StdDevTotalDamageTaken = calcMeanStdDevUint32(counters.MatchTotalDamageTaken, nil)
		if math.IsNaN(statsValues.StdDevTotalDamageTaken) {
			statsValues.StdDevTotalDamageTaken = 0
		}
		statsValues.AvgMagicDamageDealt, statsValues.StdDevMagicDamageDealt = calcMeanStdDevUint32(counters.MatchMagicDamageDealt, nil)
		if math.IsNaN(statsValues.StdDevMagicDamageDealt) {
			statsValues.StdDevMagicDamageDealt = 0
		}
		statsValues.AvgMagicDamageDealtToChampions, statsValues.StdDevMagicDamageDealtToChampions = calcMeanStdDevUint32(counters.MatchMagicDamageDealtToChampions, nil)
		if math.IsNaN(statsValues.StdDevMagicDamageDealtToChampions) {
			statsValues.StdDevMagicDamageDealtToChampions = 0
		}
		statsValues.AvgPhysicalDamageDealt, statsValues.StdDevPhysicalDamageDealt = calcMeanStdDevUint32(counters.MatchPhysicalDamageDealt, nil)
		if math.IsNaN(statsValues.StdDevPhysicalDamageDealt) {
			statsValues.StdDevPhysicalDamageDealt = 0
		}
		statsValues.AvgPhysicalDamageDealtToChampions, statsValues.StdDevPhysicalDamageDealtToChampions = calcMeanStdDevUint32(counters.MatchPhysicalDamageDealtToChampions, nil)
		if math.IsNaN(statsValues.StdDevPhysicalDamageDealtToChampions) {
			statsValues.StdDevPhysicalDamageDealtToChampions = 0
		}
		statsValues.AvgPhysicalDamageTaken, statsValues.StdDevPhysicalDamageTaken = calcMeanStdDevUint32(counters.MatchPhysicalDamageTaken, nil)
		if math.IsNaN(statsValues.StdDevPhysicalDamageTaken) {
			statsValues.StdDevPhysicalDamageTaken = 0
		}
		statsValues.AvgTrueDamageDealt, statsValues.StdDevTrueDamageDealt = calcMeanStdDevUint32(counters.MatchTrueDamageDealt, nil)
		if math.IsNaN(statsValues.StdDevTrueDamageDealt) {
			statsValues.StdDevTrueDamageDealt = 0
		}
		statsValues.AvgTrueDamageDealtToChampions, statsValues.StdDevTrueDamageDealtToChampions = calcMeanStdDevUint32(counters.MatchTrueDamageDealtToChampions, nil)
		if math.IsNaN(statsValues.StdDevTrueDamageDealtToChampions) {
			statsValues.StdDevTrueDamageDealtToChampions = 0
		}
		statsValues.AvgTrueDamageTaken, statsValues.StdDevTrueDamageTaken = calcMeanStdDevUint32(counters.MatchTrueDamageTaken, nil)
		if math.IsNaN(statsValues.StdDevTrueDamageTaken) {
			statsValues.StdDevTrueDamageTaken = 0
		}
		statsValues.AvgTotalHeal, statsValues.StdDevTotalHeal = calcMeanStdDevUint32(counters.MatchTotalHeal, nil)
		if math.IsNaN(statsValues.StdDevTotalHeal) {
			statsValues.StdDevTotalHeal = 0
		}
		statsValues.AvgDamageDealtToObjectives, statsValues.StdDevDamageDealtToObjectives = calcMeanStdDevUint32(counters.MatchDamageDealtToObjectives, nil)
		if math.IsNaN(statsValues.StdDevDamageDealtToObjectives) {
			statsValues.StdDevDamageDealtToObjectives = 0
		}
		statsValues.AvgDamageDealtToTurrets, statsValues.StdDevDamageDealtToTurrets = calcMeanStdDevUint32(counters.MatchDamageDealtToTurrets, nil)
		if math.IsNaN(statsValues.StdDevDamageDealtToTurrets) {
			statsValues.StdDevDamageDealtToTurrets = 0
		}
		statsValues.AvgTimeCCingOthers, statsValues.StdDevTimeCCingOthers = calcMeanStdDevUint32(counters.MatchTimeCCingOthers, nil)
		if math.IsNaN(statsValues.StdDevTimeCCingOthers) {
			statsValues.StdDevTimeCCingOthers = 0
		}
	}

	statsValues.MedianK, _ = calcMedianUint16(counters.MatchKills, nil)
	statsValues.MedianD, _ = calcMedianUint16(counters.MatchDeaths, nil)
	statsValues.MedianA, _ = calcMedianUint16(counters.MatchAssists, nil)

	wins := counters.Wins

	if counters.Picks > 0 {
		statsValues.WinRate = float64(wins) / float64(counters.Picks)
	} else {
		statsValues.WinRate = 0
	}

	return statsValues
}

func sumCounters(summedCounters *roleCounters, countersToAdd roleCounters) {
	summedCounters.Picks += countersToAdd.Picks
	summedCounters.Wins += countersToAdd.Wins

	summedCounters.Kills += countersToAdd.Kills
	summedCounters.Deaths += countersToAdd.Deaths
	summedCounters.Assists += countersToAdd.Assists

	summedCounters.MatchKills = append(summedCounters.MatchKills, countersToAdd.MatchKills...)
	summedCounters.MatchDeaths = append(summedCounters.MatchDeaths, countersToAdd.MatchDeaths...)
	summedCounters.MatchAssists = append(summedCounters.MatchAssists, countersToAdd.MatchAssists...)

	summedCounters.MatchGoldEarned = append(summedCounters.MatchGoldEarned, countersToAdd.MatchGoldEarned...)
	summedCounters.MatchTotalMinionsKilled = append(summedCounters.MatchTotalMinionsKilled, countersToAdd.MatchTotalMinionsKilled...)

	summedCounters.MatchTotalDamageDealt = append(summedCounters.MatchTotalDamageDealt, countersToAdd.MatchTotalDamageDealt...)
	summedCounters.MatchTotalDamageDealtToChampions = append(summedCounters.MatchTotalDamageDealtToChampions, countersToAdd.MatchTotalDamageDealtToChampions...)
	summedCounters.MatchTotalDamageTaken = append(summedCounters.MatchTotalDamageTaken, countersToAdd.MatchTotalDamageTaken...)
	summedCounters.MatchMagicDamageDealt = append(summedCounters.MatchMagicDamageDealt, countersToAdd.MatchMagicDamageDealt...)
	summedCounters.MatchMagicDamageDealtToChampions = append(summedCounters.MatchMagicDamageDealtToChampions, countersToAdd.MatchMagicDamageDealtToChampions...)
	summedCounters.MatchPhysicalDamageDealt = append(summedCounters.MatchPhysicalDamageDealt, countersToAdd.MatchPhysicalDamageDealt...)
	summedCounters.MatchPhysicalDamageDealtToChampions = append(summedCounters.MatchPhysicalDamageDealtToChampions, countersToAdd.MatchPhysicalDamageDealtToChampions...)
	summedCounters.MatchPhysicalDamageTaken = append(summedCounters.MatchPhysicalDamageTaken, countersToAdd.MatchPhysicalDamageTaken...)
	summedCounters.MatchTrueDamageDealt = append(summedCounters.MatchTrueDamageDealt, countersToAdd.MatchTrueDamageDealt...)
	summedCounters.MatchTrueDamageDealtToChampions = append(summedCounters.MatchTrueDamageDealtToChampions, countersToAdd.MatchTrueDamageDealtToChampions...)
	summedCounters.MatchTrueDamageTaken = append(summedCounters.MatchTrueDamageTaken, countersToAdd.MatchTrueDamageTaken...)

	summedCounters.MatchTotalHeal = append(summedCounters.MatchTotalHeal, countersToAdd.MatchTotalHeal...)

	summedCounters.MatchDamageDealtToObjectives = append(summedCounters.MatchDamageDealtToObjectives, countersToAdd.MatchDamageDealtToObjectives...)
	summedCounters.MatchDamageDealtToTurrets = append(summedCounters.MatchDamageDealtToTurrets, countersToAdd.MatchDamageDealtToTurrets...)
	summedCounters.MatchTimeCCingOthers = append(summedCounters.MatchTimeCCingOthers, countersToAdd.MatchTimeCCingOthers...)
}
