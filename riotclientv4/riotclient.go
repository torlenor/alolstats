// Package riotclientv4 provides the Riot API client for API version v4
package riotclientv4

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"

	riotclientrl "git.abyle.org/hps/alolstats/riotclient/ratelimit"
)

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

var now = time.Now
var apiCall = (*RiotClientV4).realAPICall

// RiotClientV4 Riot LoL API client
type RiotClientV4 struct {
	config      config.RiotClient
	httpClient  httpClient
	log         *logrus.Entry
	isStarted   bool
	workersWG   sync.WaitGroup
	stopWorkers chan struct{}
	workQueue   workQueue

	ddragon   dataDragon
	rateLimit *riotclientrl.RiotClientRL
}

type dataDragon interface {
	GetDataDragonChampions() ([]byte, error)

	GetDataDragonItems() ([]byte, error)
	GetDataDragonItemsSpecificVersionLanguage(gameVersion, language string) ([]byte, error)

	GetDataDragonSummonerSpells() ([]byte, error)
	GetDataDragonSummonerSpellsSpecificVersionLanguage(gameVersion, language string) ([]byte, error)

	GetDataDragonRunesReforged() ([]byte, error)
	GetDataDragonRunesReforgedSpecificVersionLanguage(gameVersion, language string) ([]byte, error)

	GetLoLVersions() ([]byte, error)
}

func checkConfig(cfg config.RiotClient) error {
	if cfg.APIVersion != "v4" {
		return fmt.Errorf("APIVersion is not correct, must be v4")
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
func NewClient(httpClient httpClient, cfg config.RiotClient,
	ddragon dataDragon,
	rateLimit *riotclientrl.RiotClientRL) (*RiotClientV4, error) {
	err := checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("RiotClientV4 [%s]", cfg.Region)
	c := &RiotClientV4{
		config:      cfg,
		httpClient:  httpClient,
		log:         logging.Get(name),
		isStarted:   false,
		workersWG:   sync.WaitGroup{},
		stopWorkers: make(chan struct{}),
		workQueue:   make(workQueue),

		ddragon:   ddragon,
		rateLimit: rateLimit,
	}

	cfg.Region = strings.ToLower(cfg.Region)

	return c, nil
}

// Start starts the riot client and its workers
func (c *RiotClientV4) Start() {
	if !c.isStarted {
		c.log.Println("Starting Riot Client")
		c.stopWorkers = make(chan struct{})
		c.workQueue = make(workQueue)
		c.workersWG.Add(1)
		go c.worker(c.workQueue)
		c.isStarted = true
	} else {
		c.log.Println("Riot Client already running")
	}
}

// Stop stops the riot client and its workers
func (c *RiotClientV4) Stop() {
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
func (c *RiotClientV4) IsRunning() bool {
	return c.isStarted
}

func (c *RiotClientV4) checkResponseCodeOK(response *http.Response) error {
	// Rate limit 429 is handeled in separate check function
	switch response.StatusCode {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("Status code 400 (Bad Request)")
	case 401:
		return fmt.Errorf("Status code 401 (Unauthorized)")
	case 403:
		return fmt.Errorf("Status code 403 (Forbidden)")
	case 404:
		return fmt.Errorf("Status code 404 (Not Found)")
	case 405:
		return fmt.Errorf("Status code 404 (Method not allowed)")
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

func (c *RiotClientV4) checkRateLimited(response *http.Response, method string) error {
	c.rateLimit.UpdateRateLimits(response.Header, method)

	if response.StatusCode == 429 {
		c.log.Warnf("Rate limited with header: %s", response.Header)
		return fmt.Errorf("Status code 429 (Rate Limited)")
	}

	return nil
}

func (c *RiotClientV4) realAPICall(path string, method string, body string) (r []byte, e error) {
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

	c.workQueue <- work

	select {
	case res := <-work.responseChan:
		if res.err != nil {
			c.log.Debugln("ApiCall: Error from worker:", res.err)
			return nil, res.err
		}
		c.log.Debugln("ApiCall: Successfully finished Api Call")
		return ioutil.ReadAll(res.response.Body)
	case <-time.After(180 * time.Second):
		c.log.Debugln("ApiCall: API call timed out")
		return nil, fmt.Errorf("Worker timed out")
	}
}
