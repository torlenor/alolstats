package riotclientdd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
	"github.com/sirupsen/logrus"
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

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloadFile failed, got status code %d", response.StatusCode)
	}

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

// GetLoLVersions returns all currenctly known LoL Versions
func (c *RiotClientDD) GetLoLVersions() ([]byte, error) {
	versionsURL := "https://ddragon.leagueoflegends.com/api/versions.json"

	body, err := c.downloadFile(versionsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Verions data from Data Dragon: %s", err)
	}

	return body, nil
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

// GetDataDragonChampionsSpecificVersionLanguage returns the champions for a given game version and language
func (c *RiotClientDD) GetDataDragonChampionsSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, err
	}

	championsURL := versions.Cdn + "/" + gameVersion + "/data/" + language + "/champion.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Champions data for game version %s and language %s from Data Dragon: %s", gameVersion, language, err)
	}

	return body, nil
}

// GetDataDragonSummonerSpells returns the current summoner spells available for the live game version
func (c *RiotClientDD) GetDataDragonSummonerSpells() ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, err
	}

	championsURL := versions.Cdn + "/" + versions.N.Summoner + "/data/" + versions.L + "/summoner.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Summoner Spells data from Data Dragon: %s", err)
	}

	return body, nil
}

// GetDataDragonSummonerSpellsSpecificVersionLanguage returns the Summoner Spells for a given game version and language
func (c *RiotClientDD) GetDataDragonSummonerSpellsSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, fmt.Errorf("Unable to get versions from data dragon: %s", err)
	}

	championsURL := versions.Cdn + "/" + gameVersion + "/data/" + language + "/summoner.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Summoner Spells data for game version %s and language %s from Data Dragon: %s", gameVersion, language, err)
	}

	return body, nil
}

// GetDataDragonItems returns the current items available for the live game version
func (c *RiotClientDD) GetDataDragonItems() ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, err
	}

	championsURL := versions.Cdn + "/" + versions.N.Item + "/data/" + versions.L + "/item.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Items data from Data Dragon: %s", err)
	}

	return body, nil
}

// GetDataDragonItemsSpecificVersionLanguage returns the Items  for a given game version and language
func (c *RiotClientDD) GetDataDragonItemsSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, err
	}

	championsURL := versions.Cdn + "/" + gameVersion + "/data/" + language + "/item.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Items data for game version %s and language %s from Data Dragon: %s", gameVersion, language, err)
	}

	return body, nil
}

// GetDataDragonRunesReforged returns the current Runes Reforged available for the live game version
func (c *RiotClientDD) GetDataDragonRunesReforged() ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, err
	}

	// We use the Item version as there seems to be no special version for Runes Reforged

	championsURL := versions.Cdn + "/" + versions.N.Item + "/data/" + versions.L + "/runesReforged.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Runes Reforged data from Data Dragon: %s", err)
	}

	return body, nil
}

// GetDataDragonRunesReforgedSpecificVersionLanguage returns the Runes Reforged for a given game version and language
func (c *RiotClientDD) GetDataDragonRunesReforgedSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	versions, err := c.getVersions()
	if err != nil {
		return nil, err
	}

	championsURL := versions.Cdn + "/" + gameVersion + "/data/" + language + "/runesReforged.json"

	body, err := c.downloadFile(championsURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading Runes Reforged data for game version %s and language %s from Data Dragon: %s", gameVersion, language, err)
	}

	return body, nil
}
