package riotclientdd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
)

func Test_checkConfig(t *testing.T) {
	type args struct {
		cfg config.RiotClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1 - Region empty in config",
			args: args{
				cfg: config.RiotClient{
					Region: "",
				}},
			wantErr: true,
		},
		{
			name: "Test 2 - Region set",
			args: args{
				cfg: config.RiotClient{
					Region: "euw",
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkConfig(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("checkConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		httpClient *http.Client
		cfg        config.RiotClient
	}
	tests := []struct {
		name    string
		args    args
		want    *RiotClientDD
		wantErr bool
	}{
		{
			name: "Test 1 - Create Data Dragon struct",
			args: args{
				httpClient: &http.Client{},
				cfg: config.RiotClient{
					Region: "EuW1",
				}},
			want: &RiotClientDD{
				httpClient: &http.Client{},
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
			},
			wantErr: false,
		},
		{
			name: "Test 2 - Region not set",
			args: args{
				httpClient: &http.Client{},
				cfg: config.RiotClient{
					Region: "",
				}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 3 - Region too short",
			args: args{
				httpClient: &http.Client{},
				cfg: config.RiotClient{
					Region: "e",
				}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.httpClient, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiotClientDD_getRegion(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *http.Client
		log        *logrus.Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1 - Get a set region",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
			},
			want: "euw",
		},
		{
			name: "Test 2 - Get another set region",
			fields: fields{
				config: config.RiotClient{
					Region: "na1",
				},
				log: logging.Get("RiotClientDD"),
			},
			want: "na",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			if got := c.getRegion(); got != tt.want {
				t.Errorf("RiotClientDD.getRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiotClientDD_getVersions(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		want    *currentVersions
		wantErr bool
	}{
		{
			name: "Test 1 - Get a valid versions json",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"item":"8.24.1","rune":"7.23.1","mastery":"7.23.1",
						"summoner":"8.24.1","champion":"8.24.1","profileicon":"8.24.1","map":"8.24.1","language":"8.24.1","sticker":"8.24.1"},
						"v":"8.24.1","l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn","dd":"8.24.1","lg":"8.24.1",
						"css":"8.24.1","profileiconmax":28,"store":null}`))),
					},
				},
			},
			want: &currentVersions{
				N: N{
					Item:        "8.24.1",
					Rune:        "7.23.1",
					Mastery:     "7.23.1",
					Summoner:    "8.24.1",
					Champion:    "8.24.1",
					Profileicon: "8.24.1",
					Map:         "8.24.1",
					Language:    "8.24.1",
					Sticker:     "8.24.1",
				},
				V:              "8.24.1",
				L:              "en_GB",
				Cdn:            "https://ddragon.leagueoflegends.com/cdn",
				Dd:             "8.24.1",
				Lg:             "8.24.1",
				CSS:            "8.24.1",
				Profileiconmax: 28,
				Store:          nil,
			},
		},
		{
			name: "Test 2 - Invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: nil,
					err:      fmt.Errorf("Could not get Versions"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 3 - Invalid json",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"item":"8.24.1","rune":"7.23.1","mastery":"7.23.1",
						"summoner":"8.24.1","champion":"8.24.1","profile}`))),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.getVersions()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.getVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.getVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonChampions(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	tests := []struct {
		name               string
		fields             fields
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a champions json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.24.1" + "/data/" + "en_GB" + "/champion.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.23.1" + "/data/" + "en_US" + "/champion.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonChampions()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonChampions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonChampions() = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonChampions() = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonChampions() = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonChampionsSpecificVersionLanguage(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	type args struct {
		gameVersion string
		language    string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a champions json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			args: args{
				gameVersion: "9.7.1",
				language:    "de_DE",
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.7.1" + "/data/" + "de_DE" + "/champion.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			args: args{
				gameVersion: "8.7.1",
				language:    "ab_cd",
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.7.1" + "/data/" + "ab_cd" + "/champion.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonChampionsSpecificVersionLanguage(tt.args.gameVersion, tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonChampionsSpecificVersionLanguage(args) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonChampionsSpecificVersionLanguage(args) = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonChampionsSpecificVersionLanguage(args) = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonChampionsSpecificVersionLanguage(args) = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}

func TestRiotClientDD_GetLoLVersions(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	tests := []struct {
		name              string
		fields            fields
		want              []byte
		wantErr           bool
		wantRequestString string
	}{
		{
			name: "Test 1 - Get a versions json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`blabla`))),
					},
				},
			},
			want:              []byte("blabla"),
			wantErr:           false,
			wantRequestString: "https://ddragon.leagueoflegends.com/api/versions.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error downloading file"),
				},
			},
			want:              nil,
			wantErr:           true,
			wantRequestString: "https://ddragon.leagueoflegends.com/api/versions.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetLoLVersions()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetLoLVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetLoLVersions() = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetLoLVersions() = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonSummonerSpells(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	tests := []struct {
		name               string
		fields             fields
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a Summoner Spells json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1","summoner":"8.21.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.21.1" + "/data/" + "en_GB" + "/summoner.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1","summoner":"8.22.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.22.1" + "/data/" + "en_US" + "/summoner.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonSummonerSpells()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpells() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpells() = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpells() = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpells() = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonSummonerSpellsSpecificVersionLanguage(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	type args struct {
		gameVersion string
		language    string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a Summoner Spells json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1","summoner":"8.22.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			args: args{
				gameVersion: "9.7.1",
				language:    "de_DE",
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.7.1" + "/data/" + "de_DE" + "/summoner.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1","summoner":"8.22.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			args: args{
				gameVersion: "8.7.1",
				language:    "ab_cd",
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.7.1" + "/data/" + "ab_cd" + "/summoner.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonSummonerSpellsSpecificVersionLanguage(tt.args.gameVersion, tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpellsSpecificVersionLanguage(args) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpellsSpecificVersionLanguage(args) = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpellsSpecificVersionLanguage(args) = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonSummonerSpellsSpecificVersionLanguage(args) = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonItems(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	tests := []struct {
		name               string
		fields             fields
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a Items json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1","summoner":"8.21.1","item":"9.7.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.7.1" + "/data/" + "en_GB" + "/item.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1","summoner":"8.22.1","item":"9.1.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.1.1" + "/data/" + "en_US" + "/item.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonItems() = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonItems() = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonItems() = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonItemsSpecificVersionLanguage(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	type args struct {
		gameVersion string
		language    string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a Items json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1","summoner":"8.22.1","item":"9.1.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			args: args{
				gameVersion: "9.7.1",
				language:    "de_DE",
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.7.1" + "/data/" + "de_DE" + "/item.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1","summoner":"8.22.1","item":"9.2.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			args: args{
				gameVersion: "8.7.1",
				language:    "ab_cd",
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.7.1" + "/data/" + "ab_cd" + "/item.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonItemsSpecificVersionLanguage(tt.args.gameVersion, tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonItemsSpecificVersionLanguage(args) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonItemsSpecificVersionLanguage(args) = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonItemsSpecificVersionLanguage(args) = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonItemsSpecificVersionLanguage(args) = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonRunesReforged(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	tests := []struct {
		name               string
		fields             fields
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a RunesReforged json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1","summoner":"8.21.1","item":"9.7.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.7.1" + "/data/" + "en_GB" + "/runesReforged.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1","summoner":"8.22.1","item":"9.1.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.1.1" + "/data/" + "en_US" + "/runesReforged.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonRunesReforged()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforged() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforged() = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforged() = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforged() = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}

func TestRiotClientDD_GetDataDragonRunesReforgedSpecificVersionLanguage(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *mockHTTPClient
		log        *logrus.Entry
	}
	type args struct {
		gameVersion string
		language    string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               []byte
		wantErr            bool
		wantRequestString  string
		wantRequestString2 string
	}{
		{
			name: "Test 1 - Get a RunesReforged json string",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.24.1","summoner":"8.22.1","item":"9.1.1"},"l":"en_GB","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			args: args{
				gameVersion: "9.7.1",
				language:    "de_DE",
			},
			want:               []byte("fdsfsdf"),
			wantErr:            false,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "9.7.1" + "/data/" + "de_DE" + "/runesReforged.json",
		},
		{
			name: "Test 2 - Get an invalid response",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"n":{"champion":"8.23.1","summoner":"8.22.1","item":"9.2.1"},"l":"en_US","cdn":"https://ddragon.leagueoflegends.com/cdn"}`))),
					},
					response2: &http.Response{
						StatusCode: 404,
					},
					err2: fmt.Errorf("Error downloading file"),
				},
			},
			args: args{
				gameVersion: "8.7.1",
				language:    "ab_cd",
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "https://ddragon.leagueoflegends.com/cdn/" + "8.7.1" + "/data/" + "ab_cd" + "/runesReforged.json",
		},
		{
			name: "Test 3 - Get an invalid response from versions",
			fields: fields{
				config: config.RiotClient{
					Region: "euw1",
				},
				log: logging.Get("RiotClientDD"),
				httpClient: &mockHTTPClient{
					response: &http.Response{
						StatusCode: 404,
					},
					err: fmt.Errorf("Error getting versions"),
					response2: &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`fdsfsdf`))),
					},
				},
			},
			want:               nil,
			wantErr:            true,
			wantRequestString:  "https://ddragon.leagueoflegends.com/realms/euw.json",
			wantRequestString2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.GetDataDragonRunesReforgedSpecificVersionLanguage(tt.args.gameVersion, tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforgedSpecificVersionLanguage(args) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforgedSpecificVersionLanguage(args) = %v, want %v", got, tt.want)
			}
			if tt.fields.httpClient.getURL != tt.wantRequestString {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforgedSpecificVersionLanguage(args) = %v, want request string %v", tt.fields.httpClient.getURL, tt.wantRequestString)
			}
			if tt.fields.httpClient.getURL2 != tt.wantRequestString2 {
				t.Errorf("RiotClientDD.GetDataDragonRunesReforgedSpecificVersionLanguage(args) = %v, want request string 2 %v", tt.fields.httpClient.getURL2, tt.wantRequestString2)
			}
		})
	}
}
