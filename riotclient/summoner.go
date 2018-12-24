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

// SummonerDTO models the Riot API response for summoners endpoints
type SummonerDTO struct {
	AccountID     string    `json:"accountId"`
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ProfileIcon   int       `json:"profileIconId"`
	PuuID         string    `json:"puuid"`
	SummonerLevel int64     `json:"summonerLevel"`
	RevisionDate  int64     `json:"revisionDate"`
	Timestamp     time.Time `json:"timestamp"`
}
