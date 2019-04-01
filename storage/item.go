package storage

import "github.com/torlenor/alolstats/riotclient"

// GetItems returns a list of all items for a given gameVersion and language
func (s *Storage) GetItems(gameVersion, language string) (*riotclient.ItemList, error) {
	// items, err := s.backend.GetItems(gameVersion, language)
	// if err != nil {
	items, errClient := s.riotClient.Items(gameVersion, language)
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
