package storage

import (
	"fmt"

	"github.com/torlenor/alolstats/riotclient"
)

type mockClient struct {
	failChampions      bool
	champions          riotclient.ChampionList
	championsRetrieved bool

	failFreeRotation      bool
	freeRotation          riotclient.FreeRotation
	freeRotationRetrieved bool
}

func (c *mockClient) SummonerByName(name string) (s *riotclient.Summoner, err error) {
	return nil, nil
}

func (c *mockClient) SummonerByAccountID(id uint64) (s *riotclient.Summoner, err error) {
	return nil, nil
}

func (c *mockClient) SummonerBySummonerID(id uint64) (s *riotclient.Summoner, err error) {
	return nil, nil
}

func (c *mockClient) MatchByID(id uint64) (s *riotclient.Match, err error) {
	return nil, nil
}

func (c *mockClient) MatchesByAccountID(id uint64, startIndex uint32, endIndex uint32) (s *riotclient.MatchList, err error) {
	return nil, nil
}

func (c *mockClient) reset() {
	c.failChampions = false
	c.champions = riotclient.ChampionList{}
	c.championsRetrieved = false

	c.failFreeRotation = false
	c.freeRotation = riotclient.FreeRotation{}
	c.freeRotationRetrieved = false
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

func (c *mockClient) setFreeRotation(freeRotation riotclient.FreeRotation) {
	c.freeRotation = freeRotation
}

func (c *mockClient) setFailFreeRotation(fail bool) {
	c.failFreeRotation = fail
}

func (c *mockClient) getFreeRotationRetrieved() bool {
	return c.freeRotationRetrieved
}

func (c *mockClient) FreeRotation() (*riotclient.FreeRotation, error) {
	c.freeRotationRetrieved = true

	if c.failFreeRotation {
		return &riotclient.FreeRotation{}, fmt.Errorf("Error retreiving Free Rotation")
	}

	return &c.freeRotation, nil
}
