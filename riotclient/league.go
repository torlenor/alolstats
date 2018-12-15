package riotclient

import "time"

// LeagueData is the representation of a Challenger League at a certain point in time
type LeagueData struct {
	Tier     string `json:"tier"`
	Queue    string `json:"queue"`
	LeagueID string `json:"leagueId"`
	Name     string `json:"name"`
	Entries  []struct {
		HotStreak        bool   `json:"hotStreak"`
		Wins             int    `json:"wins"`
		Veteran          bool   `json:"veteran"`
		Losses           int    `json:"losses"`
		Rank             string `json:"rank"`
		PlayerOrTeamName string `json:"playerOrTeamName"`
		Inactive         bool   `json:"inactive"`
		PlayerOrTeamID   string `json:"playerOrTeamId"`
		FreshBlood       bool   `json:"freshBlood"`
		LeaguePoints     int    `json:"leaguePoints"`
	} `json:"entries"`
	Timestamp time.Time
}
