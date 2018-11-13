package riotclient

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Summoner a summoner account data
type Summoner struct {
	AccountID    string    `json:"accountId"`
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	ProfileIcon  int       `json:"profileIconId"`
	PuuID        string    `json:"puuid"`
	Level        int       `json:"summonerLevel"`
	RevisionDate int       `json:"revisionDate"`
	Timestamp    time.Time `json:"timestamp"`
}

// SummonerByName gets summoner data by its name
func (c *RiotClient) SummonerByName(name string) (s *Summoner, err error) {
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/summoner/v4/summoners/by-name/"+name, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := Summoner{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == "" {
		return nil, fmt.Errorf("User does not exist")
	}

	return &summoner, nil
}

// SummonerByAccountID gets summoner data by its AccountID
func (c *RiotClient) SummonerByAccountID(id uint64) (s *Summoner, err error) {
	idStr := strconv.FormatUint(id, 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/summoner/v3/summoners/by-account/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := Summoner{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == "" {
		return nil, fmt.Errorf("User does not exist")
	}

	return &summoner, nil
}

// SummonerBySummonerID gets summoner data by its SummonerID
func (c *RiotClient) SummonerBySummonerID(id uint64) (s *Summoner, err error) {
	idStr := strconv.FormatUint(id, 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/summoner/v3/summoners/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := Summoner{}
	err = json.Unmarshal(data, &summoner)
	if err != nil {
		return nil, err
	} else if summoner.ID == "" {
		return nil, fmt.Errorf("User does not exist")
	}

	return &summoner, nil
}
