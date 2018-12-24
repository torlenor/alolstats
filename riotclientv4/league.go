package riotclientv4

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

func (c *RiotClientV4) leagueByQueue(leagueEndPoint string, queue string) (*riotclient.LeagueListDTO, error) {
	// https://euw1.api.riotgames.com/lol/league/v4/[leagueEndPoint]/by-queue/[QUEUE]
	if queue == "RANKED_SOLO_5x5" || queue == "RANKED_FLEX_SR" || queue == "RANKED_FLEX_TT" {
		data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/league/"+c.config.APIVersion+"/"+leagueEndPoint+"/by-queue/"+queue, "GET", "")
		if err != nil {
			return nil, fmt.Errorf("Error in API call: %s", err)
		}

		league := riotclient.LeagueListDTO{}
		err = json.Unmarshal(data, &league)
		if err != nil {
			return nil, fmt.Errorf("%s. Data was: %s", err, data)
		}

		league.Timestamp = time.Now()

		return &league, nil
	}

	return nil, fmt.Errorf("Invalid queue type %s, allowed are RANKED_SOLO_5x5, RANKED_FLEX_SR or RANKED_FLEX_TT", queue)
}

// LeagueByQueue gets the requested league list by queue type
// Known values for league are "masterleagues", "grandmasterleagues", "challengerleagues"
// Allowed values for queue are "RANKED_SOLO_5x5", "RANKED_FLEX_SR", "RANKED_FLEX_TT"
func (c *RiotClientV4) LeagueByQueue(league string, queue string) (*riotclient.LeagueListDTO, error) {
	// https://euw1.api.riotgames.com/lol/league/v4/masterleagues/by-queue/RANKED_SOLO_5x5
	// https://euw1.api.riotgames.com/lol/league/v4/grandmasterleagues/by-queue/RANKED_SOLO_5x5
	// https://euw1.api.riotgames.com/lol/league/v4/challengerleagues/by-queue/RANKED_SOLO_5x5
	return c.leagueByQueue(league, queue)
}
