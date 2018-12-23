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

type currentVersions struct {
	N struct {
		Item        string `json:"item"`
		Rune        string `json:"rune"`
		Mastery     string `json:"mastery"`
		Summoner    string `json:"summoner"`
		Champion    string `json:"champion"`
		Profileicon string `json:"profileicon"`
		Map         string `json:"map"`
		Language    string `json:"language"`
		Sticker     string `json:"sticker"`
	} `json:"n"`
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
	httpClient *http.Client
	log        *logrus.Entry
}

func checkConfig(cfg config.RiotClient) error {
	if len(cfg.Region) == 0 {
		return fmt.Errorf("Region is empty, check config file")
	}
	return nil
}

// New creates a new Riot LoL API Data Dragon client
func New(httpClient *http.Client, cfg config.RiotClient) (*RiotClientDD, error) {
	err := checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	c := &RiotClientDD{
		config:     cfg,
		httpClient: httpClient,
		log:        logging.Get("RiotClientDD"),
	}

	cfg.Region = strings.ToLower(cfg.Region)

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

	if len(region) > 0 {
		return string(region[:len(region)-1])
	}

	c.log.Errorf("Could not get region from config. Defaulting to euw")
	return "euw"
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

	championsURL := versions.Cdn + "/" + versions.N.Champion + "/data/en_US/champion.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Champions data from Data Dragon: %s", err)
	}

	return body, nil
}
