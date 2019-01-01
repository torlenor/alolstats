// Package riotclientrl is a utility package to check the Rate Limit responses of the Riot API
// and gives suggestions on how long to wait to avoid beeing rate limit
package riotclientrl

import (
	"net/http"
	"reflect"
	"testing"
	"time"

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
		{name: "Test 2 - Update MethodRateLimit",
			fields: fields{
				appRateLimit: newLimit(),
				methodRateLimits: map[string]limit{
					"champions": {
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
				},
				retryAfter: time.Time{},
			},
			args: args{
				header: http.Header{
					"X-Method-Rate-Limit":       []string{"100:1,1000:10,60000:600,360000:3600"},
					"X-Method-Rate-Limit-Count": []string{"1:1,2:10,2:600,2:3600"},
				},
				method: "champions",
			},
		},
		{name: "Test 3 - Update RetryAfter",
			fields: fields{
				appRateLimit:     newLimit(),
				methodRateLimits: make(map[string]limit),
				retryAfter:       now().Add(time.Second * time.Duration(7)),
			},
			args: args{
				header: http.Header{
					"Retry-After": []string{"7"},
				},
			},
		},
		{name: "Test 4 - Update RetryAfter - Header corrupt",
			fields: fields{
				appRateLimit:     newLimit(),
				methodRateLimits: make(map[string]limit),
				retryAfter:       now().Add(time.Second * time.Duration(10)),
			},
			args: args{
				header: http.Header{
					"Retry-After": []string{"sdsdsdsd"},
				},
			},
		},
		{name: "Test 5 - Update AppRateLimit - Header corrupt",
			fields: fields{
				appRateLimit:     newLimit(),
				methodRateLimits: make(map[string]limit),
				retryAfter:       time.Time{},
			},
			args: args{
				header: http.Header{
					"X-App-Rate-Limit":       []string{"134:df"},
					"X-App-Rate-Limit-Count": []string{"1:fff"},
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

func TestRiotClientRL_GetRateLimitRetryAt(t *testing.T) {
	type fields struct {
		appRateLimit     limit
		methodRateLimits map[string]limit
		retryAfter       time.Time
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
		{name: "Test 1 - Plain RetryAfter",
			fields: fields{
				appRateLimit:     newLimit(),
				methodRateLimits: make(map[string]limit),
				retryAfter:       now().Add(time.Second * time.Duration(43)),
			},
			args: args{
				method: "summoners",
			},
			want: now().Add(time.Second * time.Duration(43)),
		},
		{name: "Test 2 - Close to getting rate limited due to AppRateLimit",
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
						600:  60000 - 5,
						3600: 2,
					},
				},
				methodRateLimits: make(map[string]limit),
				retryAfter:       now().Add(time.Second * time.Duration(0)),
			},
			args: args{
				method: "summoners",
			},
			want: now().Add(time.Second * time.Duration(600/5)),
		},
		{name: "Test 3 - Close to getting rate limited due to AppRateLimit - still close",
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
						600:  60000 - 1,
						3600: 2,
					},
				},
				methodRateLimits: make(map[string]limit),
				retryAfter:       now().Add(time.Second * time.Duration(0)),
			},
			args: args{
				method: "summoners",
			},
			want: now().Add(time.Second * time.Duration(600/5)),
		},
		{name: "Test 4 - Close to getting rate limited due to MethodRateLimit",
			fields: fields{
				appRateLimit: newLimit(),
				methodRateLimits: map[string]limit{
					"champions": {
						rateLimits: map[uint32]uint32{
							1:    100,
							10:   1000,
							600:  60000,
							3600: 360000,
						},
						rateLimitsCount: map[uint32]uint32{
							1:    1,
							10:   1000 - 2,
							600:  2,
							3600: 2,
						},
					},
				},
				retryAfter: now().Add(time.Second * time.Duration(13)),
			},
			args: args{
				method: "champions",
			},
			want: now().Add(time.Second * time.Duration(13+10/5)),
		},
		{name: "Test 4 - Close to getting rate limited for a certain Method, but not the one we requested",
			fields: fields{
				appRateLimit: newLimit(),
				methodRateLimits: map[string]limit{
					"champions": {
						rateLimits: map[uint32]uint32{
							1:    100,
							10:   1000,
							600:  60000,
							3600: 360000,
						},
						rateLimitsCount: map[uint32]uint32{
							1:    1,
							10:   1000 - 2,
							600:  2,
							3600: 2,
						},
					},
				},
				retryAfter: now().Add(time.Second * time.Duration(23)),
			},
			args: args{
				method: "summoners",
			},
			want: now().Add(time.Second * time.Duration(23)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New()
			c.appRateLimit = tt.fields.appRateLimit
			c.methodRateLimits = tt.fields.methodRateLimits
			c.retryAfter = tt.fields.retryAfter
			if got := c.GetRateLimitRetryAt(tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientRL.GetRateLimitRetryAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
