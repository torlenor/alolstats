package storage

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/utils"

	"github.com/torlenor/alolstats/riotclient"
)

// Summoner is the storage type used for Summoner Data
type Summoner struct {
	SummonerDTO  riotclient.SummonerDTO
	SummonerName string
}

// GetSummonerByName returns a Summoner identified by name
func (s *Storage) GetSummonerByName(name string) (riotclient.SummonerDTO, error) {
	name = utils.CleanUpSummonerName(name)
	duration := time.Since(s.backend.GetSummonerByNameTimeStamp(name))
	if duration.Minutes() > float64(s.config.MaxAgeSummoner) {
		summoner, err := s.riotClient.SummonerByName(name)
		if err != nil {
			s.log.Warnln("Could not get new data from Client, trying to get it from Storage instead", err)
			summoner, err := s.backend.GetSummonerByName(name)
			if err != nil {
				s.log.Warnln("Could not get data from either Storage nor Client:", err)
				return riotclient.SummonerDTO{}, err
			}
			s.log.Debugf("Returned Summoner %s from Storage", name)
			return summoner.SummonerDTO, nil
		}
		err = s.backend.StoreSummoner(&Summoner{SummonerDTO: *summoner, SummonerName: name})
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner %s from Riot API", name)
		return *summoner, nil
	}
	summoner, err := s.backend.GetSummonerByName(name)
	if err != nil {
		summoner, errClient := s.riotClient.SummonerByName(name)
		if errClient != nil {
			s.log.Warnln("Could not get data from either Storage nor Client:", errClient)
			return riotclient.SummonerDTO{}, errClient
		}
		s.log.Warnln("Could not get Summoner from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummoner(&Summoner{SummonerDTO: *summoner, SummonerName: name})
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner %s from Riot API", name)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner %s from Storage", name)
	return summoner.SummonerDTO, nil
}

// GetSummonerBySummonerID returns a Summoner identified by Summoner ID
func (s *Storage) GetSummonerBySummonerID(summonerID string) (riotclient.SummonerDTO, error) {
	if len(summonerID) == 0 {
		return riotclient.SummonerDTO{}, fmt.Errorf("Summoner ID cannot be empty")
	}
	duration := time.Since(s.backend.GetSummonerBySummonerIDTimeStamp(summonerID))
	if duration.Minutes() > float64(s.config.MaxAgeSummoner) {
		summoner, err := s.riotClient.SummonerBySummonerID(summonerID)
		if err != nil {
			s.log.Warnln("Could not get new data from Client, trying to get it from Storage instead", err)
			summoner, err := s.backend.GetSummonerBySummonerID(summonerID)
			if err != nil {
				s.log.Warnln("Could not get data from either Storage nor Client:", err)
				return riotclient.SummonerDTO{}, err
			}
			s.log.Debugf("Returned Summoner with SummonerID %s from Storage", summonerID)
			return summoner.SummonerDTO, nil
		}
		err = s.backend.StoreSummoner(&Summoner{SummonerDTO: *summoner, SummonerName: utils.CleanUpSummonerName(summoner.Name)})
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with SummonerID %s from Riot API", summonerID)
		return *summoner, nil
	}
	summoner, err := s.backend.GetSummonerBySummonerID(summonerID)
	if err != nil {
		summoner, errClient := s.riotClient.SummonerBySummonerID(summonerID)
		if errClient != nil {
			s.log.Warnln("Could not get data from either Storage nor Client:", errClient)
			return riotclient.SummonerDTO{}, errClient
		}
		s.log.Warnln("Could not get Summoner from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummoner(&Summoner{SummonerDTO: *summoner, SummonerName: utils.CleanUpSummonerName(summoner.Name)})
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with SummonerID %s from Riot API", summonerID)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner with SummonerID %s from Storage", summonerID)
	return summoner.SummonerDTO, nil
}

// GetSummonerByAccountID returns a Summoner identified by Account ID
func (s *Storage) GetSummonerByAccountID(accountID string) (riotclient.SummonerDTO, error) {
	if len(accountID) == 0 {
		return riotclient.SummonerDTO{}, fmt.Errorf("Account ID cannot be empty")
	}
	duration := time.Since(s.backend.GetSummonerByAccountIDTimeStamp(accountID))
	if duration.Minutes() > float64(s.config.MaxAgeSummoner) {
		summoner, err := s.riotClient.SummonerByAccountID(accountID)
		if err != nil {
			s.log.Warnln("Could not get new data from Client, trying to get it from Storage instead", err)
			summoner, err := s.backend.GetSummonerByAccountID(accountID)
			if err != nil {
				s.log.Warnln("Could not get data from either Storage nor Client:", err)
				return riotclient.SummonerDTO{}, err
			}
			s.log.Debugf("Returned Summoner with AccountID %s from Storage", accountID)
			return summoner.SummonerDTO, nil
		}
		err = s.backend.StoreSummoner(&Summoner{SummonerDTO: *summoner, SummonerName: utils.CleanUpSummonerName(summoner.Name)})
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with AccountID %s from Riot API", accountID)
		return *summoner, nil
	}
	summoner, err := s.backend.GetSummonerByAccountID(accountID)
	if err != nil {
		summoner, errClient := s.riotClient.SummonerByAccountID(accountID)
		if errClient != nil {
			s.log.Warnln("Could not get data from either Storage nor Client:", errClient)
			return riotclient.SummonerDTO{}, errClient
		}
		s.log.Warnln("Could not get Summoner from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummoner(&Summoner{SummonerDTO: *summoner, SummonerName: utils.CleanUpSummonerName(summoner.Name)})
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with AccountID %s from Riot API", accountID)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner with AccountID %s from Storage", accountID)
	return summoner.SummonerDTO, nil
}
