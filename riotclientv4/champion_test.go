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

func TestRiotClientV4_ChampionRotations(t *testing.T) {
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
	tests := []struct {
		name     string
		fields   fields
		want     *riotclient.FreeRotation
		wantErr  bool
		setJSON  []byte
		setError error
	}{
		{
			name: "Test 1 - Receive valid Champions Rotation JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			want: &riotclient.FreeRotation{
				FreeChampionIds:              []int{7, 17, 74, 79, 83, 104, 115, 164, 267, 420, 421, 497, 498, 518},
				FreeChampionIdsForNewPlayers: []int{18, 81, 92, 141, 37, 238, 19, 45, 25, 64},
				MaxNewPlayerLevel:            10,
				Timestamp:                    now(),
			},
			wantErr: false,
			setJSON: []byte(`{"freeChampionIds":[7,17,74,79,83,104,115,164,267,420,421,497,498,518],"freeChampionIdsForNewPlayers":[18,81,92,141,37,238,19,45,25,64],"maxNewPlayerLevel":10}`),
		},
		{
			name: "Test 2 - Receive invald Champions Rotation JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			want:    nil,
			wantErr: true,
			setJSON: []byte(`{{{{"freeChampionIds":[7,17,74,79,83,104,115,164,267,420,421,497,498,518],"freeChampionIdsForNewPlayers":[18,81,92,141,37,238,19,45,25,64],"maxNewPlayerLevel":10}`),
		},
		{
			name: "Test 3 - Receive empty Champions Rotation JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			want:    nil,
			wantErr: true,
			setJSON: []byte(`{"freeChampionIds":[],"freeChampionIdsForNewPlayers":[],"maxNewPlayerLevel":10}`),
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
			want:     nil,
			wantErr:  true,
			setJSON:  []byte(``),
			setError: fmt.Errorf("Some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiCallReturnJSON = tt.setJSON
			apiCallReturnErr = tt.setError

			c := &RiotClientV4{
				config: tt.fields.config,
				log:    tt.fields.log,
			}
			got, err := c.ChampionRotations()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.ChampionRotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientV4.ChampionRotations() = %v, want %v", got, tt.want)
			}
		})
	}
}
