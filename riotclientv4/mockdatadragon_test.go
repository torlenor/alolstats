package riotclientv4

// MockRiotClientDD Mock of Riot LoL API DataDragon client
type MockRiotClientDD struct {
	championsJSON []byte
	itemJSON      []byte
	err           error
}

// GetDataDragonChampions returns the current champions available for the live game version
func (c *MockRiotClientDD) GetDataDragonChampions() ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.championsJSON, nil
}

// GetDataDragonItemsSpecificVersionLanguage returns the items for a specific gameVersion and langauge
func (c *MockRiotClientDD) GetDataDragonItemsSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.itemJSON, nil
}

func (c *MockRiotClientDD) GetLoLVersions() ([]byte, error) {
	return nil, nil
}

func (c *MockRiotClientDD) GetDataDragonSummonerSpells() ([]byte, error) {
	return nil, nil
}

func (c *MockRiotClientDD) GetDataDragonSummonerSpellsSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	return nil, nil
}
