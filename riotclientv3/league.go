package riotclientv3

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

func (c *RiotClientV3) leagueByQueue(leagueEndPoint string, queue string) (*riotclient.LeagueData, error) {
	// https://euw1.api.riotgames.com/lol/league/v3/LEAGUEENDPOINT/by-queue/RANKED_SOLO_5x5
	if queue == "RANKED_SOLO_5x5" || queue == "RANKED_FLEX_SR" || queue == "RANKED_FLEX_TT" {
		data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/league/v3/"+leagueEndPoint+"/by-queue/"+queue, "GET", "")
		if err != nil {
			return nil, fmt.Errorf("Error in API call: %s", err)
		}

		league := riotclient.LeagueData{}
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
func (c *RiotClientV3) MasterLeagueByQueue(queue string) (*riotclient.LeagueData, error) {
	// https://euw1.api.riotgames.com/lol/league/v3/masterleagues/by-queue/RANKED_SOLO_5x5
	return c.leagueByQueue("masterleagues", queue)
}

// ChallengerLeagueByQueue gets the current challenger league list by queue type
// Allowed values for queue are "RANKED_SOLO_5x5", "RANKED_FLEX_SR", "RANKED_FLEX_TT"
func (c *RiotClientV3) ChallengerLeagueByQueue(queue string) (*riotclient.LeagueData, error) {
	// https://euw1.api.riotgames.com/lol/league/v3/challengerleagues/by-queue/RANKED_SOLO_5x5
	return c.leagueByQueue("challengerleagues", queue)
}
