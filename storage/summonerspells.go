package storage

import (
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// GetSummonerSpells returns a list of all currently known summoner spells
// gameVersion is the game version we want to have
// language is the langauge that we want to have
func (s *Storage) GetSummonerSpells(gameVersion, language string) (riotclient.SummonerSpellsList, error) {
	summonerSpellsList, err := s.backend.GetSummonerSpells(gameVersion, language)
	if err != nil || len(summonerSpellsList) == 0 {
		summonerSpellsList, errClient := s.riotClient.SummonerSpellsSpecificVersionLanguage(gameVersion, language)
		if errClient != nil {
			s.log.Warnln("Error from client when trying to get summoner spells:", errClient)
			return nil, fmt.Errorf("Could not get Summoner Spells from Backend or Client")
		}
		s.log.Warnln("Could not get Summoner Spells from storage backend, returning from Client instead:", err)
		err = s.backend.StoreSummonerSpells(gameVersion, language, *summonerSpellsList)
		if err != nil {
			s.log.Warnln("Could not store Summoner Spells in storage backend:", err)
		}
		return *summonerSpellsList, nil
	}

	return summonerSpellsList, nil
}

//StoreSummonerSpells stores a new list of Runes Reforged for the given game version and language in the backend
func (s *Storage) StoreSummonerSpells(gameVersion, language string, summonerSpellsList riotclient.SummonerSpellsList) error {
	return s.backend.StoreSummonerSpells(gameVersion, language, summonerSpellsList)
}
