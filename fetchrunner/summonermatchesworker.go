package fetchrunner

import (
	"fmt"
	"strconv"
	"time"
)

func (f *FetchRunner) getChallengerLeagueSummonerAccountIDs(queue string, accountIDs map[uint64]bool) error {
	challengerLeague, err := f.storage.GetChallengerLeagueByQueue(queue)
	if err != nil {
		return fmt.Errorf("Error getting ChallengerLeague: %s", err)
	}

	for _, leagueEntry := range challengerLeague.Entries {
		summonerID, err := strconv.ParseUint(leagueEntry.PlayerOrTeamID, 10, 64)
		if err != nil {
			f.log.Warnf("Could not convert value %s to SummonerID", leagueEntry.PlayerOrTeamID)
			continue
		}
		summoner, err := f.storage.GetSummonerBySummonerID(summonerID)
		if err != nil {
			f.log.Warnf("Could not get Summoner for Summoner ID %d: %s", summonerID, err)
			continue
		}
		accountIDs[summoner.AccountID] = true
	}

	return nil
}

func (f *FetchRunner) getMasterLeagueSummonerAccountIDs(queue string, accountIDs map[uint64]bool) error {
	masterLeague, err := f.storage.GetMasterLeagueByQueue(queue)
	if err != nil {
		return fmt.Errorf("Error getting Master League: %s", err)
	}

	for _, leagueEntry := range masterLeague.Entries {
		summonerID, err := strconv.ParseUint(leagueEntry.PlayerOrTeamID, 10, 64)
		if err != nil {
			f.log.Warnf("Could not convert value %s to SummonerID", leagueEntry.PlayerOrTeamID)
			continue
		}
		summoner, err := f.storage.GetSummonerBySummonerID(summonerID)
		if err != nil {
			f.log.Warnf("Could not get Summoner for Summoner ID %d: %s", summonerID, err)
			continue
		}
		accountIDs[summoner.AccountID] = true
	}

	return nil
}

func (f *FetchRunner) fetchSummonerMatches(accountID uint64) {
	stop := false
	startIndex := uint32(0)
	endIndex := uint32(100)
	for !stop {
		matches, err := f.storage.GetMatchesByAccountID(accountID, startIndex, endIndex)
		if err != nil {
			f.log.Errorf("Error getting the current match list for Summoner: %s", err)
			break
		}
		for _, match := range matches.Matches {
			f.storage.FetchAndStoreMatch(uint64(match.GameID))
		}
		if len(matches.Matches) == 0 || (endIndex+1) >= uint32(matches.TotalGames) {
			stop = true
		}
		if f.config.MatchesForSummonerLastNMatches > 0 && uint64(endIndex+1) > f.config.MatchesForSummonerLastNMatches {
			stop = true
		}
		startIndex += 100
		endIndex += 100
	}
}

func (f *FetchRunner) summonerMatchesWorker() {
	f.workersWG.Add(1)
	defer f.workersWG.Done()

	var nextUpdate time.Duration

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

			if len(f.config.MatchesForSummonerAccountIDs) > 0 {
				f.log.Infof("Fetching matches for specified Account IDs...")
				for _, accountID := range f.config.MatchesForSummonerAccountIDs {
					f.fetchSummonerMatches(accountID)
				}
			}

			if len(f.config.MatchesForChallengerLeagueSummonerQueues) > 0 {
				f.log.Infof("Fetching matches for specified Challenger Leagues...")
				accountIDs := make(map[uint64]bool)

				for _, queue := range f.config.MatchesForChallengerLeagueSummonerQueues {
					err := f.getChallengerLeagueSummonerAccountIDs(queue, accountIDs)
					if err != nil {
						f.log.Errorf("Error fetching Account IDs for Challenger League queue %s: %s", queue, err)
						continue
					}
				}

				f.log.Infof("Found %d unique Account IDs in specified Challenger Leagues. Fetching matches...", len(accountIDs))
				for accountID := range accountIDs {
					f.fetchSummonerMatches(accountID)
				}
			}

			if len(f.config.MatchesForMasterLeagueSummonerQueues) > 0 {
				f.log.Infof("Fetching matches for specified Master Leagues...")
				accountIDs := make(map[uint64]bool)

				for _, queue := range f.config.MatchesForMasterLeagueSummonerQueues {
					err := f.getMasterLeagueSummonerAccountIDs(queue, accountIDs)
					if err != nil {
						f.log.Errorf("Error fetching Account IDs for Master League queue %s: %s", queue, err)
						continue
					}
				}

				f.log.Infof("Found %d unique Account IDs in specified Master Leagues. Fetching matches...", len(accountIDs))
				for accountID := range accountIDs {
					f.fetchSummonerMatches(accountID)
				}
			}

			nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)

			elapsed := time.Since(start)
			f.log.Infof("Finished SummonerMatchesWorker run. Took %s", elapsed)
		}
	}
}
