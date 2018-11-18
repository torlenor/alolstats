package riotclient

import (
	"strconv"
	"strings"
	"time"
)

type rateLimit struct {
	// appRateLimits key = time period in seconds, value = number of calls allowed in that time period
	appRateLimits map[uint32]uint32
	// appRateLimitsCount key = time period in seconds, value = number of used calls in that time period
	appRateLimitsCount map[uint32]uint32
	// Time at which to send the request at earliest
	retryAt time.Time
}

// 100 calls per 1 second
// 1,000 calls per 10 seconds
// 60,000 calls per 10 minutes (600 seconds)
// 360,000 calls per 1 hour (3,600 seconds)

// X-App-Rate-Limit: 100:1,1000:10,60000:600,360000:3600
// X-App-Rate-Limit-Count: 1:1,1:10,1:600,1:3600
// Retry-After: 3 // seconds

// updateRateLimits update the allowed rate limits
func (c *RiotClient) updateAppRateLimits(limits string) {
	// c.log.Debugln("RateLimit: Updating App Rate Limits:", limits)
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	if len(limits) > 0 {
		values := strings.Split(limits, ",")
		c.rateLimit.appRateLimits = make(map[uint32]uint32)
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
				c.rateLimit.appRateLimits[uint32(period)] = uint32(calls)
			}
		}
	}
}

// updateRateLimitsCount update the current rate limit counts
func (c *RiotClient) updateAppRateLimitsCount(counts string) {
	// c.log.Debugln("RateLimit: Updating App Rate Limits Count", counts)
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	if len(counts) > 0 {
		values := strings.Split(counts, ",")
		c.rateLimit.appRateLimitsCount = make(map[uint32]uint32)
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
				c.rateLimit.appRateLimitsCount[uint32(period)] = uint32(calls)
			}
		}
	}
}

// updateRateLimitRetryAt sets the retryAt time based on current time + seconds specified
func (c *RiotClient) updateRateLimitRetryAt(seconds uint32) {
	// c.log.Debugln("RateLimit: Updating App Rate Limits Retry At")
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	c.rateLimit.retryAt = time.Now().Add(time.Second * time.Duration(seconds))
}

func (c *RiotClient) getRateLimitRetryAt() time.Time {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	return c.rateLimit.retryAt
}

func (c *RiotClient) getAdditionalWaitTime() time.Duration {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	var addWaitTime time.Duration
	addWaitTime = 0

	for key, val := range c.rateLimit.appRateLimitsCount {
		maxSlots := c.rateLimit.appRateLimits[key]
		emptySlots := int32(maxSlots) - int32(val)
		if emptySlots < 5 {
			addWaitTime += time.Second * time.Duration(key) / 5
		}
	}

	return addWaitTime
}
