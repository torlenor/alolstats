// Package riotclientrl is a utility package to check the Rate Limit responses of the Riot API
// and gives suggestions on how long to wait to avoid beeing rate limit
package riotclientrl

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/logging"
)

var now = time.Now

// A limit represents a Riot API rate limit set of limits and counts
type limit struct {
	// appRateLimits key = time period in seconds, value = number of calls allowed in that time period
	rateLimits map[uint32]uint32
	// appRateLimitsCount key = time period in seconds, value = number of used calls in that time period
	rateLimitsCount map[uint32]uint32
}

func newLimit() limit {
	l := limit{
		rateLimits:      make(map[uint32]uint32),
		rateLimitsCount: make(map[uint32]uint32),
	}
	return l
}

// RiotClientRL Riot LoL API Rate Limit checker
type RiotClientRL struct {
	log *logrus.Entry

	// Stores the Application Rate Limit informations
	appRateLimit limit
	// Stores the Method Rate Limit informations key = method name, value limit
	methodRateLimits map[string]limit
	// Time after which it is allowed to send a new request to the API
	retryAfter time.Time

	rateLimitMutex sync.Mutex
}

// New creates a new Riot LoL API Rate Limit checker
func New() (*RiotClientRL, error) {

	rl := &RiotClientRL{
		log:              logging.Get("RiotClientRL"),
		appRateLimit:     newLimit(),
		methodRateLimits: make(map[string]limit),
	}

	return rl, nil
}

// 100 calls per 1 second
// 1,000 calls per 10 seconds
// 60,000 calls per 10 minutes (600 seconds)
// 360,000 calls per 1 hour (3,600 seconds)

// X-App-Rate-Limit: 100:1,1000:10,60000:600,360000:3600
// X-App-Rate-Limit-Count: 1:1,1:10,1:600,1:3600
// Retry-After: 3 // seconds

// UpdateRateLimits updates the rate limits from the given http header
// header: http header containing rate limit information
// method: method name of the api endpoint used
func (c *RiotClientRL) UpdateRateLimits(header http.Header, method string) {
	if val, ok := header["X-App-Rate-Limit"]; ok {
		if len(val) > 0 {
			c.updateAppRateLimits(val[0])
		}
	}
	if val, ok := header["X-App-Rate-Limit-Count"]; ok {
		if len(val) > 0 {
			c.updateAppRateLimitsCount(val[0])
		}
	}

	if val, ok := header["X-Method-Rate-Limit"]; ok {
		if len(val) > 0 {
			c.updateMethodRateLimits(val[0], method)
		}
	}
	if val, ok := header["X-Method-Rate-Limit-Count"]; ok {
		if len(val) > 0 {
			c.updateMethodRateLimitsCount(val[0], method)
		}
	}

	if val, ok := header["Retry-After"]; ok {
		if len(val) > 0 {
			seconds, err := strconv.ParseUint(val[0], 10, 32)
			if err != nil {
				c.log.Warnf("Could not convert value %s to rate limit retry at seconds", val[0])
				c.updateRateLimitRetryAt(10)
			} else {
				c.updateRateLimitRetryAt(uint32(seconds))
			}
		}
	}
}

func (c *RiotClientRL) updateAppRateLimits(limits string) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	if len(limits) > 0 {
		values := strings.Split(limits, ",")
		c.appRateLimit.rateLimits = make(map[uint32]uint32)
		for _, entry := range values {
			rate := strings.Split(entry, ":")
			if len(rate) == 2 {
				period, err := strconv.ParseUint(rate[1], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit period", rate[1])
					continue
				}
				calls, err := strconv.ParseUint(rate[0], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit count", rate[0])
					continue
				}
				c.appRateLimit.rateLimits[uint32(period)] = uint32(calls)
			}
		}
	}
}

// updateRateLimitsCount update the current rate limit counts
func (c *RiotClientRL) updateAppRateLimitsCount(counts string) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	if len(counts) > 0 {
		values := strings.Split(counts, ",")
		c.appRateLimit.rateLimitsCount = make(map[uint32]uint32)
		for _, entry := range values {
			rate := strings.Split(entry, ":")
			if len(rate) == 2 {
				period, err := strconv.ParseUint(rate[1], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit period", rate[1])
					continue
				}
				calls, err := strconv.ParseUint(rate[0], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit count", rate[0])
					continue
				}
				c.appRateLimit.rateLimitsCount[uint32(period)] = uint32(calls)
			}
		}
	}
}

func (c *RiotClientRL) updateMethodRateLimits(limits string, method string) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	if len(limits) > 0 {
		values := strings.Split(limits, ",")
		if _, ok := c.methodRateLimits[method]; !ok {
			c.methodRateLimits[method] = newLimit()
		}
		for _, entry := range values {
			rate := strings.Split(entry, ":")
			if len(rate) == 2 {
				period, err := strconv.ParseUint(rate[1], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit period", rate[1])
					continue
				}
				calls, err := strconv.ParseUint(rate[0], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit count", rate[0])
					continue
				}
				c.methodRateLimits[method].rateLimits[uint32(period)] = uint32(calls)
			}
		}
	}
}

func (c *RiotClientRL) updateMethodRateLimitsCount(counts string, method string) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	if len(counts) > 0 {
		values := strings.Split(counts, ",")
		if _, ok := c.methodRateLimits[method]; !ok {
			c.methodRateLimits[method] = newLimit()
		}
		for _, entry := range values {
			rate := strings.Split(entry, ":")
			if len(rate) == 2 {
				period, err := strconv.ParseUint(rate[1], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit period", rate[1])
					continue
				}
				calls, err := strconv.ParseUint(rate[0], 10, 32)
				if err != nil {
					c.log.Warnf("Could not convert value %s to rate limit count", rate[0])
					continue
				}
				c.methodRateLimits[method].rateLimitsCount[uint32(period)] = uint32(calls)
			}
		}
	}
}

// updateRateLimitRetryAt sets the retryAt time based on current time + seconds specified
func (c *RiotClientRL) updateRateLimitRetryAt(seconds uint32) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	c.retryAfter = now().Add(time.Second * time.Duration(seconds))
}

// GetRateLimitRetryAt gets a time back at which it is allowed to do a new API call
// method: method name used for that API call
func (c *RiotClientRL) GetRateLimitRetryAt(method string) time.Time {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	if c.retryAfter.Sub(now()) <= time.Duration(0) {
		c.retryAfter = now()
	}

	return c.retryAfter.Add(c.getAdditionalWaitTime(method))
}

func (c *RiotClientRL) getAdditionalWaitTime(method string) time.Duration {
	var addWaitTime time.Duration
	addWaitTime = 0

	for key, val := range c.appRateLimit.rateLimitsCount {
		maxSlots := c.appRateLimit.rateLimits[key]
		emptySlots := int32(maxSlots) - int32(val)
		if emptySlots <= 5 {
			addWaitTime += time.Second * time.Duration(key) / 5
		}
	}

	if limit, ok := c.methodRateLimits[method]; ok {
		for key, val := range limit.rateLimitsCount {
			maxSlots := limit.rateLimits[key]
			emptySlots := int32(maxSlots) - int32(val)
			if emptySlots <= 5 {
				addWaitTime += time.Second * time.Duration(key) / 5
			}
		}
	}

	return addWaitTime
}
