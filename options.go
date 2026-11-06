package trakt

import (
	"net/http"
	"net/url"
	"strings"
)

// ClientOption configures a Client instance.
type ClientOption func(*Client)

// WithHTTPClient overrides the default http.Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.client = httpClient
		}
	}
}

// WithBaseURL overrides the default Trakt API base URL (e.g. for staging or tests).
func WithBaseURL(rawURL string) ClientOption {
	return func(c *Client) {
		if parsed, err := url.Parse(rawURL); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			c.baseURL = parsed
		}
	}
}

// WithToken sets the user OAuth2 access token for authenticated requests.
// Leading and trailing whitespace is stripped.
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = strings.TrimSpace(token)
	}
}

// WithUserAgent overrides the default User-Agent header.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = strings.TrimSpace(userAgent)
	}
}

// WithMaxBodyBytes overrides the maximum bytes read from a response body.
func WithMaxBodyBytes(limit int64) ClientOption {
	return func(c *Client) {
		if limit > 0 {
			c.maxBytes = limit
		}
	}
}
