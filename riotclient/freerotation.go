package riotclient

import (
	"encoding/json"
	"fmt"
	"time"
)

// FreeRotation contains the data of the free champions rotation from Riot API
type FreeRotation struct {
	FreeChampionIds              []int     `json:"freeChampionIds"`
	FreeChampionIdsForNewPlayers []int     `json:"freeChampionIdsForNewPlayers"`
	MaxNewPlayerLevel            int       `json:"maxNewPlayerLevel"`
	Timestamp                    time.Time `json:"timestamp"`
}

// FreeRotation returns the current free champions rotation from Riot API
func (c *RiotClient) FreeRotation() (*FreeRotation, error) {
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/platform/"+c.config.APIVersion+"/champion-rotations", "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	freeRotation := FreeRotation{}
	err = json.Unmarshal(data, &freeRotation)
	if err != nil {
		return nil, err
	} else if len(freeRotation.FreeChampionIds) == 0 {
		return nil, fmt.Errorf("Received empty free rotation list: %s", data)
	}

	freeRotation.Timestamp = time.Now()

	return &freeRotation, nil
}
