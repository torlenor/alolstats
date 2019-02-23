package statsrunner

import (
	"fmt"
	"strconv"
	"time"

	"github.com/torlenor/alolstats/storage"
	"github.com/torlenor/alolstats/utils"
)

func (sr *StatsRunner) champStatsCalcWorker() {
	sr.workersWG.Add(1)
	defer sr.workersWG.Done()

	var nextUpdate time.Duration

	for {
		select {
		case <-sr.stopWorkers:
			sr.log.Printf("Stopping champStatsCalcWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			sr.log.Infof("Performing champStatsCalcWorker run")

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

				sr.log.Debugf("champStatsCalcWorker calculation for Game Version %s started", gameVersion)

				for _, champ := range champions {
					if sr.shouldWorkersStop {
						return
					}
					champKey, err := strconv.ParseUint(champ.Key, 10, 32)
					if err != nil {
						sr.log.Warnf("Could not convert value %s to Champion Key", champ.Key)
						return
					}
					stats, err := sr.getChampionStatsByID(champKey, version[0], version[1])
					if err == nil {
						err := sr.storage.StoreChampionStats(stats)
						if err != nil {
							sr.log.Warnf("Something went wrong storing the Champion Stats: %s", err)
						}
					}
				}
				sr.log.Debugf("champStatsCalcWorker calculation for Game Version %s done", gameVersion)
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
			sr.log.Infof("Finished champStatsCalcWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
		}
	}
}
