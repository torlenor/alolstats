package riotclientv3

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// Used for parsing the data coming from data dragon
type championData struct {
	Type    string                         `json:"type"`
	Format  string                         `json:"format"`
	Version string                         `json:"version"`
	Data    map[string]riotclient.Champion `json:"data"`
}

// Champions gets all champions
func (c *RiotClientV3) Champions() (s *riotclient.ChampionList, err error) {
	championsData, err := c.getDataDragonChampions()
	if err != nil {
		return nil, fmt.Errorf("Error getting champions from Data Dragon: %s", err)
	}

	championsDat := championData{}
	err = json.Unmarshal(championsData, &championsDat)
	if err != nil {
		return nil, err
	}

	champions := riotclient.ChampionList{
		Champions: championsDat.Data,
	}

	now := time.Now()

	for id, champion := range champions.Champions {
		champion.Timestamp = now
		champions.Champions[id] = champion
	}

	return &champions, nil
}
