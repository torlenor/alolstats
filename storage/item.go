package storage

import (
	"github.com/torlenor/alolstats/riotclient"
)

// GetItems returns a list of all currently known summoner spells
// forceUpdate will try to update the champion, if it is false the config settings will be considered if update is required
func (s *Storage) GetItems(forceUpdate bool) (*riotclient.ItemList, error) {
	// duration := time.Since(s.backend.GetItemsTimeStamp())
	// if (duration.Minutes() > float64(s.config.MaxAgeItems)) || forceUpdate {
	// 	items, err := s.riotClient.Items()
	// 	if err != nil {
	// 		s.log.Warnln(err)
	// 		items, err := s.backend.GetItems()
	// 		if err != nil {
	// 			s.log.Warnln(err)
	// 			return nil
	// 		}
	// 		s.log.Debugf("Could not get Items from Client, returning from Storage Backend instead")
	// 		return items
	// 	}
	// 	err = s.backend.StoreItems(items)
	// 	if err != nil {
	// 		s.log.Warnln("Could not store Items in storage backend:", err)
	// 	}
	// 	return items
	// }
	// items, err := s.backend.GetItems()
	// if err != nil {
	// 	items, errClient := s.riotClient.Items()
	// 	if errClient != nil {
	// 		s.log.Warnln(errClient)
	// 		return nil
	// 	}
	// 	s.log.Warnln("Could not get Items from storage backend, returning from Client instead:", err)
	// 	err = s.backend.StoreItems(items)
	// 	if err != nil {
	// 		s.log.Warnln("Could not store Items in storage backend:", err)
	// 	}
	// 	return items
	// }
	return s.riotClient.Items()
}

// GetItemsSpecificVersionLanguage returns a list of all items for a given gameVersion and language
func (s *Storage) GetItemsSpecificVersionLanguage(gameVersion, language string, forceUpdate bool) (*riotclient.ItemList, error) {
	// items, err := s.backend.GetItems(gameVersion, language)
	// if err != nil {
	items, errClient := s.riotClient.ItemsSpecificVersionLanguage(gameVersion, language)
	if errClient != nil {
		s.log.Warnln(errClient)
		return nil, errClient
	}
	// s.log.Warnln("Could not get Items from storage backend, returning from Client instead:", err)
	// err = s.backend.StoreItems(items)
	// if err != nil {
	// 	s.log.Warnln("Could not store Items in storage backend:", err)
	// }
	return items, nil
	// }
	// return items
}
