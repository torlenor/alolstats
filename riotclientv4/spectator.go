package riotclientv4

import (
	"encoding/json"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// ActiveGameBySummonerID returns the active game (live game) for the given Summoner ID
func (c *RiotClientV4) ActiveGameBySummonerID(summonerID string) (*riotclient.CurrentGameInfoDTO, error) {
	// /lol/spectator/v4/active-games/by-summoner/{encryptedSummonerId}
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/spectator/"+c.config.APIVersion+"/active-games/by-summoner/"+summonerID, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	currentGame := riotclient.CurrentGameInfoDTO{}
	err = json.Unmarshal(data, &currentGame)
	if err != nil {
		return nil, err
	}

	return &currentGame, nil
}

// FeaturedGames returns the currently features games from Riot
func (c *RiotClientV4) FeaturedGames() (*riotclient.FeaturedGamesDTO, error) {
	// /lol/spectator/v4/featured-games
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/spectator/"+c.config.APIVersion+"/featured-games", "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	featuredGames := riotclient.FeaturedGamesDTO{}
	err = json.Unmarshal(data, &featuredGames)
	if err != nil {
		return nil, err
	}

	return &featuredGames, nil
}
