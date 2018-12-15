package riotclient

import "time"

// Summoner a summoner account data
type Summoner struct {
	AccountID    uint64    `json:"accountId"`
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	ProfileIcon  uint32    `json:"profileIconId"`
	PuuID        uint64    `json:"puuid"`
	Level        uint64    `json:"summonerLevel"`
	RevisionDate uint64    `json:"revisionDate"`
	Timestamp    time.Time `json:"timestamp"`
}
