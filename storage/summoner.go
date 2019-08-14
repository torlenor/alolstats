package storage

import (
	"fmt"
	"time"

	"git.abyle.org/hps/alolstats/utils"

	"git.abyle.org/hps/alolstats/riotclient"
)

// Summoner is the storage type used for Summoner Data
type Summoner struct {
	SummonerDTO  riotclient.SummonerDTO
	SummonerName string
	SummonerID   string
	AccountID    string
	PUUID        string
}

func (s *Storage) storeSummoner(summoner riotclient.SummonerDTO) error {
	return s.backend.StoreSummoner(&Summoner{
		SummonerDTO:  summoner,
		SummonerName: utils.CleanUpSummonerName(summoner.Name),
		SummonerID:   summoner.ID,
		AccountID:    summoner.AccountID,
		PUUID:        summoner.PuuID,
	})
}

// getSummonerByNameFromClient returns a Summoner identified by name
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) getSummonerByNameFromClient(client riotclient.Client, name string, forceUpdate bool) (riotclient.SummonerDTO, error) {
	name = utils.CleanUpSummonerName(name)
	duration := time.Since(s.backend.GetSummonerByNameTimeStamp(name))
	if (duration.Minutes() > float64(s.config.MaxAgeSummoner)) || forceUpdate {
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
		err = s.storeSummoner(*summoner)
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
		err = s.storeSummoner(*summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner %s from Riot API", name)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner %s from Storage", name)
	return summoner.SummonerDTO, nil
}

// GetSummonerByName returns a Summoner identified by name
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetSummonerByName(name string, forceUpdate bool) (riotclient.SummonerDTO, error) {
	return s.getSummonerByNameFromClient(s.riotClient, name, forceUpdate)
}

// GetRegionalSummonerByName returns a Summoner identified by name for a specific region
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetRegionalSummonerByName(region string, name string, forceUpdate bool) (riotclient.SummonerDTO, error) {
	if client, ok := s.riotClients[region]; ok {
		return s.getSummonerByNameFromClient(client, name, forceUpdate)
	}
	return riotclient.SummonerDTO{}, fmt.Errorf("Invalid region specified: %s", region)
}

// getSummonerBySummonerIDFromClient returns a Summoner identified by Summoner ID
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) getSummonerBySummonerIDFromClient(client riotclient.Client, summonerID string, forceUpdate bool) (riotclient.SummonerDTO, error) {
	if len(summonerID) == 0 {
		return riotclient.SummonerDTO{}, fmt.Errorf("Summoner ID cannot be empty")
	}
	duration := time.Since(s.backend.GetSummonerBySummonerIDTimeStamp(summonerID))
	if (duration.Minutes() > float64(s.config.MaxAgeSummoner)) || forceUpdate {
		summoner, err := client.SummonerBySummonerID(summonerID)
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
		err = s.storeSummoner(*summoner)
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
		err = s.storeSummoner(*summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with SummonerID %s from Riot API", summonerID)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner with SummonerID %s from Storage", summonerID)
	return summoner.SummonerDTO, nil
}

// GetSummonerBySummonerID returns a Summoner identified by Summoner ID
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetSummonerBySummonerID(summonerID string, forceUpdate bool) (riotclient.SummonerDTO, error) {
	return s.getSummonerBySummonerIDFromClient(s.riotClient, summonerID, forceUpdate)
}

// GetRegionalSummonerBySummonerID returns a Summoner identified by Summoner ID for a specific region
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetRegionalSummonerBySummonerID(region string, summonerID string, forceUpdate bool) (riotclient.SummonerDTO, error) {
	if client, ok := s.riotClients[region]; ok {
		return s.getSummonerBySummonerIDFromClient(client, summonerID, forceUpdate)
	}
	return riotclient.SummonerDTO{}, fmt.Errorf("Invalid region specified: %s", region)
}

// GetSummonerByAccountID returns a Summoner identified by Account ID
func (s *Storage) GetSummonerByAccountID(accountID string, forceUpdate bool) (riotclient.SummonerDTO, error) {
	if len(accountID) == 0 {
		return riotclient.SummonerDTO{}, fmt.Errorf("Account ID cannot be empty")
	}
	duration := time.Since(s.backend.GetSummonerByAccountIDTimeStamp(accountID))
	if (duration.Minutes() > float64(s.config.MaxAgeSummoner)) || forceUpdate {
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
		err = s.storeSummoner(*summoner)
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
		err = s.storeSummoner(*summoner)
		if err != nil {
			s.log.Warnln("Could not store Summoner in storage backend:", err)
		}
		s.log.Debugf("Returned Summoner with AccountID %s from Riot API", accountID)
		return *summoner, nil
	}
	s.log.Debugf("Returned Summoner with AccountID %s from Storage", accountID)
	return summoner.SummonerDTO, nil
}
