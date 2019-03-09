package storage

import "github.com/torlenor/alolstats/api"

func (s *Storage) registerAPIChampions(api *api.API) {
	api.AttachModuleGet("/champions", s.championsEndpoint)

	api.AttachModuleGet("/champion/bykey", s.championByKeyEndpoint)
	api.AttachModuleGet("/champion/byid", s.championByIDEndpoint)

	api.AttachModuleGet("/champion-rotations", s.freeRotationEndpoint)
}

func (s *Storage) registerAPIMatch(api *api.API) {
	api.AttachModuleGet("/match", s.getMatchEndpoint)
}

func (s *Storage) registerAPISpectator(api *api.API) {
	api.AttachModuleGet("/active-game", s.getActiveGameBySummonerNameEndpoint)
	api.AttachModuleGet("/featured-games", s.getFeaturedGamesEndpoint)
}

func (s *Storage) registerAPISummoner(api *api.API) {
	api.AttachModuleGet("/summoner/byname", s.summonerByNameEndpoint)
}

func (s *Storage) registerAPIStorage(api *api.API) {
	api.AttachModuleGet("/storage/summary", s.storageSummaryEndpoint)
}

func (s *Storage) registerAPIStats(api *api.API) {
	api.AttachModuleGet("/stats/champion/byid", s.championStatsByIDEndpoint)
	api.AttachModuleGet("/stats/championhistory/byid", s.championStatsHistoryByIDEndpoint)
	api.AttachModuleGet("/stats/champions", s.championStats)
	api.AttachModuleGet("/stats/versions", s.getKnownVersionsEndpoint)
	api.AttachModuleGet("/stats/leagues", s.getStatLeaguesEndpoint)
}

// RegisterAPI registers all endpoints from storage to the RestAPI
func (s *Storage) RegisterAPI(api *api.API) {
	s.registerAPIChampions(api)
	s.registerAPIMatch(api)
	s.registerAPISpectator(api)
	s.registerAPISummoner(api)
	s.registerAPIStorage(api)
	s.registerAPIStats(api)
}
