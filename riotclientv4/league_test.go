package riotclientv4

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/riotclient"
)

func TestRiotClientV4_leagueByQueue(t *testing.T) {
	// Override real API call with our fake one
	apiCall = (*RiotClientV4).mockAPICall

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.0000"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config config.RiotClient
		log    *logrus.Entry
	}
	type args struct {
		leagueEndPoint string
		queue          string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *riotclient.LeagueListDTO
		wantErr  bool
		setJSON  []byte
		setError error
	}{
		{
			name: "Test 1 - Valid request - Receive valid Leagues JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				leagueEndPoint: "mastersleague",
				queue:          "RANKED_SOLO_5x5",
			},
			want: &riotclient.LeagueListDTO{
				Tier:     "GRANDMASTER",
				LeagueID: "a14485bd-709d-3c4c-94f3-b4d6e92a7372",
				Entries: []riotclient.LeagueItemDTO{
					{
						SummonerName: "S1 Atrocez",
						Wins:         544,
						Losses:       489,
						Rank:         "I",
						SummonerID:   "CM1D8AS2W9uPKJa2AjCaUkXSakYqmYaYvosL-tSjBaMD8ck",
						LeaguePoints: 191,
					},
					{
						SummonerName: "MINDFLAY",
						Wins:         154,
						Losses:       97,
						Rank:         "I",
						SummonerID:   "LbWA5BBedoFOYOzHY56oxGaqrTW9NMYEP9Ftfi6UHtxk_LQ",
						LeaguePoints: 32,
					},
				},
			},
			wantErr: false,
			setJSON: []byte(`{"tier":"GRANDMASTER","leagueId":"a14485bd-709d-3c4c-94f3-b4d6e92a7372","entries":[{"summonerName":"S1 Atrocez","wins":544,"losses":489,"rank":"I","summonerId":"CM1D8AS2W9uPKJa2AjCaUkXSakYqmYaYvosL-tSjBaMD8ck","leaguePoints":191},{"summonerName":"MINDFLAY","wins":154,"losses":97,"rank":"I","summonerId":"LbWA5BBedoFOYOzHY56oxGaqrTW9NMYEP9Ftfi6UHtxk_LQ","leaguePoints":32}]}`),
		},
		{
			name: "Test 2 - Valid request - Receive invalid Leagues JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				leagueEndPoint: "challengersleague",
				queue:          "RANKED_SOLO_5x5",
			},
			want:    nil,
			wantErr: true,

			setJSON: []byte(`{{{{"tier":"GRANDMASTER","leagueId":"a14485bd-709d-3c4c-94f3-b4d6e92a7372","entries":[{"summonerName":"S1 Atrocez","wins":544,"losses":489,"rank":"I","summonerId":"CM1D8AS2W9uPKJa2AjCaUkXSakYqmYaYvosL-tSjBaMD8ck","leaguePoints":191},{"summonerName":"MINDFLAY","wins":154,"losses":97,"rank":"I","summonerId":"LbWA5BBedoFOYOzHY56oxGaqrTW9NMYEP9Ftfi6UHtxk_LQ","leaguePoints":32}]}`),
		},
		{
			name: "Test 3 - Invalid queue specified",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				leagueEndPoint: "grandmastersleague",
				queue:          "RANKED_BLA_BLA",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 4 - API Call returns error",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				leagueEndPoint: "grandmastersleague",
				queue:          "RANKED_SOLO_5x5",
			},
			want:     nil,
			wantErr:  true,
			setError: fmt.Errorf("Some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config: tt.fields.config,
				log:    tt.fields.log,
			}

			apiCallReturnJSON = tt.setJSON
			apiCallReturnErr = tt.setError

			{
				got, err := c.leagueByQueue(tt.args.leagueEndPoint, tt.args.queue)
				if (err != nil) != tt.wantErr {
					t.Errorf("RiotClientV4.leagueByQueue() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("RiotClientV4.leagueByQueue() = %v, want %v", got, tt.want)
				}
			}
			{
				got, err := c.LeagueByQueue(tt.args.leagueEndPoint, tt.args.queue)
				if (err != nil) != tt.wantErr {
					t.Errorf("RiotClientV4.LeagueByQueue() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("RiotClientV4.LeagueByQueue() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
