package riotclient

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Summoner a summoner account data
type Summoner struct {
	AccountID    uint64    `json:"accountId"`
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	ProfileIcon  uint32    `json:"profileIconId"`
	PuuID        uint64    `json:"puuid"`
	Level        uint64    `json:"summonerLevel"`
	RevisionDate uint64    `json:"revisionDate"`
	Timestamp    time.Time `json:"timestamp"`
}

// SummonerByName gets summoner data by its name
func (c *RiotClient) SummonerByName(name string) (s *Summoner, err error) {
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/summoner/v3/summoners/by-name/"+name, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	summoner := Summoner{}
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
	} else if summoner.ID == 0 {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = time.Now()

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
	} else if summoner.ID == 0 {
		return nil, fmt.Errorf("User does not exist")
	}

	summoner.Timestamp = time.Now()

	return &summoner, nil
}
