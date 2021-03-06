// Package riotclient provides the Riot API Client interfaces
package riotclient

// ClientBase defines basic Riot Client operations
type ClientBase interface {
	Start()
	Stop()
}

// ClientChampion defines an interface to Champion API calls
type ClientChampion interface {
	Champions() (s ChampionsList, err error)
	ChampionRotations() (s *FreeRotation, err error)
}

// ClientItem defines an interface to Item API calls
type ClientItem interface {
	Items() (*ItemList, error)
	ItemsSpecificVersionLanguage(gameVersion, language string) (*ItemList, error)
}

// ClientRunesReforged defines an interface to Item API calls
type ClientRunesReforged interface {
	RunesReforged() (*RunesReforgedList, error)
	RunesReforgedSpecificVersionLanguage(gameVersion, language string) (*RunesReforgedList, error)
}

// ClientSummonerSpells defines an interface to Summoner Spells API calls
type ClientSummonerSpells interface {
	SummonerSpells() (s *SummonerSpellsList, err error)
	SummonerSpellsSpecificVersionLanguage(gameVersion, language string) (*SummonerSpellsList, error)
}

// ClientSummoner defines an interface to Summoner API calls
type ClientSummoner interface {
	SummonerByName(name string) (s *SummonerDTO, err error)
	SummonerByAccountID(accountID string) (s *SummonerDTO, err error)
	SummonerBySummonerID(summonerID string) (s *SummonerDTO, err error)
	SummonerByPUUID(PUUID string) (s *SummonerDTO, err error)
}

// ClientMatch defines an interface to Match API calls
type ClientMatch interface {
	MatchByID(matchID uint64) (s *MatchDTO, err error)
	MatchesByAccountID(accountID string, args map[string]string) (s *MatchlistDTO, err error)
	MatchTimeLineByID(matchID uint64) (t *MatchTimelineDTO, err error)
}

// ClientLeague defines an interface to League API calls
type ClientLeague interface {
	LeagueByQueue(league string, queue string) (*LeagueListDTO, error)
	LeaguesForSummoner(encSummonerID string) (*LeaguePositionDTOList, error)
}

// ClientSpectator defines an interface to Spectator API calls
type ClientSpectator interface {
	ActiveGameBySummonerID(summonerID string) (*CurrentGameInfoDTO, error)
	FeaturedGames() (*FeaturedGamesDTO, error)
}

// Client defines the interface for a Riot API Client
type Client interface {
	ClientBase
	ClientChampion
	ClientItem
	ClientLeague
	ClientMatch
	ClientRunesReforged
	ClientSpectator
	ClientSummoner
	ClientSummonerSpells
}
