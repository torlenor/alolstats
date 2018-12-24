package riotclient

import "time"

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
