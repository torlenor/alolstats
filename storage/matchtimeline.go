package storage

import (
	"fmt"
	"strings"

	"git.abyle.org/hps/alolstats/riotclient"
)

// fetchAndStoreTimeLineFromClient gets a timeline from Riot Client and stores it in storage backend if it doesn't exist, yet
func (s *Storage) fetchAndStoreMatchTimeLineFromClient(client riotclient.Client, match *riotclient.MatchDTO) (*riotclient.MatchTimelineDTO, error) {
	_, err := s.backend.GetMatchTimeLine(uint64(match.GameID))
	if err != nil {
		timeline, err := client.MatchTimeLineByID(uint64(match.GameID))
		if err != nil {
			s.log.Warnln(err)
			return nil, err
		}
		s.log.Debugf("Storing Match TimeLine %d from Riot API in Backend", uint64(match.GameID))
		err = s.backend.StoreMatchTimeLine(match, timeline)
		if err != nil {
			s.log.Errorf("Error storing Match TimeLine %s from Riot API in Backend: ", uint64(match.GameID), err)
		}
		return timeline, nil
	}

	s.log.Debugf("Match TimeLine already found in storage, did not fetch it")

	return nil, nil
}

// RegionalFetchAndStoreMatchTimeLine gets a timeline from Riot Client for a given match and stores it in storage backend if it doesn't exist, yet
func (s *Storage) RegionalFetchAndStoreMatchTimeLine(match *riotclient.MatchDTO) (*riotclient.MatchTimelineDTO, error) {
	if client, ok := s.riotClients[strings.ToLower(match.PlatformID)]; ok {
		return s.fetchAndStoreMatchTimeLineFromClient(client, match)
	}
	return nil, fmt.Errorf("Invalid region specified: %s", match.PlatformID)
}
