package storage

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetChampions returns a list of all currently known champions
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetChampions(forceUpdate bool) riotclient.ChampionsList {
	duration := time.Since(s.backend.GetChampionsTimeStamp())
	if (duration.Minutes() > float64(s.config.MaxAgeChampion)) || forceUpdate {
		champions, err := s.riotClient.Champions()
		if err != nil {
			s.log.Warnln(err)
			champions, err := s.backend.GetChampions()
			if err != nil {
				s.log.Warnln(err)
				return nil
			}
			s.log.Debugf("Could not get Champions from Client, returning from Storage Backend instead")
			return champions
		}
		err = s.backend.StoreChampions(champions)
		if err != nil {
			s.log.Warnln("Could not store Champions in storage backend:", err)
		}
		return champions
	}
	champions, err := s.backend.GetChampions()
	if err != nil {
		champions, errClient := s.riotClient.Champions()
		if errClient != nil {
			s.log.Warnln(err)
			return nil
		}
		s.log.Warnln("Could not get Champions from storage backend, returning from Client instead:", err)
		err = s.backend.StoreChampions(champions)
		if err != nil {
			s.log.Warnln("Could not store Champions in storage backend:", err)
		}
		return champions
	}
	return champions
}

// GetChampionByID returns a champion identified by its ID
// ID is a simplified name or internal Riot name for the champion
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetChampionByID(ID string, forceUpdate bool) (riotclient.Champion, error) {
	champions := s.GetChampions(forceUpdate)

	champion := riotclient.Champion{}
	found := false
	for _, champion = range champions {
		if ID == champion.ID {
			found = true
			break
		}
	}

	if found {
		return champion, nil
	}

	return champion, fmt.Errorf("Champion with ID %s not found", ID)
}

// GetChampionByName returns a champion identified by its Name
// name is the official name for the champion and it must match exactly, e.g. Cho'Gath
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetChampionByName(name string, forceUpdate bool) (riotclient.Champion, error) {
	champions := s.GetChampions(forceUpdate)

	champion := riotclient.Champion{}
	found := false
	for _, champion = range champions {
		if name == champion.Name {
			found = true
			break
		}
	}

	if found {
		return champion, nil
	}

	return champion, fmt.Errorf("Champion with Name %s not found", name)
}

// GetChampionByKey returns a champion identified by its key
// key is a, for now, numeric identifier from Riot for a champion
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetChampionByKey(key string, forceUpdate bool) (riotclient.Champion, error) {
	champions := s.GetChampions(forceUpdate)

	champion := riotclient.Champion{}
	found := false
	for _, champion = range champions {
		if key == champion.Key {
			found = true
			break
		}
	}

	if found {
		return champion, nil
	}

	return champion, fmt.Errorf("Champion with Key %s not found", key)
}
