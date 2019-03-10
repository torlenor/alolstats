package fetchrunner

import (
	"fmt"
	"strings"
	"time"

	"github.com/torlenor/alolstats/utils"
)

func (f *FetchRunner) getLeagueSummonerAccountIDs(league string, queue string, accountIDs map[string]bool) error {
	leagueData, err := f.storage.GetRegionalLeagueByQueue(f.config.Region, league, queue)
	if err != nil {
		return fmt.Errorf("Error getting League %s: %s", league, err)
	}

	for _, leagueEntry := range leagueData.Entries {
		summoner, err := f.storage.GetRegionalSummonerBySummonerID(f.config.Region, leagueEntry.SummonerID, false)
		if err != nil {
			f.log.Warnf("Could not get Summoner for Summoner ID %s: %s", leagueEntry.SummonerID, err)
			continue
		}
		accountIDs[summoner.AccountID] = true
	}

	return nil
}

func (f *FetchRunner) fetchSummonerMatchesByName(summonerName string, number uint32, seenAccountIDs map[string]bool, knownLatestVersion string) {
	summoner, err := f.storage.GetRegionalSummonerByName(f.config.Region, summonerName, false)
	if err != nil {
		f.log.Errorf("Error fetching summoner matches: Could not get Summoner Data for Summoner %s", summonerName)
		return
	}
	accountID := summoner.AccountID

	f.fetchSummonerMatchesByAccountID(accountID, number, seenAccountIDs, knownLatestVersion)
}

func (f *FetchRunner) fetchSummonerMatchesByAccountID(accountID string, number uint32, seenAccountIDs map[string]bool, knownLatestVersion string) {
	stop := false
	startIndex := uint32(0)
	endIndex := uint32(100)
	if number < endIndex {
		endIndex = number
	}
	for !stop {
		matches, err := f.storage.GetRegionalMatchesByAccountID(f.config.Region, accountID, startIndex, endIndex)
		if err != nil {
			f.log.Errorf("Error getting the current match list for Summoner: %s", err)
			break
		}
		for _, matchInfo := range matches.Matches {
			match, err := f.storage.RegionalFetchAndStoreMatch(f.config.Region, uint64(matchInfo.GameID))
			if match != nil && err == nil && seenAccountIDs != nil {
				for _, participant := range match.ParticipantIdentities {
					if participant.Player.AccountID != accountID {
						seenAccountIDs[participant.Player.AccountID] = true
					}
				}
			}
			if f.config.FetchOnlyLatestGameVersion && match != nil && err == nil {
				if !f.checkGameVersionsEqual(match.GameVersion, knownLatestVersion) {
					f.log.Debugf("Skipping remaining matches for Summoner %s because we encountered a game version not beeing the latest (latest: %s, seen %s)", accountID, knownLatestVersion, match.GameVersion)
					return
				}
			}
		}
		if len(matches.Matches) == 0 || (endIndex+1) >= uint32(matches.TotalGames) {
			stop = true
		}
		if number > 0 && uint32(endIndex+1) > number {
			stop = true
		}
		startIndex += 100
		endIndex += 100
		if endIndex > number {
			endIndex = number
		}
	}
}

func (f *FetchRunner) checkGameVersionsEqual(latestSeenGameVersion string, knownLatestVersion string) bool {
	version, err := utils.SplitNumericMatchVersion(latestSeenGameVersion)
	if err != nil {
		f.log.Warnf("Could not extract Match Game Version from string %s: %s", latestSeenGameVersion, err)
		return true
	}
	seenVersionStr := fmt.Sprintf("%d.%d", version[0], version[1])
	if err == nil {
		if strings.Compare(seenVersionStr, knownLatestVersion) != 0 {
			return false
		}
	}
	return true
}

func (f *FetchRunner) summonerMatchesWorker() {
	f.workersWG.Add(1)
	defer f.workersWG.Done()

	var nextUpdate time.Duration

WaitLoop:
	for {
		select {
		case <-f.stopWorkers:
			f.log.Printf("Stopping SummonerMatchesWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			f.log.Infof("Performing SummonerMatchesWorker run")

			start := time.Now()

			var knownLatestVersion string
			if f.config.FetchOnlyLatestGameVersion {
				versions, err := utils.SplitNumericVersion(f.config.LatestGameVersionForFetching)
				if err != nil {
					f.log.Warnf("LatestGameVersionForFetching specified in config is invalid, disabling FetchOnlyLatestGameVersion, err was: %s", err)
					f.config.FetchOnlyLatestGameVersion = false
				} else {
					knownLatestVersion = fmt.Sprintf("%d.%d", versions[0], versions[1])
					f.log.Infof("Fetching only for latest known game version %s", knownLatestVersion)
				}
			}

			additionalAccountIDs := make(map[string]bool)
			if len(f.config.FetchMatchesForSummoners) > 0 {
				f.log.Infof("Fetching matches for specified Summoners")
				for _, summonerName := range f.config.FetchMatchesForSummoners {
					if f.shouldWorkersStop {
						elapsed := time.Since(start)
						f.log.Infof("Canceled SummonerMatchesWorker run. Took %s", elapsed)
						nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)
						continue WaitLoop
					}
					f.fetchSummonerMatchesByName(summonerName, uint32(f.config.FetchMatchesForSummonersNumber), additionalAccountIDs, knownLatestVersion)
				}
			}

			f.log.Infof("Fetching matches for specified Leagues")
			accountIDs := make(map[string]bool)
			for _, league := range f.config.FetchMatchesForLeagues {
				if len(f.config.FetchMatchesForLeagueQueues) > 0 {

					for _, queue := range f.config.FetchMatchesForLeagueQueues {
						f.log.Infof("Getting Summoner Account IDs for League %s and Queue %s", league, queue)
						err := f.getLeagueSummonerAccountIDs(league, queue, accountIDs)
						if err != nil {
							f.log.Errorf("Error fetching Account IDs for league %s queue %s: %s", league, queue, err)
							continue
						}
						if f.shouldWorkersStop {
							elapsed := time.Since(start)
							f.log.Infof("Canceled SummonerMatchesWorker run. Took %s", elapsed)
							nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)
							continue WaitLoop
						}
					}
				}
			}

			f.log.Infof("Found %d unique Account IDs in specified Leagues. Fetching matches", len(accountIDs))
			for accountID := range accountIDs {
				if f.shouldWorkersStop {
					elapsed := time.Since(start)
					f.log.Infof("Canceled SummonerMatchesWorker run. Took %s", elapsed)
					nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)
					continue WaitLoop
				}
				f.fetchSummonerMatchesByAccountID(accountID, uint32(f.config.FetchMatchesForLeaguesNumber), additionalAccountIDs, knownLatestVersion)
			}

			for accountID := range accountIDs {
				delete(additionalAccountIDs, accountID)
			}

			if f.config.FetchMatchesForSeenSummoners {
				f.log.Infof("Found %d additional unique Account IDs in fetched matches. Fetching matches", len(additionalAccountIDs))
				for accountID := range additionalAccountIDs {
					if f.shouldWorkersStop {
						elapsed := time.Since(start)
						f.log.Infof("Canceled SummonerMatchesWorker run. Took %s", elapsed)
						nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)
						continue WaitLoop
					}
					f.fetchSummonerMatchesByAccountID(accountID, uint32(f.config.FetchMatchesForLeaguesNumber), nil, knownLatestVersion)
				}
			}

			nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)

			elapsed := time.Since(start)
			f.log.Infof("Finished SummonerMatchesWorker run. Took %s", elapsed)
		}
	}
}
