package riotclientv4

import (
	"encoding/json"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// ChampionRotations returns the current free champions rotation from Riot API
func (c *RiotClientV4) ChampionRotations() (*riotclient.FreeRotation, error) {
	// still v3
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/platform/v3/champion-rotations", "GET", "")
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

	freeRotation.Timestamp = now()

	return &freeRotation, nil
}
