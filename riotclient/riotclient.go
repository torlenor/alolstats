// Package riotclient provides the Riot API client
package riotclient

// Client defines the interface for a Riot API client
type Client interface {
	Start()
	Stop()

	SummonerByName(name string) (s *Summoner, err error)
	SummonerByAccountID(id uint64) (s *Summoner, err error)
	SummonerBySummonerID(id uint64) (s *Summoner, err error)

	Champions() (s *ChampionList, err error)

	FreeRotation() (*FreeRotation, error)

	MatchByID(id uint64) (s *Match, err error)

	MatchesByAccountID(id uint64, startIndex uint32, endIndex uint32) (s *MatchList, err error)

	ChallengerLeagueByQueue(queue string) (*LeagueData, error)
	MasterLeagueByQueue(queue string) (*LeagueData, error)
}
