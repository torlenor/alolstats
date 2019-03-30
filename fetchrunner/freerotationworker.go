package fetchrunner

import (
	"time"
)

func (f *FetchRunner) freeRotationWorker() {
	f.workersWG.Add(1)
	defer f.workersWG.Done()

	var nextUpdate time.Duration

	for {
		select {
		case <-f.stopWorkers:
			f.log.Printf("Stopping FreeRotationWorker")
			return
		default:
			if nextUpdate > 0 {
				time.Sleep(time.Second * 1)
				nextUpdate -= 1 * time.Second
				continue
			}
			f.log.Infof("Performing FreeRotationWorker run")

			start := time.Now()

			f.storage.GetFreeRotation(true)

			nextUpdate = time.Minute * time.Duration(f.config.UpdateIntervalFreeRotation)

			elapsed := time.Since(start)
			f.log.Infof("Finished FreeRotationWorker run. Took %s", elapsed)
		}
	}
}
