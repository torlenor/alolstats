package riotclient

import (
	"encoding/json"
	"fmt"
	"time"
)

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

func (c *RiotClient) leagueByQueue(leagueEndPoint string, queue string) (*LeagueData, error) {
	// https://euw1.api.riotgames.com/lol/league/v3/LEAGUEENDPOINT/by-queue/RANKED_SOLO_5x5
	if queue == "RANKED_SOLO_5x5" || queue == "RANKED_FLEX_SR" || queue == "RANKED_FLEX_TT" {
		data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/league/v3/"+leagueEndPoint+"/by-queue/"+queue, "GET", "")
		if err != nil {
			return nil, fmt.Errorf("Error in API call: %s", err)
		}

		league := LeagueData{}
		err = json.Unmarshal(data, &league)
		if err != nil {
			return nil, fmt.Errorf("%s. Data was: %s", err, data)
		}

		league.Timestamp = time.Now()

		return &league, nil
	}

	return nil, fmt.Errorf("Invalid queue type %s, allowed are RANKED_SOLO_5x5, RANKED_FLEX_SR or RANKED_FLEX_TT", queue)
}

// MasterLeagueByQueue gets the current challenger league list by queue type
// Allowed values for queue are "RANKED_SOLO_5x5", "RANKED_FLEX_SR", "RANKED_FLEX_TT"
func (c *RiotClient) MasterLeagueByQueue(queue string) (*LeagueData, error) {
	// https://euw1.api.riotgames.com/lol/league/v3/masterleagues/by-queue/RANKED_SOLO_5x5
	return c.leagueByQueue("masterleagues", queue)
}

// ChallengerLeagueByQueue gets the current challenger league list by queue type
// Allowed values for queue are "RANKED_SOLO_5x5", "RANKED_FLEX_SR", "RANKED_FLEX_TT"
func (c *RiotClient) ChallengerLeagueByQueue(queue string) (*LeagueData, error) {
	// https://euw1.api.riotgames.com/lol/league/v3/challengerleagues/by-queue/RANKED_SOLO_5x5
	return c.leagueByQueue("challengerleagues", queue)
}
