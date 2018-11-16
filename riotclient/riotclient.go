package riotclient

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

// Client defines the interface for a Riot API client
type Client interface {
	SummonerByName(name string) (s *Summoner, err error)
	SummonerByAccountID(id uint64) (s *Summoner, err error)
	SummonerBySummonerID(id uint64) (s *Summoner, err error)

	Champions() (s *ChampionList, err error)

	FreeRotation() (*FreeRotation, error)

	MatchByID(id uint64) (s *Match, err error)
}

// RiotClient Riot LoL API client
type RiotClient struct {
	config         config.RiotClient
	httpClient     *http.Client
	log            *logrus.Entry
	rateLimit      rateLimit
	isStarted      bool
	rateLimitMutex sync.Mutex
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
func NewClient(httpClient *http.Client, cfg config.RiotClient) (*RiotClient, error) {
	err := checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	c := &RiotClient{
		config:     cfg,
		httpClient: httpClient,
		log:        logging.Get("RiotClient"),
		isStarted:  false,
	}

	cfg.Region = strings.ToLower(cfg.Region)

	return c, nil
}

// Start starts the riot client and its workers
func (c *RiotClient) Start() {
	if !c.isStarted {
		c.log.Println("Starting Riot Client")
		go c.worker()
		c.isStarted = true
	} else {
		c.log.Println("Riot Client already running")
	}
}

// Stop stops the riot client and its workers
func (c *RiotClient) Stop() {
	if c.isStarted {
		c.log.Println("Stopping Riot Client")
		// TODO
		c.isStarted = false
	} else {
		c.log.Println("Riot Client already stopped")
	}
}

func (c *RiotClient) checkResponseCodeOK(response *http.Response) error {
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

func (c *RiotClient) checkRateLimited(response *http.Response) error {
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
			c.updateRateLimitRetryAt(10)
		}
		return fmt.Errorf("Status code 429 (Rate Limited)")
	}

	return nil
}

func (c *RiotClient) apiCall(path string, method string, body string) (r []byte, e error) {
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

	c.log.Debugln("ApiCall: Pushing request into work queue")
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
	case <-time.After(30 * time.Second):
		c.log.Debugln("ApiCall: API call timed out")
		return nil, fmt.Errorf("Worker timed out")
	}
}
