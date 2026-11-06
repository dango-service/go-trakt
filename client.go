package trakt

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL      = "https://api.trakt.tv"
	defaultUserAgent    = "go-trakt/0.1.0 (+https://github.com/dango-service/go-trakt)"
	defaultMaxBodyBytes = 10 * 1024 * 1024 // 10 MiB
	apiVersion          = "2"
)

// Client coordinates communication with the Trakt.tv API.
type Client struct {
	client    *http.Client
	baseURL   *url.URL
	clientID  string
	token     string
	userAgent string
	maxBytes  int64

	Search   *SearchService
	Shows    *ShowsService
	Movies   *MoviesService
	Sync     *SyncService
	Scrobble *ScrobbleService
	Users    *UsersService
	Auth     *AuthService
}

// NewClient initializes a new Trakt API client with the given client ID.
func NewClient(clientID string, opts ...ClientOption) *Client {
	parsedBaseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   parsedBaseURL,
		clientID:  strings.TrimSpace(clientID),
		userAgent: defaultUserAgent,
		maxBytes:  defaultMaxBodyBytes,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	c.Search = &SearchService{client: c}
	c.Shows = &ShowsService{client: c}
	c.Movies = &MoviesService{client: c}
	c.Sync = &SyncService{client: c}
	c.Scrobble = &ScrobbleService{client: c}
	c.Users = &UsersService{client: c}
	c.Auth = &AuthService{client: c}

	return c
}
