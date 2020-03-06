package riotclientv4

import (
	"encoding/json"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// Used for parsing the data coming from data dragon
type summonerSpellsData struct {
	Type    string                        `json:"type"`
	Version string                        `json:"version"`
	Data    riotclient.SummonerSpellsList `json:"data"`
}

// SummonerSpells gets all items from Data Dragon
func (c *RiotClientV4) SummonerSpells() (*riotclient.SummonerSpellsList, error) {
	data, err := c.ddragon.GetDataDragonSummonerSpells()
	if err != nil {
		return nil, fmt.Errorf("Error getting SummonerSpells from Data Dragon: %s", err)
	}

	return c.parseSummonerSpells(data)
}

// SummonerSpellsSpecificVersionLanguage gets all items for a specific gameVersion and language from Data Dragon
func (c *RiotClientV4) SummonerSpellsSpecificVersionLanguage(gameVersion, language string) (*riotclient.SummonerSpellsList, error) {
	data, err := c.ddragon.GetDataDragonSummonerSpellsSpecificVersionLanguage(gameVersion, language)
	if err != nil {
		return nil, fmt.Errorf("Error getting SummonerSpells for gameversion %s and language %s from Data Dragon: %s", gameVersion, language, err)
	}

	return c.parseSummonerSpells(data)
}

func (c *RiotClientV4) parseSummonerSpells(data []byte) (*riotclient.SummonerSpellsList, error) {
	summonerSpellsDat := summonerSpellsData{}
	err := json.Unmarshal(data, &summonerSpellsDat)
	if err != nil {
		return nil, fmt.Errorf("parseSummonerSpells failed with error: %s, data was: %s", err, data)
	}

	summonerSpells := summonerSpellsDat.Data

	now := now()

	for id, summonerSpell := range summonerSpells {
		summonerSpell.Timestamp = now
		summonerSpells[id] = summonerSpell
	}

	return &summonerSpells, nil
}
