package riotclientv4

// MockRiotClientDD Mock of Riot LoL API DataDragon client
type MockRiotClientDD struct {
	championsJSON []byte
	err           error
}

// GetDataDragonChampions returns the current champions available for the live game version
func (c *MockRiotClientDD) GetDataDragonChampions() ([]byte, error) {
	if c.err != nil {
		return []byte(""), c.err
	}

	return c.championsJSON, nil
}
