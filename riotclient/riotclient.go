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
	SummonerByName(name string) (s *Summoner, err error)
	SummonerByAccountID(id uint64) (s *Summoner, err error)
	SummonerBySummonerID(id uint64) (s *Summoner, err error)
}

// Client defines the interface for a Riot API Client
type Client interface {
	ClientBase
	ClientChampion
	ClientSummoner

	MatchByID(id uint64) (s *Match, err error)

	MatchesByAccountID(id uint64, startIndex uint32, endIndex uint32) (s *MatchList, err error)

	ChallengerLeagueByQueue(queue string) (*LeagueData, error)
	MasterLeagueByQueue(queue string) (*LeagueData, error)
}
