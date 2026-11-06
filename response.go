package trakt

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Response wraps the standard http.Response with parsed rate limit details.
type Response struct {
	*http.Response
	RateLimit RateLimitInfo
}

// RateLimitInfo contains quota and throttling metadata returned by the Trakt API.
type RateLimitInfo struct {
	Name       string        `json:"name,omitempty"`
	Period     int           `json:"period,omitempty"`
	Limit      int           `json:"limit,omitempty"`
	Remaining  int           `json:"remaining,omitempty"`
	Until      time.Time     `json:"-"`
	RetryAfter time.Duration `json:"-"`
}

type rateLimitJSON struct {
	Name      string `json:"name"`
	Period    int    `json:"period"`
	Limit     int    `json:"limit"`
	Remaining int    `json:"remaining"`
	Until     string `json:"until"`
}

// parseRateLimit extracts rate limit details from response headers.
func parseRateLimit(h http.Header) RateLimitInfo {
	var info RateLimitInfo

	// Parse Retry-After
	if retryAfterStr := h.Get("Retry-After"); retryAfterStr != "" {
		if seconds, err := strconv.Atoi(retryAfterStr); err == nil {
			info.RetryAfter = time.Duration(seconds) * time.Second
		} else if parsedTime, err := http.ParseTime(retryAfterStr); err == nil {
			if duration := time.Until(parsedTime); duration > 0 {
				info.RetryAfter = duration
			}
		}
	}

	// Parse X-Ratelimit JSON header (e.g. {"name":"...","period":300,"limit":1000,"remaining":999,"until":"..."})
	if rawLimit := h.Get("X-Ratelimit"); rawLimit != "" {
		var parsed rateLimitJSON
		if err := json.NewDecoder(strings.NewReader(rawLimit)).Decode(&parsed); err == nil {
			info.Name = parsed.Name
			info.Period = parsed.Period
			info.Limit = parsed.Limit
			info.Remaining = parsed.Remaining
			if parsed.Until != "" {
				if t, err := time.Parse(time.RFC3339, parsed.Until); err == nil {
					info.Until = t
				}
			}
		}
	}

	return info
}
