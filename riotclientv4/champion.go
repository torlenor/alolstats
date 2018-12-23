package riotclientv4

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// ChampionRotations returns the current free champions rotation from Riot API
func (c *RiotClientV4) ChampionRotations() (*riotclient.FreeRotation, error) {
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/platform/"+c.config.APIVersion+"/champion-rotations", "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	freeRotation := riotclient.FreeRotation{}
	err = json.Unmarshal(data, &freeRotation)
	if err != nil {
		return nil, err
	} else if len(freeRotation.FreeChampionIds) == 0 {
		return nil, fmt.Errorf("Received empty free rotation list: %s", data)
	}

	freeRotation.Timestamp = time.Now()

	return &freeRotation, nil
}
