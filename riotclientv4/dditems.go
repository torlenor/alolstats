package riotclientv4

import (
	"encoding/json"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// itemData is used for parsing the data coming from data dragon
type itemData struct {
	Type    string               `json:"type"`
	Version string               `json:"version"`
	Basic   riotclient.BasicItem `json:"basic"`
	Data    riotclient.ItemList  `json:"data"`
	Groups  []struct {
		ID              string `json:"id"`
		MaxGroupOwnable string `json:"MaxGroupOwnable"`
	} `json:"groups"`
	Tree []struct {
		Header string   `json:"header"`
		Tags   []string `json:"tags"`
	} `json:"tree"`
}

// Items gets all items from Data Dragon
func (c *RiotClientV4) Items() (*riotclient.ItemList, error) {
	data, err := c.ddragon.GetDataDragonItems()
	if err != nil {
		return nil, fmt.Errorf("Error getting Items from Data Dragon: %s", err)
	}

	return c.parseItems(data)
}

// ItemsSpecificVersionLanguage gets all items for a specific gameVersion and language from Data Dragon
func (c *RiotClientV4) ItemsSpecificVersionLanguage(gameVersion, language string) (*riotclient.ItemList, error) {
	data, err := c.ddragon.GetDataDragonItemsSpecificVersionLanguage(gameVersion, language)
	if err != nil {
		return nil, fmt.Errorf("Error getting Items from Data Dragon: %s", err)
	}

	return c.parseItems(data)
}

func (c *RiotClientV4) parseItems(data []byte) (*riotclient.ItemList, error) {
	itemsDat := itemData{}
	err := json.Unmarshal(data, &itemsDat)
	if err != nil {
		return nil, err
	}

	items := itemsDat.Data

	now := now()

	for key, item := range items {
		item.Key = key
		item.Timestamp = now
		items[key] = item
	}

	return &items, nil
}
