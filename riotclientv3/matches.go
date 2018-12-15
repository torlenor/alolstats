package riotclientv3

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/torlenor/alolstats/riotclient"
)

// MatchByID gets a match by its ID
func (c *RiotClientV3) MatchByID(id uint64) (s *riotclient.Match, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v3/matches/3827449823
	idStr := strconv.FormatUint(id, 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/match/v3/matches/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	match := riotclient.Match{}
	err = json.Unmarshal(data, &match)
	if err != nil {
		return nil, err
	} else if match.GameID == 0 {
		return nil, fmt.Errorf("Match GameID invalid, probably empty data")
	}

	return &match, nil
}

// MatchesByAccountID gets a match by AccountID. Provide a start and end index to fetch matches.
// The matches will be fetched as [beginIndex, endIndex) and the indices start at 0.
//
// For example to fetch the last match of an account enter startIndex=0, endIndex=1
// and to fetch the first match of an account enter startIndex=(totalNumMatches-1), endIndex=totalNumMatches.
func (c *RiotClientV3) MatchesByAccountID(id uint64, startIndex uint32, endIndex uint32) (s *riotclient.MatchList, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v3/matchlists/by-account/40722898?beginIndex=0&endIndex=100
	idStr := strconv.FormatUint(id, 10)
	startIndexStr := strconv.FormatUint(uint64(startIndex), 10)
	endIndexStr := strconv.FormatUint(uint64(endIndex), 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/match/v3/matchlists/by-account/"+idStr+"?beginIndex="+startIndexStr+"&endIndex="+endIndexStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	matchList := riotclient.MatchList{}
	err = json.Unmarshal(data, &matchList)
	if err != nil {
		return nil, err
	}

	return &matchList, nil
}
