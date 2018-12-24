package riotclientv4

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/torlenor/alolstats/riotclient"
)

// MatchByID gets a match by its ID
func (c *RiotClientV4) MatchByID(id uint64) (s *riotclient.MatchDTO, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v4/matches/3827449823
	idStr := strconv.FormatUint(id, 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/match/"+c.config.APIVersion+"/matches/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	match := riotclient.MatchDTO{}
	err = json.Unmarshal(data, &match)
	if err != nil {
		return nil, fmt.Errorf("MatchById error unmarshaling: %s", err)
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
func (c *RiotClientV4) MatchesByAccountID(acountID string, startIndex uint32, endIndex uint32) (s *riotclient.MatchList, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v4/matchlists/by-account/1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g?endIndex=100&beginIndex=0
	startIndexStr := strconv.FormatUint(uint64(startIndex), 10)
	endIndexStr := strconv.FormatUint(uint64(endIndex), 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/match/"+c.config.APIVersion+"/matchlists/by-account/"+acountID+"?beginIndex="+startIndexStr+"&endIndex="+endIndexStr, "GET", "")
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
