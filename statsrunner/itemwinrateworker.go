package statsrunner

import (
	"fmt"
	"time"

	"git.abyle.org/hps/alolstats/riotclient"
	"git.abyle.org/hps/alolstats/statsrunner/analyzer"
	"git.abyle.org/hps/alolstats/storage"
	"git.abyle.org/hps/alolstats/utils"
)

func (sr *StatsRunner) itemWinRateWorker() {
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

					analyzer := analyzer.NewItemAnalyzer(int(gameVersion[0]), int(gameVersion[1]))

					sr.log.Infof("Calculation of itemWinRateWorker for Game Version %s and Queue %s started", gameVersion, queue)

					cur, err := sr.storage.GetMatchesCursorByGameVersionMapQueueID(majorMinor, mapID, queueID)
					if err != nil {
						sr.log.Errorf("Error performing itemWinRateWorker calculation for Game Version %s: %s", gameVersion, err)
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
							sr.log.Errorf("Error deconding match: %s", err)
							continue
						}

						if currentMatch.MapID != 11 || currentMatch.QueueID != int(queueID) {
							sr.log.Warnf("Found match which should not have been returned from storage, skipping...")
							continue
						}

						analyzer.FeedMatch(currentMatch)
						cnt++
					}

					cur.Close()
					sr.log.Infof("Calculation of itemWinRateWorker for Game Version %s and Queue %s done. Analyzed %d matches", gameVersion, queue, cnt)

					_ = analyzer.Analyze()
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

			nextUpdate = time.Minute * time.Duration(sr.config.ChampionsStats.UpdateInverval)
			elapsed := time.Since(start)
			sr.log.Infof("Finished itemWinRateWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
			sr.calculationMutex.Unlock()
		}
	}
}
