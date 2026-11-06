package trakt

import (
	"net/url"
	"strconv"
)

// SyncOption customizes a sync query.
type SyncOption func(url.Values)

// WithSyncPage sets the page number of sync results to fetch.
func WithSyncPage(page int) SyncOption {
	return func(v url.Values) {
		if page > 0 {
			v.Set("page", strconv.Itoa(page))
		}
	}
}

// WithSyncLimit sets the limit of items returned per page.
func WithSyncLimit(limit int) SyncOption {
	return func(v url.Values) {
		if limit > 0 {
			v.Set("limit", strconv.Itoa(limit))
		}
	}
}
