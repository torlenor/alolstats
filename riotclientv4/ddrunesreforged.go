package riotclientv4

import (
	"encoding/json"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// RunesReforged gets all items from Data Dragon
func (c *RiotClientV4) RunesReforged() (*riotclient.RunesReforgedList, error) {
	data, err := c.ddragon.GetDataDragonRunesReforged()
	if err != nil {
		return nil, fmt.Errorf("Error getting RunesReforged from Data Dragon: %s", err)
	}

	return c.parseRunesReforged(data)
}

// RunesReforgedSpecificVersionLanguage gets all items for a specific gameVersion and language from Data Dragon
func (c *RiotClientV4) RunesReforgedSpecificVersionLanguage(gameVersion, language string) (*riotclient.RunesReforgedList, error) {
	data, err := c.ddragon.GetDataDragonRunesReforgedSpecificVersionLanguage(gameVersion, language)
	if err != nil {
		return nil, fmt.Errorf("Error getting RunesReforged from Data Dragon: %s", err)
	}

	return c.parseRunesReforged(data)
}

func (c *RiotClientV4) parseRunesReforged(data []byte) (*riotclient.RunesReforgedList, error) {
	runesReforged := riotclient.RunesReforgedList{}
	err := json.Unmarshal(data, &runesReforged)
	if err != nil {
		return nil, err
	}

	now := now()

	for _, item := range runesReforged {
		item.Timestamp = now
	}

	return &runesReforged, nil
}
