package storage

import "github.com/torlenor/alolstats/riotclient"

// GetLeagueByQueue returns a league identified by its name for a specific queue name
func (s *Storage) GetLeagueByQueue(league string, queue string) (*riotclient.LeagueListDTO, error) {
	s.log.Debugf("GetLeagueByQueue(%s, %s)", league, queue)
	leagueData, errClient := s.riotClient.LeagueByQueue(league, queue)
	if errClient != nil {
		s.log.Warnln("Could not get data from Riot API:", errClient)
		return nil, errClient
	}
	s.log.Debugf("Returned %s for queue %s from Riot API", league, queue)
	return leagueData, nil
}
