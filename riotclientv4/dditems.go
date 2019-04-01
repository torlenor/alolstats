package riotclientv4

import (
	"encoding/json"
	"fmt"

	"github.com/torlenor/alolstats/riotclient"
)

// Items gets all items for a specific gameVersion and language from Data Dragon
func (c *RiotClientV4) Items(gameVersion, language string) (*riotclient.ItemList, error) {
	itemsData, err := c.ddragon.GetDataDragonItemsSpecificVersionLanguage(gameVersion, language)
	if err != nil {
		return nil, fmt.Errorf("Error getting Items from Data Dragon: %s", err)
	}

	itemsDat := riotclient.ItemList{}
	err = json.Unmarshal(itemsData, &itemsDat)
	if err != nil {
		return nil, err
	}

	return &itemsDat, nil
}
