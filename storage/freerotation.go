package storage

import (
	"time"

	"git.abyle.org/hps/alolstats/riotclient"
)

// GetFreeRotation returns the current free rotation from storage
// When forceUpdate is true, it always fetches it vom Riot API, if
// it is false it depends on how old the free rotation is if it gets fetched
// from Riot API
func (s *Storage) GetFreeRotation(forceUpdate bool) riotclient.FreeRotation {
	duration := time.Since(s.backend.GetFreeRotationTimeStamp())
	if duration.Minutes() > float64(s.config.MaxAgeChampionRotation) || forceUpdate {
		freeRotation, err := s.riotClient.ChampionRotations()
		if err != nil {
			s.log.Warnln(err)
			freeRotation, err := s.backend.GetFreeRotation()
			if err != nil {
				s.log.Warnln(err)
				return riotclient.FreeRotation{}
			}
			return *freeRotation
		}
		s.backend.StoreFreeRotation(freeRotation)
		return *freeRotation
	}

	freeRotation, err := s.backend.GetFreeRotation()
	if err != nil {
		s.log.Warnln(err)
		return riotclient.FreeRotation{}
	}
	return *freeRotation
}
