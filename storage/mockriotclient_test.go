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

	failSummoner         bool
	summoner             riotclient.SummonerDTO
	wasSummonerRetrieved bool
}

func (c *mockClient) Start() {
	// nothing
}

func (c *mockClient) Stop() {
	// nothing
}

func (c *mockClient) setSummoner(summoner riotclient.SummonerDTO) {
	c.summoner = summoner
}

func (c *mockClient) setFailSummoner(fail bool) {
	c.failSummoner = fail
}

func (c *mockClient) getWasSummonerRetrieved() bool {
	return c.wasSummonerRetrieved
}

func (c *mockClient) SummonerByName(name string) (s *riotclient.SummonerDTO, err error) {
	c.wasSummonerRetrieved = true

	if c.failSummoner {
		return &riotclient.SummonerDTO{}, fmt.Errorf("Error retreiving summoner")
	}

	return &c.summoner, nil
}

func (c *mockClient) SummonerByAccountID(accountID string) (s *riotclient.SummonerDTO, err error) {
	c.wasSummonerRetrieved = true

	if c.failSummoner {
		return &riotclient.SummonerDTO{}, fmt.Errorf("Error retreiving summoner")
	}

	return &c.summoner, nil
}

func (c *mockClient) SummonerBySummonerID(summonerID string) (s *riotclient.SummonerDTO, err error) {
	c.wasSummonerRetrieved = true

	if c.failSummoner {
		return &riotclient.SummonerDTO{}, fmt.Errorf("Error retreiving summoner")
	}

	return &c.summoner, nil
}

func (c *mockClient) MatchByID(id uint64) (s *riotclient.MatchDTO, err error) {
	return nil, nil
}

func (c *mockClient) MatchesByAccountID(accountID string, startIndex uint32, endIndex uint32) (s *riotclient.MatchList, err error) {
	return nil, nil
}

func (c *mockClient) LeagueByQueue(league string, queue string) (*riotclient.LeagueListDTO, error) {
	return nil, nil
}

func (c *mockClient) reset() {
	c.failChampions = false
	c.champions = riotclient.ChampionList{}
	c.championsRetrieved = false

	c.failFreeRotation = false
	c.freeRotation = riotclient.FreeRotation{}
	c.freeRotationRetrieved = false

	c.failSummoner = false
	c.summoner = riotclient.SummonerDTO{}
	c.wasSummonerRetrieved = false
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

func (c *mockClient) ChampionRotations() (*riotclient.FreeRotation, error) {
	c.freeRotationRetrieved = true

	if c.failFreeRotation {
		return &riotclient.FreeRotation{}, fmt.Errorf("Error retreiving Free Rotation")
	}

	return &c.freeRotation, nil
}

func (c *mockClient) ActiveGameBySummonerID(summonerID string) (*riotclient.CurrentGameInfoDTO, error) {
	return nil, fmt.Errorf("Not implemented")
}
func (c *mockClient) FeaturedGames() (*riotclient.FeaturedGamesDTO, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (c *mockClient) LeaguesForSummoner(encSummonerID string) (*riotclient.LeaguePositionDTOList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (c *mockClient) MatchTimeLineByID(matchID uint64) (t *riotclient.MatchTimelineDTO, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func (c *mockClient) SummonerByPUUID(PUUID string) (s *riotclient.SummonerDTO, err error) {
	return nil, fmt.Errorf("Not implemented")
}
