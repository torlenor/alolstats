package fetchrunner

import (
	"fmt"
	"time"
)

func (f *FetchRunner) getLeagueSummonerAccountIDs(league string, queue string, accountIDs map[string]bool) error {
	leagueData, err := f.storage.GetLeagueByQueue(league, queue)
	if err != nil {
		return fmt.Errorf("Error getting League %s: %s", league, err)
	}

	for _, leagueEntry := range leagueData.Entries {
		summoner, err := f.storage.GetSummonerBySummonerID(leagueEntry.SummonerID)
		if err != nil {
			f.log.Warnf("Could not get Summoner for Summoner ID %s: %s", leagueEntry.SummonerID, err)
			continue
		}
		accountIDs[summoner.AccountID] = true
	}

	return nil
}

func (f *FetchRunner) fetchSummonerMatchesByName(summonerName string, number uint32) {
	summoner, err := f.storage.GetSummonerByName(summonerName)
	if err != nil {
		f.log.Errorf("Error fetching summoner matches: Could not get Summoner Data for Summoner %s", summonerName)
		return
	}
	accountID := summoner.AccountID

	f.fetchSummonerMatchesByAccountID(accountID, number)
}

func (f *FetchRunner) fetchSummonerMatchesByAccountID(accountID string, number uint32) {
	stop := false
	startIndex := uint32(0)
	endIndex := uint32(100)
	if number < endIndex {
		endIndex = number
	}
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

			if len(f.config.FetchMatchesForSummoners) > 0 {
				f.log.Infof("Fetching matches for specified Summoners...")
				for _, summonerName := range f.config.FetchMatchesForSummoners {
					f.fetchSummonerMatchesByName(summonerName, uint32(f.config.FetchMatchesForSummonersNumber))
					if f.shouldWorkersStop {
						elapsed := time.Since(start)
						f.log.Infof("Canceled SummonerMatchesWorker run. Took %s", elapsed)
						nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)
						continue WaitLoop
					}
				}
			}

			accountIDs := make(map[string]bool)
			for _, league := range f.config.FetchMatchesForLeagues {
				if len(f.config.FetchMatchesForLeagueQueues) > 0 {
					f.log.Infof("Fetching matches for specified %s Leagues...", league)

					for _, queue := range f.config.FetchMatchesForLeagueQueues {
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

			f.log.Infof("Found %d unique Account IDs in specified Leagues. Fetching matches...", len(accountIDs))
			for accountID := range accountIDs {
				f.fetchSummonerMatchesByAccountID(accountID, uint32(f.config.FetchMatchesForLeaguesNumber))
				if f.shouldWorkersStop {
					elapsed := time.Since(start)
					f.log.Infof("Canceled SummonerMatchesWorker run. Took %s", elapsed)
					nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)
					continue WaitLoop
				}
			}

			nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)

			elapsed := time.Since(start)
			f.log.Infof("Finished SummonerMatchesWorker run. Took %s", elapsed)
		}
	}
}
