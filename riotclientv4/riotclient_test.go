// Package riotclientv4 provides the Riot API client for API version v4
package riotclientv4

import (
	"net/http"
	"reflect"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
	"git.abyle.org/hps/alolstats/riotclient/datadragon"
	"git.abyle.org/hps/alolstats/riotclient/ratelimit"
)

func TestNewClient(t *testing.T) {
	httpClient := &http.Client{}
	ddragon, _ := riotclientdd.New(httpClient, config.RiotClient{})
	rateLimit, _ := riotclientrl.New()
	client, err := NewClient(httpClient, config.RiotClient{}, ddragon, rateLimit)
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though APIVersion is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1"}, ddragon, rateLimit)
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though Key is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1", Key: "abcd"}, ddragon, rateLimit)
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though Region is missing from config")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v1", Key: "abcd", Region: "euw1"}, ddragon, rateLimit)
	if err == nil || client != nil {
		t.Fatalf("Could get a new client even though APIVersion is wrong")
	}
	client, err = NewClient(&http.Client{}, config.RiotClient{APIVersion: "v4", Key: "abcd", Region: "euw1"}, ddragon, rateLimit)
	if err != nil || client == nil {
		t.Fatalf("Could not get a new client")
	}

	if client.IsRunning() != false {
		t.Fatalf("Client claims to be running even though we did not start it")
	}

	client.Start()
	if client.IsRunning() != true {
		t.Fatalf("Client is not running even though we started it")
	}
	client.Start()
	if client.IsRunning() != true {
		t.Fatalf("Client should still be running if we start it again")
	}

	client.Stop()
	if client.IsRunning() != false {
		t.Fatalf("Client not stopped even though we stopped it")
	}
	client.Stop()
	if client.IsRunning() != false {
		t.Fatalf("Client should still be stopped if we stop it again")
	}
}

func TestRiotClientV4_checkResponseCodeOK(t *testing.T) {
	type fields struct {
		log       *logrus.Entry
		isStarted bool
	}
	type args struct {
		response *http.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test 1 - 200 response OK",
			fields: fields{
				log:       logging.Get("RiotClientV4"),
				isStarted: true,
			},
			args: args{
				response: &http.Response{
					Status:     "OK",
					StatusCode: 200,
				},
			},
			wantErr: false,
		},
		{
			name: "Test 2 - Error response",
			fields: fields{
				log:       logging.Get("RiotClientV4"),
				isStarted: true,
			},
			args: args{
				response: &http.Response{
					Status:     "Not OK",
					StatusCode: 404,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				log:       tt.fields.log,
				isStarted: tt.fields.isStarted,
			}
			if err := c.checkResponseCodeOK(tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.checkResponseCodeOK() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRiotClientV4_checkRateLimited(t *testing.T) {
	type fields struct {
		log       *logrus.Entry
		isStarted bool
	}
	type args struct {
		response *http.Response
		method   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test 1 - Not rate limited",
			fields: fields{
				log:       logging.Get("RiotClientV4"),
				isStarted: true,
			},
			args: args{
				response: &http.Response{
					Status:     "OK",
					StatusCode: 200,
				},
			},
			wantErr: false,
		},
		{
			name: "Test 2 - Rate limited",
			fields: fields{
				log:       logging.Get("RiotClientV4"),
				isStarted: true,
			},
			args: args{
				response: &http.Response{
					Status:     "Rate Limited",
					StatusCode: 429,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				log:       tt.fields.log,
				isStarted: tt.fields.isStarted,
			}
			if err := c.checkRateLimited(tt.args.response, tt.args.method); (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.checkRateLimited() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRiotClientV4_realAPICall(t *testing.T) {
	type fields struct {
		config      config.RiotClient
		httpClient  httpClient
		log         *logrus.Entry
		isStarted   bool
		workersWG   sync.WaitGroup
		stopWorkers chan struct{}
		workQueue   workQueue
		ddragon     *riotclientdd.RiotClientDD
		rateLimit   *riotclientrl.RiotClientRL
	}
	type args struct {
		path   string
		method string
		body   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config:      tt.fields.config,
				httpClient:  tt.fields.httpClient,
				log:         tt.fields.log,
				isStarted:   tt.fields.isStarted,
				workersWG:   tt.fields.workersWG,
				stopWorkers: tt.fields.stopWorkers,
				workQueue:   tt.fields.workQueue,
				ddragon:     tt.fields.ddragon,
				rateLimit:   tt.fields.rateLimit,
			}
			gotR, err := c.realAPICall(tt.args.path, tt.args.method, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.realAPICall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("RiotClientV4.realAPICall() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
