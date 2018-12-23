package riotclientdd

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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

func TestRiotClientDD_downloadFile(t *testing.T) {
	type fields struct {
		config     config.RiotClient
		httpClient *http.Client
		log        *logrus.Entry
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientDD{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
				log:        tt.fields.log,
			}
			got, err := c.downloadFile(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientDD.downloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientDD.downloadFile() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		httpClient *http.Client
		log        *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		want    *currentVersions
		wantErr bool
	}{
		// TODO: Add test cases.
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
		httpClient *http.Client
		log        *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
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
		})
	}
}
