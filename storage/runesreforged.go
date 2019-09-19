package storage

import (
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// GetRunesReforged returns a list of all currently known summoner spells
// gameVersion is the game version we want to have
// language is the langauge that we want to have
func (s *Storage) GetRunesReforged(gameVersion, language string) (riotclient.RunesReforgedList, error) {
	runesReforgedList, err := s.backend.GetRunesReforged(gameVersion, language)
	if err != nil {
		runesReforgedList, errClient := s.riotClient.RunesReforgedSpecificVersionLanguage(gameVersion, language)
		if errClient != nil {
			s.log.Warnln(errClient)
			return nil, fmt.Errorf("Could not get Runes Reforged from Backend or Client: %s", err)
		}
		s.log.Warnln("Could not get Summoner Spells from storage backend, returning from Client instead:", err)
		err = s.backend.StoreRunesReforged(gameVersion, language, *runesReforgedList)
		if err != nil {
			s.log.Warnln("Could not store Summoner Spells in storage backend:", err)
		}
		return *runesReforgedList, nil
	}

	return runesReforgedList, nil
}

//StoreRunesReforged stores a new list of Runes Reforged for the given game version and language in the backend
func (s *Storage) StoreRunesReforged(gameVersion, language string, runesReforgedList riotclient.RunesReforgedList) error {
	return s.backend.StoreRunesReforged(gameVersion, language, runesReforgedList)
}
