package riotclientv4

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
	"git.abyle.org/hps/alolstats/riotclient"
)

func TestRiotClientV4_Versions(t *testing.T) {
	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config  config.RiotClient
		log     *logrus.Entry
		ddragon *MockRiotClientDD
	}
	tests := []struct {
		name    string
		fields  fields
		wantS   riotclient.Versions
		wantErr bool
	}{
		{
			name: "Test 1 - Receive valid Versions JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					versionsJSON: []byte(`["9.7.1","9.6.1","8.24.1","0.154.2","0.153.2","lolpatch_7.20","lolpatch_7.19"]`),
				},
			},
			wantS: riotclient.Versions{
				"9.7.1",
				"9.6.1",
				"8.24.1",
				"0.154.2",
				"0.153.2",
				"lolpatch_7.20",
				"lolpatch_7.19",
			},
			wantErr: false,
		},
		{
			name:    "Test 2 - Receive invald Versions JSON",
			wantErr: true,
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					versionsJSON: []byte(`[{{{{"9}.7.1","9.6.1","8.24.1","0.154.2","0.153.2","lolpatch_7.20","lolpatch_7.19"]`),
				},
			},
		},
		{
			name:    "Test 3 - API Call returns error",
			wantErr: true,
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					summonerSpellsJSON: []byte(""),
					err:                fmt.Errorf("Some error"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config:  tt.fields.config,
				log:     tt.fields.log,
				ddragon: tt.fields.ddragon,
			}
			gotS, err := c.Versions()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.Versions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.Versions() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
