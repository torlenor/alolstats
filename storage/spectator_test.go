package storage

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/riotclient"
)

func TestStorage_GetActiveGameBySummonerID(t *testing.T) {
	type fields struct {
		config     config.LoLStorage
		riotClient riotclient.Client
		log        *logrus.Entry
		stats      stats
		backend    Backend
	}
	type args struct {
		summonerID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *riotclient.CurrentGameInfoDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				config:     tt.fields.config,
				riotClient: tt.fields.riotClient,
				log:        tt.fields.log,
				stats:      tt.fields.stats,
				backend:    tt.fields.backend,
			}
			got, err := s.GetActiveGameBySummonerID(tt.args.summonerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetActiveGameBySummonerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.GetActiveGameBySummonerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetFeaturedGames(t *testing.T) {
	type fields struct {
		config     config.LoLStorage
		riotClient riotclient.Client
		log        *logrus.Entry
		stats      stats
		backend    Backend
	}
	tests := []struct {
		name    string
		fields  fields
		want    *riotclient.FeaturedGamesDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				config:     tt.fields.config,
				riotClient: tt.fields.riotClient,
				log:        tt.fields.log,
				stats:      tt.fields.stats,
				backend:    tt.fields.backend,
			}
			got, err := s.GetFeaturedGames()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetFeaturedGames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.GetFeaturedGames() = %v, want %v", got, tt.want)
			}
		})
	}
}
