package riotclientv4

import (
	"encoding/json"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// Used for parsing the data coming from data dragon
type championData struct {
	Type    string                   `json:"type"`
	Format  string                   `json:"format"`
	Version string                   `json:"version"`
	Data    riotclient.ChampionsList `json:"data"`
}

// Champions gets all champions from Data Dragon
func (c *RiotClientV4) Champions() (s riotclient.ChampionsList, err error) {
	championsData, err := c.ddragon.GetDataDragonChampions()
	if err != nil {
		return nil, fmt.Errorf("Error getting champions from Data Dragon: %s", err)
	}

	championsDat := championData{}
	err = json.Unmarshal(championsData, &championsDat)
	if err != nil {
		return nil, err
	}

	champions := championsDat.Data

	now := now()

	for id, champion := range champions {
		champion.Timestamp = now
		champions[id] = champion
	}

	return champions, nil
}
