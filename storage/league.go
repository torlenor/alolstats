package storage

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
	"github.com/torlenor/alolstats/utils"
)

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

// SummonerLeagues is the storage type used for Summoner Leagues Data
type SummonerLeagues struct {
	LeaguePositionDTOList riotclient.LeaguePositionDTOList
	SummonerName          string
	SummonerID            string
}

func (s *Storage) storeLeaguesForSummoner(leagues *riotclient.LeaguePositionDTOList) error {
	return s.backend.StoreLeaguesForSummoner(&SummonerLeagues{
		LeaguePositionDTOList: *leagues,
		SummonerName:          utils.CleanUpSummonerName(leagues.LeaguePosition[0].SummonerName),
		SummonerID:            leagues.LeaguePosition[0].SummonerID,
	})
}

// GetLeaguesForSummonerBySummonerID returns all Leagues a Summoner is placed in, identified by Summoner ID
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetLeaguesForSummonerBySummonerID(summonerID string, forceUpdate bool) (riotclient.LeaguePositionDTOList, error) {
	if len(summonerID) == 0 {
		return riotclient.LeaguePositionDTOList{}, fmt.Errorf("Summoner ID cannot be empty")
	}
	duration := time.Since(s.backend.GetSummonerBySummonerIDTimeStamp(summonerID))
	if (duration.Minutes() > float64(s.config.MaxAgeSummoner)) || forceUpdate {
		leagues, err := s.riotClient.LeaguesForSummoner(summonerID)
		if err != nil {
			s.log.Warnln("Could not get new data from Client, trying to get it from Storage instead", err)
			leagues, err := s.backend.GetLeaguesForSummonerBySummonerID(summonerID)
			if err != nil {
				s.log.Warnln("Could not get data from either Storage nor Client:", err)
				return riotclient.LeaguePositionDTOList{}, err
			}
			s.log.Debugf("Returned Leagues for Summoner with SummonerID %s from Storage", summonerID)
			return leagues.LeaguePositionDTOList, nil
		}
		err = s.storeLeaguesForSummoner(leagues)
		if err != nil {
			s.log.Warnln("Could not store Leagues for Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Leagues for Summoner with SummonerID %s from Riot API", summonerID)
		return *leagues, nil
	}
	leagues, err := s.backend.GetLeaguesForSummonerBySummonerID(summonerID)
	if err != nil {
		leagues, errClient := s.riotClient.LeaguesForSummoner(summonerID)
		if errClient != nil {
			s.log.Warnln("Could not get data from either Storage nor Client:", errClient)
			return riotclient.LeaguePositionDTOList{}, errClient
		}
		s.log.Warnln("Could not get Leagues for Summoner from storage backend, returning from Client instead:", err)
		err = s.storeLeaguesForSummoner(leagues)
		if err != nil {
			s.log.Warnln("Could not store Leagues for Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with SummonerID %s from Riot API", summonerID)
		return *leagues, nil
	}
	s.log.Debugf("Returned Summoner with SummonerID %s from Storage", summonerID)
	return leagues.LeaguePositionDTOList, nil
}
