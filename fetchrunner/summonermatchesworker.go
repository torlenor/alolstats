package fetchrunner

import "time"

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
			f.log.Printf("Stopping summonerMatchesWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			f.log.Infof("Performing SummonerMatchesWorker run")

			start := time.Now()

			for _, accountID := range f.config.MatchesForSummonerAccountIDs {
				f.fetchSummonerMatches(accountID)
			}

			nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalSummonerMatches)

			elapsed := time.Since(start)
			f.log.Infof("Finished SummonerMatchesWorker run. Took %s", elapsed)
		}
	}
}
