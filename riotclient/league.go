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

	Timestamp time.Time `json:"timestamp"`
}

// MiniSeriesDTO no description available
type MiniSeriesDTO struct {
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Target   int    `json:"target"`
	Progress string `json:"progress"`
}

// LeagueItemDTO contains a certain entry in a league response
type LeagueItemDTO struct {
	SummonerName string        `json:"summonerName"`
	Wins         int           `json:"wins"`
	Losses       int           `json:"losses"`
	Rank         string        `json:"rank"`
	SummonerID   string        `json:"summonerId"`
	LeaguePoints int           `json:"leaguePoints"`
	MiniSeries   MiniSeriesDTO `json:"miniSeries"`
}

// LeagueListDTO models the response of the Riot API for league informations
type LeagueListDTO struct {
	Tier     string          `json:"tier"`
	Queue    string          `json:"queue"`
	LeagueID string          `json:"leagueId"`
	Name     string          `json:"name"`
	Entries  []LeagueItemDTO `json:"entries"`

	Timestamp time.Time `json:"timestamp"`
}

// LeaguePositionDTO models the response of the Riot API for leagues for a certain summoner (one entry in result array)
type LeaguePositionDTO struct {
	QueueType    string `json:"queueType"`
	SummonerName string `json:"summonerName"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Rank         string `json:"rank"`
	LeagueName   string `json:"leagueName"`
	LeagueID     string `json:"leagueId"`
	Tier         string `json:"tier"`
	SummonerID   string `json:"summonerId"`
	LeaguePoints int    `json:"leaguePoints"`

	Timestamp time.Time `json:"timestamp"`
}

// LeaguePositionDTOList models response of the Riot API for leagues for a certain summoner
type LeaguePositionDTOList struct {
	LeaguePosition []LeaguePositionDTO
}
