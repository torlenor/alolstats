package riotclientv4

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"git.abyle.org/hps/alolstats/riotclient"
)

// MatchByID gets a match by its ID
func (c *RiotClientV4) MatchByID(id uint64) (s *riotclient.MatchDTO, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v4/matches/3827449823
	idStr := strconv.FormatUint(id, 10)
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/match/"+c.config.APIVersion+"/matches/"+idStr, "GET", "")
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

// MatchesByAccountID gets a match by AccountID
// args: List of arguments to the query. They are directly passed to the request.
// Refer to https://developer.riotgames.com/api-methods/#match-v4/GET_getMatchlist for details.
func (c *RiotClientV4) MatchesByAccountID(accountID string, args map[string]string) (s *riotclient.MatchlistDTO, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v4/matchlists/by-account/1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g?endIndex=100&beginIndex=0
	basicAPICall := "https://" + c.config.Region + ".api.riotgames.com/lol/match/" + c.config.APIVersion + "/matchlists/by-account/" + accountID
	fullAPICall := basicAPICall
	if len(args) > 0 {
		fullAPICall = fullAPICall + "?"

		keys := make([]string, 0, len(args))
		for k := range args {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fullAPICall = fullAPICall + k + "=" + args[k] + "&"
		}
		if last := len(fullAPICall) - 1; last >= 0 && fullAPICall[last] == '&' {
			fullAPICall = fullAPICall[:last]
		}
	}
	data, err := apiCall(c, fullAPICall, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	matchList := riotclient.MatchlistDTO{}
	err = json.Unmarshal(data, &matchList)
	if err != nil {
		return nil, err
	}

	return &matchList, nil
}

// MatchTimeLineByID gets the Match TimeLine for a certain match identified by its MatchID/GameID
func (c *RiotClientV4) MatchTimeLineByID(matchID uint64) (t *riotclient.MatchTimelineDTO, err error) {
	// /lol/match/v4/timelines/by-match/{matchId}
	idStr := strconv.FormatUint(matchID, 10)
	data, err := apiCall(c, "https://"+c.config.Region+".api.riotgames.com/lol/match/"+c.config.APIVersion+"/timelines/by-match/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	matchTimeLine := riotclient.MatchTimelineDTO{}
	err = json.Unmarshal(data, &matchTimeLine)
	if err != nil {
		return nil, fmt.Errorf("MatchTimeLineByID error unmarshaling: %s", err)
	}

	return &matchTimeLine, nil
}
