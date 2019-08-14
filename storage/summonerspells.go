package storage

import (
	"time"

	"git.abyle.org/hps/alolstats/riotclient"
)

// GetSummonerSpells returns a list of all currently known summoner spells
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetSummonerSpells(forceUpdate bool) *riotclient.SummonerSpellsList {
	duration := time.Since(s.backend.GetSummonerSpellsTimeStamp())
	if (duration.Minutes() > float64(s.config.MaxAgeSummonerSpells)) || forceUpdate {
		summonerSpells, err := s.riotClient.SummonerSpells()
		if err != nil {
			s.log.Warnln(err)
			summonerSpells, err := s.backend.GetSummonerSpells()
			if err != nil {
				s.log.Warnln(err)
				return nil
			}
			s.log.Debugf("Could not get Summoner Spells from Client, returning from Storage Backend instead")
			return summonerSpells
		}
		err = s.backend.StoreSummonerSpells(summonerSpells)
		if err != nil {
			s.log.Warnln("Could not store Summoner Spells in storage backend:", err)
		}
		return summonerSpells
	}
	summonerSpells, err := s.backend.GetSummonerSpells()
	if err != nil {
		summonerSpells, errClient := s.riotClient.SummonerSpells()
		if errClient != nil {
			s.log.Warnln(errClient)
			return nil
		}
		s.log.Warnln("Could not get Summoner Spells from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummonerSpells(summonerSpells)
		if err != nil {
			s.log.Warnln("Could not store Summoner Spells in storage backend:", err)
		}
		return summonerSpells
	}
	return summonerSpells
}
