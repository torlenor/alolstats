package statsrunner

import (
	"os/exec"
	"path/filepath"
	"time"
)

func (sr *StatsRunner) rScriptWorker() {
	defer sr.workersWG.Done()

	var nextUpdate time.Duration

	for {
		select {
		case <-sr.stopWorkers:
			sr.log.Printf("Stopping rScriptWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			sr.log.Infof("Performing rScriptWorker run")

			start := time.Now()

			script := sr.config.RScriptPath + string(filepath.Separator) + "champion_stats_from_alolstats_api.R"

			cmd := exec.Command("Rscript", script, "-u", "http://localhost:8000", "-o", sr.config.RPlotsOutputPath, "-v", "8.24")
			sr.log.Printf("Running command and waiting for it to finish...")
			err := cmd.Run()
			if err != nil {
				sr.log.Warnf("Command finished with error: %v", err)
			}

			cmd = exec.Command("Rscript", script, "-u", "http://localhost:8000", "-o", sr.config.RPlotsOutputPath, "-v", "8.23")
			sr.log.Printf("Running command and waiting for it to finish...")
			err = cmd.Run()
			if err != nil {
				sr.log.Warnf("Command finished with error: %v", err)
			}

			nextUpdate = time.Minute * time.Duration(sr.config.RScriptsUpdateInterval)

			elapsed := time.Since(start)
			sr.log.Infof("Finished rScriptWorker run. Took %s. Next run in %s", elapsed, nextUpdate)
		}
	}
}
