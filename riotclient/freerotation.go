package riotclient

import "time"

// FreeRotation contains the data of the free champions rotation from Riot API
type FreeRotation struct {
	FreeChampionIds              []int `json:"freeChampionIds"`
	FreeChampionIdsForNewPlayers []int `json:"freeChampionIdsForNewPlayers"`
	MaxNewPlayerLevel            int   `json:"maxNewPlayerLevel"`

	Timestamp time.Time `json:"timestamp"`
}
