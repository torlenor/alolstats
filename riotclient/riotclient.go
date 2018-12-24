// Package riotclient provides the Riot API Client interfaces
package riotclient

// ClientBase defines basic Riot Client operations
type ClientBase interface {
	Start()
	Stop()
}

// ClientChampion defines an interface to Champion API calls
type ClientChampion interface {
	Champions() (s *ChampionList, err error)
	ChampionRotations() (s *FreeRotation, err error)
}

// ClientSummoner defines an interface to Summoner API calls
type ClientSummoner interface {
	SummonerByName(name string) (s *SummonerDTO, err error)
	SummonerByAccountID(accountID string) (s *SummonerDTO, err error)
	SummonerBySummonerID(summonerID string) (s *SummonerDTO, err error)
}

// Client defines the interface for a Riot API Client
type Client interface {
	ClientBase
	ClientChampion
	ClientSummoner

	MatchByID(matchID uint64) (s *MatchDTO, err error)

	MatchesByAccountID(accountID string, startIndex uint32, endIndex uint32) (s *MatchList, err error)

	LeagueByQueue(league string, queue string) (*LeagueListDTO, error)
}
