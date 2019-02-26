package statsrunner

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/torlenor/alolstats/riotclient"

	"github.com/torlenor/alolstats/utils"
)

type roleCounters struct {
	Picks uint64
	Wins  uint64
}

type championCounters struct {
	ChampionID int

	GameVersion string

	TotalPicks uint64
	TotalWins  uint64

	TotalBans uint64

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
				sr.log.Debugf("matchAnalysisWorker calculation for Game Version %s started", gameVersion)

				// Prepare championsCountersPerTier
				champsCountersPerTier := make(championsCountersPerTier)
				champsCountersAllTiers := sr.newChampionsCounters(champions, gameVersion)

				cur, err := sr.storage.GetStoredMatchesCursorByGameVersion(gameVersion)
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

					matchTier := determineMatchTier(currentMatch.Participants)

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
						perRole.Picks++
						perRoleAll.Picks++

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
					for _, team := range currentMatch.Teams {
						for _, ban := range team.Bans {
							cid := ban.ChampionID

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
					}

					cnt++
				}

				sr.log.Debugf("matchAnalysisWorker calculation for Game Version %s done. Analyzed %d matches", gameVersion, cnt)
			}

			nextUpdate = time.Minute * time.Duration(sr.config.RScriptsUpdateInterval)
			elapsed := time.Since(start)
			sr.log.Infof("Finished matchAnalysisWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
		}
	}
}
