// Package riotclientrl is a utility package to check the Rate Limit responses of the Riot API
// and gives suggestions on how long to wait to avoid beeing rate limit
package riotclientrl

import (
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/logging"
)

func Test_newLimit(t *testing.T) {
	tests := []struct {
		name string
		want limit
	}{
		{"Succesful creation", limit{rateLimits: make(map[uint32]uint32), rateLimitsCount: make(map[uint32]uint32)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLimit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *RiotClientRL
		wantErr bool
	}{
		{"Succesful creation", &RiotClientRL{
			log:              logging.Get("RiotClientRL"),
			appRateLimit:     newLimit(),
			methodRateLimits: make(map[string]limit),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
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

func TestRiotClientRL_UpdateRateLimits(t *testing.T) {
	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.0000"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
	}
	type args struct {
		header http.Header
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "Test 1 - Update AppRateLimit",
			fields: fields{
				appRateLimit: limit{
					rateLimits: map[uint32]uint32{
						1:    100,
						10:   1000,
						600:  60000,
						3600: 360000,
					},
					rateLimitsCount: map[uint32]uint32{
						1:    1,
						10:   2,
						600:  2,
						3600: 2,
					},
				},
				methodRateLimits: make(map[string]limit),
				retryAfter:       time.Time{},
			},
			args: args{
				header: http.Header{
					"X-App-Rate-Limit":       []string{"100:1,1000:10,60000:600,360000:3600"},
					"X-App-Rate-Limit-Count": []string{"1:1,2:10,2:600,2:3600"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New()
			c.UpdateRateLimits(tt.args.header, tt.args.method)
			if !reflect.DeepEqual(c.appRateLimit, tt.fields.appRateLimit) {
				t.Errorf("New() = %v, want %v", c.appRateLimit, tt.fields.appRateLimit)
			}
			if !reflect.DeepEqual(c.methodRateLimits, tt.fields.methodRateLimits) {
				t.Errorf("New() = %v, want %v", c.methodRateLimits, tt.fields.methodRateLimits)
			}
			if !reflect.DeepEqual(c.retryAfter, tt.fields.retryAfter) {
				t.Errorf("New() = %v, want %v", c.retryAfter, tt.fields.retryAfter)
			}
		})
	}
}

func TestRiotClientRL_updateAppRateLimits(t *testing.T) {
	type fields struct {
		log              *logrus.Entry
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
		rateLimitMutex   sync.Mutex
	}
	type args struct {
		limits string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientRL{
				log:              tt.fields.log,
				appRateLimit:     tt.fields.appRateLimit,
				methodRateLimits: tt.fields.methodRateLimits,
				retryAfter:       tt.fields.retryAfter,
				rateLimitMutex:   tt.fields.rateLimitMutex,
			}
			c.updateAppRateLimits(tt.args.limits)
		})
	}
}

func TestRiotClientRL_updateAppRateLimitsCount(t *testing.T) {
	type fields struct {
		log              *logrus.Entry
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
		rateLimitMutex   sync.Mutex
	}
	type args struct {
		counts string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientRL{
				log:              tt.fields.log,
				appRateLimit:     tt.fields.appRateLimit,
				methodRateLimits: tt.fields.methodRateLimits,
				retryAfter:       tt.fields.retryAfter,
				rateLimitMutex:   tt.fields.rateLimitMutex,
			}
			c.updateAppRateLimitsCount(tt.args.counts)
		})
	}
}

func TestRiotClientRL_updateMethodRateLimits(t *testing.T) {
	type fields struct {
		log              *logrus.Entry
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
		rateLimitMutex   sync.Mutex
	}
	type args struct {
		limits string
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientRL{
				log:              tt.fields.log,
				appRateLimit:     tt.fields.appRateLimit,
				methodRateLimits: tt.fields.methodRateLimits,
				retryAfter:       tt.fields.retryAfter,
				rateLimitMutex:   tt.fields.rateLimitMutex,
			}
			c.updateMethodRateLimits(tt.args.limits, tt.args.method)
		})
	}
}

func TestRiotClientRL_updateMethodRateLimitsCount(t *testing.T) {
	type fields struct {
		log              *logrus.Entry
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
		rateLimitMutex   sync.Mutex
	}
	type args struct {
		counts string
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientRL{
				log:              tt.fields.log,
				appRateLimit:     tt.fields.appRateLimit,
				methodRateLimits: tt.fields.methodRateLimits,
				retryAfter:       tt.fields.retryAfter,
				rateLimitMutex:   tt.fields.rateLimitMutex,
			}
			c.updateMethodRateLimitsCount(tt.args.counts, tt.args.method)
		})
	}
}

func TestRiotClientRL_updateRateLimitRetryAt(t *testing.T) {
	type fields struct {
		log              *logrus.Entry
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
		rateLimitMutex   sync.Mutex
	}
	type args struct {
		seconds uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientRL{
				log:              tt.fields.log,
				appRateLimit:     tt.fields.appRateLimit,
				methodRateLimits: tt.fields.methodRateLimits,
				retryAfter:       tt.fields.retryAfter,
				rateLimitMutex:   tt.fields.rateLimitMutex,
			}
			c.updateRateLimitRetryAt(tt.args.seconds)
		})
	}
}

func TestRiotClientRL_GetRateLimitRetryAt(t *testing.T) {
	type fields struct {
		log              *logrus.Entry
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
		rateLimitMutex   sync.Mutex
	}
	type args struct {
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientRL{
				log:              tt.fields.log,
				appRateLimit:     tt.fields.appRateLimit,
				methodRateLimits: tt.fields.methodRateLimits,
				retryAfter:       tt.fields.retryAfter,
				rateLimitMutex:   tt.fields.rateLimitMutex,
			}
			if got := c.GetRateLimitRetryAt(tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientRL.GetRateLimitRetryAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiotClientRL_getAdditionalWaitTime(t *testing.T) {
	type fields struct {
		log              *logrus.Entry
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
		rateLimitMutex   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientRL{
				log:              tt.fields.log,
				appRateLimit:     tt.fields.appRateLimit,
				methodRateLimits: tt.fields.methodRateLimits,
				retryAfter:       tt.fields.retryAfter,
				rateLimitMutex:   tt.fields.rateLimitMutex,
			}
			if got := c.getAdditionalWaitTime(); got != tt.want {
				t.Errorf("RiotClientRL.getAdditionalWaitTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
