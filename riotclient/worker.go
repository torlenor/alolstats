package riotclient

import (
	"fmt"
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
}

// A buffered channel that we can send work requests on.
var workQueue = make(chan workOrder)

func (c *RiotClient) worker() {
	c.log.Debugln("Worker: Starting")

	c.workersWG.Add(1)
	defer c.workersWG.Done()

	for {
		select {
		case work := <-workQueue:
			c.log.Debugln("Worker: Got new work order to process")
			sleepUntil := time.Until(c.getRateLimitRetryAt())
			c.log.Debugln("Worker: Sleeping until", sleepUntil.String())
			time.Sleep(sleepUntil)
			additionalSleep := c.getAdditionalWaitTime()
			if additionalSleep > 0 {
				c.log.Debugln("Worker: Sleeping additional", additionalSleep.Seconds(), "seconds to make sure we are not getting rate limited")
				time.Sleep(additionalSleep)
			}

			c.log.Debugln("Worker: Performing http request")
			response, err := c.httpClient.Do(work.request)
			if err != nil {
				work.responseChan <- workResponseData{response: nil,
					err: err}
				continue
			}
			c.log.Debugln("Worker: Checking for rate limit")
			err = c.checkRateLimited(response)
			if err != nil {
				work.responseChan <- workResponseData{response: nil,
					err: fmt.Errorf("Got rate limited from API: %s", err)}
				continue
			}
			c.log.Debugln("Worker: Checking response")
			err = c.checkResponseCodeOK(response)
			if err != nil {
				work.responseChan <- workResponseData{response: nil,
					err: fmt.Errorf("Got an invalid response code from API: %s", err)}
				continue
			}
			c.log.Debugln("Worker: Send back result")
			work.responseChan <- workResponseData{response: response,
				err: nil}
			c.log.Debugln("Worker: Done processing work order")
		case <-c.stopWorkers:
			c.log.Printf("Stopping worker")
			return
		}
	}
}
