package trakt

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("my-trakt-id")

	if c.clientID != "my-trakt-id" {
		t.Errorf("expected clientID 'my-trakt-id', got %q", c.clientID)
	}
	if c.baseURL.String() != defaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", defaultBaseURL, c.baseURL.String())
	}
	if c.userAgent != defaultUserAgent {
		t.Errorf("expected userAgent %q, got %q", defaultUserAgent, c.userAgent)
	}
	if c.maxBytes != defaultMaxBodyBytes {
		t.Errorf("expected maxBytes %d, got %d", defaultMaxBodyBytes, c.maxBytes)
	}
	if c.Search == nil || c.Shows == nil || c.Movies == nil || c.Sync == nil || c.Scrobble == nil || c.Users == nil || c.Auth == nil {
		t.Error("expected all sub-services to be initialized")
	}
}

func TestNewClient_Options(t *testing.T) {
	customHTTP := &http.Client{Timeout: 5 * time.Second}
	c := NewClient("  trimmed-id  ",
		WithHTTPClient(customHTTP),
		WithBaseURL("https://api-staging.trakt.tv"),
		WithToken("  secret-token  "),
		WithUserAgent("CustomTrakt/1.0"),
		WithMaxBodyBytes(4096),
	)

	if c.clientID != "trimmed-id" {
		t.Errorf("expected trimmed clientID, got %q", c.clientID)
	}
	if c.client != customHTTP {
		t.Errorf("expected custom HTTP client")
	}
	if c.baseURL.String() != "https://api-staging.trakt.tv" {
		t.Errorf("expected custom baseURL, got %q", c.baseURL.String())
	}
	if c.token != "secret-token" {
		t.Errorf("expected trimmed token, got %q", c.token)
	}
	if c.userAgent != "CustomTrakt/1.0" {
		t.Errorf("expected custom userAgent, got %q", c.userAgent)
	}
	if c.maxBytes != 4096 {
		t.Errorf("expected custom maxBytes 4096, got %d", c.maxBytes)
	}
}
