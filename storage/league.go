package storage

import "github.com/torlenor/alolstats/riotclient"

// GetChallengerLeagueByQueue returns a ChallengerLeague identified by its queue name
func (s *Storage) GetChallengerLeagueByQueue(queue string) (riotclient.LeagueData, error) {
	challengerLeague, errClient := s.riotClient.ChallengerLeagueByQueue(queue)
	if errClient != nil {
		s.log.Warnln("Could not get data from either Storage nor Riot API:", errClient)
		return riotclient.LeagueData{}, errClient
	}
	s.log.Debugf("Returned ChallengerLeague for queue %s from Riot API", queue)
	return *challengerLeague, nil
}

// GetMasterLeagueByQueue returns a MasterLeague identified by its queue name
func (s *Storage) GetMasterLeagueByQueue(queue string) (riotclient.LeagueData, error) {
	masterLeague, errClient := s.riotClient.MasterLeagueByQueue(queue)
	if errClient != nil {
		s.log.Warnln("Could not get data from either Storage nor Riot API:", errClient)
		return riotclient.LeagueData{}, errClient
	}
	s.log.Debugf("Returned MasterLeague for queue %s from Riot API", queue)
	return *masterLeague, nil
}
