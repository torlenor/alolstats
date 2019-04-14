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

func TestRiotClientV4_SummonerByName(t *testing.T) {
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
		name string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantS             *riotclient.SummonerDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid Summoner JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				name: "Suirotas",
			},
			wantS: &riotclient.SummonerDTO{
				ProfileIcon:   3176,
				Name:          "Suirotas",
				PuuID:         "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
				SummonerLevel: 93,
				AccountID:     "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				ID:            "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
				RevisionDate:  1545818141000,
				Timestamp:     now(),
			},
			setJSON:           []byte(`{"profileIconId":3176,"name":"Suirotas","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/Suirotas",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid Summoner JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				name: "SomeSummoner",
			},
			setJSON:           []byte(`{{{{{"profileIconId":3176,"name":"SomeSummoner","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/SomeSummoner",
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
				name: "AnotherSummoner",
			},
			setJSON:           []byte(""),
			setError:          fmt.Errorf("Some API error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/AnotherSummoner",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 4 - Received empty JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				name: "AnotherSummoner",
			},
			setJSON:           []byte("{}"),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/AnotherSummoner",
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

			gotS, err := c.SummonerByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.SummonerByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.SummonerByName() = %v, want %v", gotS, tt.wantS)
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

func TestRiotClientV4_SummonerByAccountID(t *testing.T) {
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
		accountID string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantS             *riotclient.SummonerDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid Summoner JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			},
			wantS: &riotclient.SummonerDTO{
				ProfileIcon:   3176,
				Name:          "Suirotas",
				PuuID:         "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
				SummonerLevel: 93,
				AccountID:     "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				ID:            "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
				RevisionDate:  1545818141000,
				Timestamp:     now(),
			},
			setJSON:           []byte(`{"profileIconId":3176,"name":"Suirotas","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid Summoner JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			},
			setJSON:           []byte(`{{{{{"profileIconId":3176,"name":"SomeSummoner","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
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
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			},
			setJSON:           []byte(""),
			setError:          fmt.Errorf("Some API error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 4 - Received empty JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			},
			setJSON:           []byte("{}"),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
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

			gotS, err := c.SummonerByAccountID(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.SummonerByAccountID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.SummonerByAccountID() = %v, want %v", gotS, tt.wantS)
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

func TestRiotClientV4_SummonerBySummonerID(t *testing.T) {
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
		wantS             *riotclient.SummonerDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid Summoner JSON",
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
			wantS: &riotclient.SummonerDTO{
				ProfileIcon:   3176,
				Name:          "Suirotas",
				PuuID:         "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
				SummonerLevel: 93,
				AccountID:     "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				ID:            "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
				RevisionDate:  1545818141000,
				Timestamp:     now(),
			},
			setJSON:           []byte(`{"profileIconId":3176,"name":"Suirotas","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid Summoner JSON",
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
			setJSON:           []byte(`{{{{{"profileIconId":3176,"name":"SomeSummoner","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
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
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 4 - Received empty JSON",
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
			setJSON:           []byte("{}"),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
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

			gotS, err := c.SummonerBySummonerID(tt.args.summonerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.SummonerBySummonerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.SummonerBySummonerID() = %v, want %v", gotS, tt.wantS)
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

func TestRiotClientV4_SummonerByPUUID(t *testing.T) {
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
		PUUID string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantS             *riotclient.SummonerDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid Summoner JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				PUUID: "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
			},
			wantS: &riotclient.SummonerDTO{
				ProfileIcon:   3176,
				Name:          "Suirotas",
				PuuID:         "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
				SummonerLevel: 93,
				AccountID:     "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				ID:            "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
				RevisionDate:  1545818141000,
				Timestamp:     now(),
			},
			setJSON:           []byte(`{"profileIconId":3176,"name":"Suirotas","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid Summoner JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				PUUID: "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
			},
			setJSON:           []byte(`{{{{{"profileIconId":3176,"name":"SomeSummoner","puuid":"Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA","summonerLevel":93,"accountId":"C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw","id":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE","revisionDate":1545818141000}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
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
				PUUID: "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
			},
			setJSON:           []byte(""),
			setError:          fmt.Errorf("Some API error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 4 - Received empty JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				PUUID: "Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
			},
			setJSON:           []byte("{}"),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/Mmu-pkyHrl6Cx5cz9PaHZym5NUF1Wk-cLNj8k8djtcV9zulj_dow-dtKnWYfQ8VbTHeM7uDR15oirA",
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

			gotS, err := c.SummonerByPUUID(tt.args.PUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.SummonerByPUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.SummonerByPUUID() = %v, want %v", gotS, tt.wantS)
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
