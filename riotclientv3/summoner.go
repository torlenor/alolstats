package riotclientv3

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// SummonerByName gets summoner data by its name
func (c *RiotClientV3) SummonerByName(name string) (s *riotclient.Summoner, err error) {
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/summoner/v3/summoners/by-name/"+name, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := riotclient.Summoner{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, fmt.Errorf("%s. Data was: %s", err, data)
	} else if summoner.ID == 0 {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = time.Now()

	return &summoner, nil
}

// SummonerByAccountID gets summoner data by its AccountID
func (c *RiotClientV3) SummonerByAccountID(id uint64) (s *riotclient.Summoner, err error) {
	idStr := strconv.FormatUint(id, 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/summoner/v3/summoners/by-account/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := riotclient.Summoner{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == 0 {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = time.Now()

	return &summoner, nil
}

// SummonerBySummonerID gets summoner data by its SummonerID
func (c *RiotClientV3) SummonerBySummonerID(id uint64) (s *riotclient.Summoner, err error) {
	idStr := strconv.FormatUint(id, 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/summoner/v3/summoners/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := riotclient.Summoner{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == 0 {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = time.Now()

	return &summoner, nil
}
