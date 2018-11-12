package storage

import (
	"fmt"

	"github.com/torlenor/alolstats/riotclient"
)

type mockClient struct {
	failChampions      bool
	champions          riotclient.ChampionList
	championsRetrieved bool
}

func (c *mockClient) SummonerByName(name string) (s *riotclient.Summoner, err error) {
	return nil, nil
}

func (c *mockClient) reset() {
	c.failChampions = false
	c.champions = riotclient.ChampionList{}
	c.championsRetrieved = false
}

func (c *mockClient) setChampions(champions riotclient.ChampionList) {
	c.champions = champions
}

func (c *mockClient) setFailChampions(fail bool) {
	c.failChampions = fail
}

func (c *mockClient) getChampionsRetrieved() bool {
	return c.championsRetrieved
}

func (c *mockClient) Champions() (s *riotclient.ChampionList, err error) {
	c.championsRetrieved = true

	if c.failChampions {
		return &riotclient.ChampionList{}, fmt.Errorf("Error retreiving champions")
	}

	return &c.champions, nil
}

func (c *mockClient) FreeRotation() (*riotclient.FreeRotation, error) {
	return nil, nil
}
