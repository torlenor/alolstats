package storage

import (
	"github.com/torlenor/alolstats/riotclient"
)

// GetRunesReforged returns a list of all currently known summoner spells
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetRunesReforged(forceUpdate bool) *riotclient.RunesReforgedList {
	// duration := time.Since(s.backend.GetRunesReforgedTimeStamp())
	// if (duration.Minutes() > float64(s.config.MaxAgeRunesReforged)) || forceUpdate {
	summonerSpells, err := s.riotClient.RunesReforged()
	if err != nil {
		return nil
	}
	// 	if err != nil {
	// 		s.log.Warnln(err)
	// 		summonerSpells, err := s.backend.GetRunesReforged()
	// 		if err != nil {
	// 			s.log.Warnln(err)
	// 			return nil
	// 		}
	// 		s.log.Debugf("Could not get Summoner Spells from Client, returning from Storage Backend instead")
	// 		return summonerSpells
	// 	}
	// 	err = s.backend.StoreRunesReforged(summonerSpells)
	// 	if err != nil {
	// 		s.log.Warnln("Could not store Summoner Spells in storage backend:", err)
	// 	}
	// 	return summonerSpells
	// }
	// summonerSpells, err := s.backend.GetRunesReforged()
	// if err != nil {
	// 	summonerSpells, errClient := s.riotClient.RunesReforged()
	// 	if errClient != nil {
	// 		s.log.Warnln(errClient)
	// 		return nil
	// 	}
	// 	s.log.Warnln("Could not get Summoner Spells from storage backend, returning from Client instead:", err)
	// 	err = s.backend.StoreRunesReforged(summonerSpells)
	// 	if err != nil {
	// 		s.log.Warnln("Could not store Summoner Spells in storage backend:", err)
	// 	}
	// 	return summonerSpells
	// }
	return summonerSpells
}
