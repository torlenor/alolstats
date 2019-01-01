package riotclientv4

import (
	"net/http"
	"time"
)

type workResponseData struct {
	response *http.Response
	err      error
}

type workResponseChan chan workResponseData

type workOrder struct {
	request      *http.Request
	responseChan workResponseChan
	method       string
}

// A buffered channel that we can send work requests on.
type workQueue = chan workOrder

func (c *RiotClientV4) worker(workQueue workQueue) {
	c.log.Debugln("Worker: Starting")

	c.workersWG.Add(1)
	defer c.workersWG.Done()

	for {
		select {
		case work := <-workQueue:
			tryAgain := true
			tries := 0
			var response *http.Response
			var err error
			for tryAgain && tries <= 3 {
				tryAgain = false
				tries++

				sleepFor := time.Until(c.rateLimit.GetRateLimitRetryAt(work.method))
				if sleepFor > 0 {
					c.log.Debugln("Worker: Sleeping for", sleepFor.String(), "to adhere to rate limit")
					time.Sleep(sleepFor)
				}

				response, err = c.httpClient.Do(work.request)
				if err != nil {
					continue
				}
				err = c.checkRateLimited(response, work.method)
				if err != nil {
					tryAgain = true
					c.log.Debugln("Worker: Repeating request")
					continue
				}
				err = c.checkResponseCodeOK(response)
				if err != nil {
					continue
				}
			}

			work.responseChan <- workResponseData{response: response,
				err: err}

			c.log.Debugln("Worker: Done processing work order")
		case <-c.stopWorkers:
			c.log.Printf("Stopping worker")
			return
		}
	}
}
