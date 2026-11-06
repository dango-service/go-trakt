package trakt

import (
	"net/url"
	"strconv"
)

// SearchOption customizes search queries.
type SearchOption func(url.Values)

// WithSearchLimit sets the maximum number of items per page (default 10, max 100).
func WithSearchLimit(limit int) SearchOption {
	return func(v url.Values) {
		if limit > 0 {
			v.Set("limit", strconv.Itoa(limit))
		}
	}
}

// WithSearchPage sets the page number of results to fetch.
func WithSearchPage(page int) SearchOption {
	return func(v url.Values) {
		if page > 0 {
			v.Set("page", strconv.Itoa(page))
		}
	}
}

// WithSearchExtended requests extended metadata (full overview, ratings, etc.).
func WithSearchExtended(extended bool) SearchOption {
	return func(v url.Values) {
		if extended {
			v.Set("extended", "full")
		}
	}
}

// WithSearchType filters external ID lookups by media type ("movie", "show", "episode").
func WithSearchType(mediaType string) SearchOption {
	return func(v url.Values) {
		if mediaType != "" {
			v.Set("type", mediaType)
		}
	}
}
