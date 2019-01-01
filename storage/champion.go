package storage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetChampions returns a list of all currently known champions
func (s *Storage) GetChampions() riotclient.ChampionList {
	duration := time.Since(s.backend.GetChampionsTimeStamp())
	if duration.Minutes() > float64(s.config.MaxAgeChampion) {
		champions, err := s.riotClient.Champions()
		if err != nil {
			s.log.Warnln(err)
			champions, err := s.backend.GetChampions()
			if err != nil {
				s.log.Warnln(err)
				return riotclient.ChampionList{}
			}
			s.log.Debugf("Could not get Champions from Client, returning from Storage Backend instead")
			return *champions
		}
		err = s.backend.StoreChampions(champions)
		if err != nil {
			s.log.Warnln("Could not store Champions in storage backend:", err)
		}
		s.log.Debugf("Returned Champions from Client")
		return *champions
	}
	champions, err := s.backend.GetChampions()
	if err != nil {
		champions, errClient := s.riotClient.Champions()
		if errClient != nil {
			s.log.Warnln(err)
			return riotclient.ChampionList{}
		}
		s.log.Warnln("Could not get Champions from storage backend, returning from Client instead:", err)
		err = s.backend.StoreChampions(champions)
		if err != nil {
			s.log.Warnln("Could not store Champions in storage backend:", err)
		}
		return *champions
	}
	s.log.Debugf("Returned Champions from Storage Backend")
	return *champions
}

// GetChampionByID returns a champion identified by its ID
func (s *Storage) GetChampionByID(id uint32) (riotclient.Champion, error) {
	idStr := strconv.FormatUint(uint64(id), 10)

	champions := s.GetChampions()

	champion := riotclient.Champion{}
	found := false
	for _, champion = range champions.Champions {
		if idStr == champion.ID {
			found = true
			break
		}
	}

	if found {
		return champion, nil
	}

	return champion, fmt.Errorf("Champion with ID %d not found", id)
}
