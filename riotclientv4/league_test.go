package riotclientv4

import (
	"reflect"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/riotclient"

	"github.com/torlenor/alolstats/riotclient/datadragon"
	"github.com/torlenor/alolstats/riotclient/ratelimit"
)

func TestRiotClientV4_MasterLeagueByQueue(t *testing.T) {
	type fields struct {
		config         config.RiotClient
		httpClient     httpClient
		log            *logrus.Entry
		isStarted      bool
		rateLimitMutex sync.Mutex
		workersWG      sync.WaitGroup
		stopWorkers    chan struct{}
		workQueue      workQueue
		ddragon        *riotclientdd.RiotClientDD
		rateLimit      *riotclientrl.RiotClientRL
	}
	type args struct {
		queue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *riotclient.LeagueListDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config:         tt.fields.config,
				httpClient:     tt.fields.httpClient,
				log:            tt.fields.log,
				isStarted:      tt.fields.isStarted,
				rateLimitMutex: tt.fields.rateLimitMutex,
				workersWG:      tt.fields.workersWG,
				stopWorkers:    tt.fields.stopWorkers,
				workQueue:      tt.fields.workQueue,
				ddragon:        tt.fields.ddragon,
				rateLimit:      tt.fields.rateLimit,
			}
			got, err := c.MasterLeagueByQueue(tt.args.queue)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.MasterLeagueByQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientV4.MasterLeagueByQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiotClientV4_GrandMasterLeagueByQueue(t *testing.T) {
	type fields struct {
		config         config.RiotClient
		httpClient     httpClient
		log            *logrus.Entry
		isStarted      bool
		rateLimitMutex sync.Mutex
		workersWG      sync.WaitGroup
		stopWorkers    chan struct{}
		workQueue      workQueue
		ddragon        *riotclientdd.RiotClientDD
		rateLimit      *riotclientrl.RiotClientRL
	}
	type args struct {
		queue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *riotclient.LeagueListDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config:         tt.fields.config,
				httpClient:     tt.fields.httpClient,
				log:            tt.fields.log,
				isStarted:      tt.fields.isStarted,
				rateLimitMutex: tt.fields.rateLimitMutex,
				workersWG:      tt.fields.workersWG,
				stopWorkers:    tt.fields.stopWorkers,
				workQueue:      tt.fields.workQueue,
				ddragon:        tt.fields.ddragon,
				rateLimit:      tt.fields.rateLimit,
			}
			got, err := c.GrandMasterLeagueByQueue(tt.args.queue)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.GrandMasterLeagueByQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientV4.GrandMasterLeagueByQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiotClientV4_ChallengerLeagueByQueue(t *testing.T) {
	type fields struct {
		config         config.RiotClient
		httpClient     httpClient
		log            *logrus.Entry
		isStarted      bool
		rateLimitMutex sync.Mutex
		workersWG      sync.WaitGroup
		stopWorkers    chan struct{}
		workQueue      workQueue
		ddragon        *riotclientdd.RiotClientDD
		rateLimit      *riotclientrl.RiotClientRL
	}
	type args struct {
		queue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *riotclient.LeagueListDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config:         tt.fields.config,
				httpClient:     tt.fields.httpClient,
				log:            tt.fields.log,
				isStarted:      tt.fields.isStarted,
				rateLimitMutex: tt.fields.rateLimitMutex,
				workersWG:      tt.fields.workersWG,
				stopWorkers:    tt.fields.stopWorkers,
				workQueue:      tt.fields.workQueue,
				ddragon:        tt.fields.ddragon,
				rateLimit:      tt.fields.rateLimit,
			}
			got, err := c.ChallengerLeagueByQueue(tt.args.queue)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.ChallengerLeagueByQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientV4.ChallengerLeagueByQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}
