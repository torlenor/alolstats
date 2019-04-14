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
		str := "2018-12-22T13:00:00.000Z"
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

				Timestamp: now(),
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

func TestRiotClientV4_LeaguesForSummoner(t *testing.T) {
	// Override real API call with our fake one
	apiCall = (*RiotClientV4).mockAPICall

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config config.RiotClient
		log    *logrus.Entry
	}
	type args struct {
		summonerID string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantS             *riotclient.LeaguePositionDTOList
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid League Positions JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				summonerID: "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			},
			wantS: &riotclient.LeaguePositionDTOList{
				LeaguePosition: []riotclient.LeaguePositionDTO{
					{
						QueueType:    "RANKED_SOLO_5x5",
						SummonerName: "Suirotas",
						Wins:         392,
						Losses:       339,
						Rank:         "I",
						LeagueName:   "Mordekaiser's Maulers",
						LeagueID:     "a14485bd-709d-3c4c-94f3-b4d6e92a7372",
						Tier:         "GRANDMASTER",
						SummonerID:   "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
						LeaguePoints: 157,

						Timestamp: now(),
					},
				},
			},
			setJSON:           []byte(`[{"queueType":"RANKED_SOLO_5x5","summonerName":"Suirotas","wins":392,"losses":339,"rank":"I","leagueName":"Mordekaiser's Maulers","leagueId":"a14485bd-709d-3c4c-94f3-b4d6e92a7372","tier":"GRANDMASTER","summonerId":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","leaguePoints":157}]`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/league/v4/positions/by-summoner/7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid League Positions JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				summonerID: "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			},
			setJSON:           []byte(`{{{[{"queueType":"RANKED_SOLO_5x5","summonerName":"Suirotas","wins":392,"losses":339,"rank":"I","leagueName":"Mordekaiser's Maulers","leagueId":"a14485bd-709d-3c4c-94f3-b4d6e92a7372","tier":"GRANDMASTER","summonerId":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","leaguePoints":157}]`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/league/v4/positions/by-summoner/7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 3 - API Call failure",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				summonerID: "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			},
			setJSON:           []byte(""),
			setError:          fmt.Errorf("Some API error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/league/v4/positions/by-summoner/7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
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

			gotS, err := c.LeaguesForSummoner(tt.args.summonerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.LeaguesForSummoner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.LeaguesForSummoner() = %v, want %v", gotS, tt.wantS)
			}

			if lastAPICallPath != tt.wantAPICallPath {
				t.Errorf("lastAPICallPath = %v, want %v", lastAPICallPath, tt.wantAPICallPath)
			}
			if lastAPICallBody != tt.wantAPICallBody {
				t.Errorf("lastAPICallBody = %v, want %v", lastAPICallBody, tt.wantAPICallBody)
			}
			if lastAPICallMethod != tt.wantAPICallMethod {
				t.Errorf("lastAPICallMethod = %v, want %v", lastAPICallMethod, tt.wantAPICallMethod)
			}
		})
	}
}
