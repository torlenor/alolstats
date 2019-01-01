package riotclientv4

import (
	"encoding/json"
	"fmt"

	"github.com/torlenor/alolstats/riotclient"
)

// SummonerByName gets summoner data by its name
func (c *RiotClientV4) SummonerByName(name string) (s *riotclient.SummonerDTO, err error) {
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/summoner/"+c.config.APIVersion+"/summoners/by-name/"+name, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := riotclient.SummonerDTO{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, fmt.Errorf("%s. Data was: %s", err, data)
	} else if summoner.ID == "" {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = now()

	return &summoner, nil
}

// SummonerByAccountID gets summoner data by its AccountID
func (c *RiotClientV4) SummonerByAccountID(accountID string) (s *riotclient.SummonerDTO, err error) {
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/summoner/"+c.config.APIVersion+"/summoners/by-account/"+accountID, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := riotclient.SummonerDTO{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == "" {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = now()

	return &summoner, nil
}

// SummonerBySummonerID gets summoner data by its SummonerID
func (c *RiotClientV4) SummonerBySummonerID(summonerID string) (s *riotclient.SummonerDTO, err error) {
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/summoner/"+c.config.APIVersion+"/summoners/"+summonerID, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := riotclient.SummonerDTO{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == "" {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = now()

	return &summoner, nil
}

// SummonerByPUUID gets summoner data by its PUUID
func (c *RiotClientV4) SummonerByPUUID(PUUID string) (s *riotclient.SummonerDTO, err error) {
	// /lol/summoner/v4/summoners/by-puuid/{encryptedPUUID}
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/summoner/"+c.config.APIVersion+"/summoners/by-puuid/"+PUUID, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := riotclient.SummonerDTO{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == "" {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = now()

	return &summoner, nil
}
