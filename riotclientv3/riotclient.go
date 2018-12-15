// Package riotclientv3 provides the Riot API client for API version v3
package riotclientv3

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
)

// RiotClientV3 Riot LoL API client
type RiotClientV3 struct {
	config         config.RiotClient
	httpClient     *http.Client
	log            *logrus.Entry
	rateLimit      rateLimit
	isStarted      bool
	rateLimitMutex sync.Mutex
	workersWG      sync.WaitGroup
	stopWorkers    chan struct{}
}

func checkConfig(cfg config.RiotClient) error {
	if len(cfg.APIVersion) == 0 {
		return fmt.Errorf("APIVersion is empty, check config file")
	}
	if len(cfg.Key) == 0 {
		return fmt.Errorf("Key is empty, check config file")
	}
	if len(cfg.Region) == 0 {
		return fmt.Errorf("Region is empty, check config file")
	}
	return nil
}

// NewClient creates a new Riot LoL API client
func NewClient(httpClient *http.Client, cfg config.RiotClient) (*RiotClientV3, error) {
	err := checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	c := &RiotClientV3{
		config:      cfg,
		httpClient:  httpClient,
		log:         logging.Get("RiotClientV3"),
		isStarted:   false,
		workersWG:   sync.WaitGroup{},
		stopWorkers: make(chan struct{}),
	}

	cfg.Region = strings.ToLower(cfg.Region)

	return c, nil
}

// Start starts the riot client and its workers
func (c *RiotClientV3) Start() {
	if !c.isStarted {
		c.log.Println("Starting Riot Client")
		c.stopWorkers = make(chan struct{})
		go c.worker()
		c.isStarted = true
	} else {
		c.log.Println("Riot Client already running")
	}
}

// Stop stops the riot client and its workers
func (c *RiotClientV3) Stop() {
	if c.isStarted {
		c.log.Println("Stopping Riot Client")
		close(c.stopWorkers)
		c.workersWG.Wait()
		c.isStarted = false
	} else {
		c.log.Println("Riot Client already stopped")
	}
}

// IsRunning returns if the Riot Client is currently started
func (c *RiotClientV3) IsRunning() bool {
	return c.isStarted
}

func (c *RiotClientV3) checkResponseCodeOK(response *http.Response) error {
	// Rate limit 429 is handeled in separate check function
	switch response.StatusCode {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("Status code 400 (Bad Request)")
	case 401:
		return fmt.Errorf("Status code 401 (Unauthoritzed)")
	case 403:
		return fmt.Errorf("Status code 403 (Forbidden)")
	case 404:
		return fmt.Errorf("Status code 404 (Not Found)")
	case 415:
		return fmt.Errorf("Status code 415 (Unsupported Media Type)")
	case 500:
		return fmt.Errorf("Status code 500 (Internal Server Error)")
	case 503:
		return fmt.Errorf("Status code 503 (Service Unavailable)")
	default:
		return fmt.Errorf("Status code %d (Unknown Error)", response.StatusCode)
	}
}

func (c *RiotClientV3) checkRateLimited(response *http.Response) error {
	if val, ok := response.Header["X-App-Rate-Limit"]; ok {
		if len(val) > 0 {
			c.updateAppRateLimits(val[0])
		}
	}
	if val, ok := response.Header["X-App-Rate-Limit-Count"]; ok {
		if len(val) > 0 {
			c.updateAppRateLimitsCount(val[0])
		}
	}
	if val, ok := response.Header["X-Method-Rate-Limit"]; ok {
		c.log.Warnln("TODO: X-Method-Rate-Limit header not yet processed", val)
	}
	if val, ok := response.Header["X-Method-Rate-Limit-Count"]; ok {
		c.log.Warnln("TODO: X-Method-Rate-Limit-Count header not yet processed", val)
	}
	if response.StatusCode == 429 {
		if val, ok := response.Header["Retry-After"]; ok {
			if len(val) > 0 {
				seconds, err := strconv.ParseUint(val[0], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit retry at seconds", val[0])
					c.updateRateLimitRetryAt(10)
				}
				c.updateRateLimitRetryAt(uint32(seconds))
			}
		} else {
			// https://developer.riotgames.com/rate-limiting.html
			// If the rate limit was enforced by the underlying service to which the request was proxied,
			// rather than the API edge, then the above headers will not be included. In that case, your
			// code cannot use the same mechanism to handle these responses. Instead, your code would simply
			// need to back off for a reasonable amount of time (e.g., 1 second) before trying again the
			// same request.
			c.updateRateLimitRetryAt(2)
		}
		c.log.Warnf("Rate limited with header: %s", response.Header)
		return fmt.Errorf("Status code 429 (Rate Limited)")
	}

	return nil
}

func (c *RiotClientV3) apiCall(path string, method string, body string) (r []byte, e error) {
	if !c.isStarted {
		return nil, fmt.Errorf("Riot Client not started. Start by calling the Start() function")
	}

	c.log.Debugln("ApiCall: Got new Api Call:", path)

	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Riot-Token", c.config.Key)
	req.Header.Add("Content-Type", "application/json")

	work := workOrder{request: req,
		responseChan: make(workResponseChan)}

	workQueue <- work

	c.log.Debugln("ApiCall: Waiting for request to finish processing")
	select {
	case res := <-work.responseChan:
		if res.err != nil {
			c.log.Debugln("ApiCall: Error from worker:", res.err)
			return nil, res.err
		}
		c.log.Debugln("ApiCall: Succesfully finished Api Call")
		return ioutil.ReadAll(res.response.Body)
	case <-time.After(180 * time.Second):
		c.log.Debugln("ApiCall: API call timed out")
		return nil, fmt.Errorf("Worker timed out")
	}
}
