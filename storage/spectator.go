package storage

import (
	"github.com/torlenor/alolstats/riotclient"
)

// GetActiveGameBySummonerID returns the active game (live game) for the given Summoner ID
func (s *Storage) GetActiveGameBySummonerID(summonerID string) (*riotclient.CurrentGameInfoDTO, error) {
	s.log.Debugf("GetActiveGameBySummonerID(%s)", summonerID)

	currentGame, err := s.riotClient.ActiveGameBySummonerID(summonerID)
	if err != nil {
		s.log.Warnln("Could not get data from Riot API:", err)
		return nil, err
	}

	s.log.Debugf("Returned Active Game for SummonerID %s from Riot API", summonerID)
	return currentGame, nil
}

// GetFeaturedGames returns the currently features games from Riot
func (s *Storage) GetFeaturedGames() (*riotclient.FeaturedGamesDTO, error) {
	s.log.Debugf("GetFeaturedGames()")

	featuredGames, err := s.riotClient.FeaturedGames()
	if err != nil {
		s.log.Warnln("Could not get data from Riot API:", err)
		return nil, err
	}

	s.log.Debugf("Returned Featured Games from Riot API")
	return featuredGames, nil
}
