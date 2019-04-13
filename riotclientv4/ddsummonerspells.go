package riotclientv4

import (
	"encoding/json"
	"fmt"

	"github.com/torlenor/alolstats/riotclient"
)

// Used for parsing the data coming from data dragon
type summonerSpellsData struct {
	Type    string                        `json:"type"`
	Version string                        `json:"version"`
	Data    riotclient.SummonerSpellsList `json:"data"`
}

// SummonerSpells gets all Summoner Spells from Data Dragon
func (c *RiotClientV4) SummonerSpells() (s *riotclient.SummonerSpellsList, err error) {
	data, err := c.ddragon.GetDataDragonSummonerSpells()
	if err != nil {
		return nil, fmt.Errorf("Error getting Summoner Spells from Data Dragon: %s", err)
	}

	summonerSpellsDat := summonerSpellsData{}
	err = json.Unmarshal(data, &summonerSpellsDat)
	if err != nil {
		return nil, err
	}

	summonerSpells := summonerSpellsDat.Data

	now := now()

	for id, summonerSpell := range summonerSpells {
		summonerSpell.Timestamp = now
		summonerSpells[id] = summonerSpell
	}

	return &summonerSpells, nil
}
