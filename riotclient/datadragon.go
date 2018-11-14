package riotclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
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

func (c *RiotClient) downloadFile(url string) ([]byte, error) {
	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

func (c *RiotClient) getRegion() string {
	region := strings.ToLower(c.config.Region)

	if len(region) > 0 {
		return string(region[:len(region)-1])
	}

	c.log.Errorf("Could not get region from config. Defaulting to euw")
	return "euw"
}

func (c *RiotClient) getVersions() (*currentVersions, error) {

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

func (c *RiotClient) getDataDragonChampions() ([]byte, error) {
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
