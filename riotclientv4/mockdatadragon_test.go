package riotclientv4

// MockRiotClientDD Mock of Riot LoL API DataDragon client
type MockRiotClientDD struct {
	championsJSON      []byte
	itemJSON           []byte
	summonerSpellsJSON []byte
	versionsJSON       []byte
	err                error
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

// GetDataDragonItems returns the items for the live game version
func (c *MockRiotClientDD) GetDataDragonItems() ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.itemJSON, nil
}

func (c *MockRiotClientDD) GetLoLVersions() ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.versionsJSON, nil
}

func (c *MockRiotClientDD) GetDataDragonSummonerSpells() ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.summonerSpellsJSON, nil
}

func (c *MockRiotClientDD) GetDataDragonSummonerSpellsSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.summonerSpellsJSON, nil
}

// GetDataDragonRunesReforgedSpecificVersionLanguage returns the Runes Reforged for a specific gameVersion and langauge
func (c *MockRiotClientDD) GetDataDragonRunesReforgedSpecificVersionLanguage(gameVersion, language string) ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.itemJSON, nil
}

// GetDataDragonRunesReforged returns the Runes Reforged for the live game version
func (c *MockRiotClientDD) GetDataDragonRunesReforged() ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.itemJSON, nil
}
