package storage

import (
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// GetItems returns a list of all currently known summoner spells
// gameVersion is the game version we want to have
// language is the langauge that we want to have
func (s *Storage) GetItems(gameVersion, language string) (riotclient.ItemList, error) {
	itemsList, err := s.backend.GetItems(gameVersion, language)
	if err != nil {
		itemsList, errClient := s.riotClient.ItemsSpecificVersionLanguage(gameVersion, language)
		if errClient != nil {
			s.log.Warnln(errClient)
			return nil, fmt.Errorf("Could not get Items from Backend or Client")
		}
		s.log.Warnln("Could not get Summoner Spells from storage backend, returning from Client instead:", err)
		err = s.backend.StoreItems(gameVersion, language, *itemsList)
		if err != nil {
			s.log.Warnln("Could not store Summoner Spells in storage backend:", err)
		}
		return *itemsList, nil
	}

	return itemsList, nil
}

//StoreItems stores a new list of Runes Reforged for the given game version and language in the backend
func (s *Storage) StoreItems(gameVersion, language string, itemsList riotclient.ItemList) error {
	return s.backend.StoreItems(gameVersion, language, itemsList)
}
