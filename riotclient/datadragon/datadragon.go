package riotclientdd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
)

// N holds the actual versions
type N struct {
	Item        string `json:"item"`
	Rune        string `json:"rune"`
	Mastery     string `json:"mastery"`
	Summoner    string `json:"summoner"`
	Champion    string `json:"champion"`
	Profileicon string `json:"profileicon"`
	Map         string `json:"map"`
	Language    string `json:"language"`
	Sticker     string `json:"sticker"`
}

type currentVersions struct {
	N              N           `json:"n"`
	V              string      `json:"v"`
	L              string      `json:"l"`
	Cdn            string      `json:"cdn"`
	Dd             string      `json:"dd"`
	Lg             string      `json:"lg"`
	CSS            string      `json:"css"`
	Profileiconmax int         `json:"profileiconmax"`
	Store          interface{} `json:"store"`
}

// RiotClientDD Riot LoL API DataDragon client
type RiotClientDD struct {
	config     config.RiotClient
	httpClient httpClient
	log        *logrus.Entry
}

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

func checkConfig(cfg config.RiotClient) error {
	if len(cfg.Region) < 2 {
		return fmt.Errorf("Region does not comply to Riot region conventions, check config file")
	}
	return nil
}

// New creates a new Riot LoL API Data Dragon client
func New(httpClient httpClient, cfg config.RiotClient) (*RiotClientDD, error) {
	err := checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Region = strings.ToLower(cfg.Region)

	c := &RiotClientDD{
		config:     cfg,
		httpClient: httpClient,
		log:        logging.Get("RiotClientDD"),
	}

	return c, nil
}

func (c *RiotClientDD) downloadFile(url string) ([]byte, error) {
	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

func (c *RiotClientDD) getRegion() string {
	region := strings.ToLower(c.config.Region)

	return string(region[:len(region)-1])
}

func (c *RiotClientDD) getVersions() (*currentVersions, error) {

	versionURL := "https://ddragon.leagueoflegends.com/realms/" + c.getRegion() + ".json"

	versionData, err := c.downloadFile(versionURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading versions data from Data Dragon: %s", err)
	}

	versions := currentVersions{}
	err = json.Unmarshal(versionData, &versions)
	if err != nil {
		return nil, err
	}

	return &versions, nil
}

// GetDataDragonChampions returns the current champions available for the live game version
func (c *RiotClientDD) GetDataDragonChampions() ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, err
	}

	championsURL := versions.Cdn + "/" + versions.N.Champion + "/data/" + versions.L + "/champion.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Champions data from Data Dragon: %s", err)
	}

	return body, nil
}