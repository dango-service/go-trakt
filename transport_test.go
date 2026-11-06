package trakt

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTransport_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("trakt-api-version") != "2" {
			t.Errorf("missing or invalid trakt-api-version: %s", r.Header.Get("trakt-api-version"))
		}
		if r.Header.Get("trakt-api-key") != "test-id" {
			t.Errorf("missing trakt-api-key: %s", r.Header.Get("trakt-api-key"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Authorization: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("User-Agent") != defaultUserAgent {
			t.Errorf("unexpected User-Agent: %s", r.Header.Get("User-Agent"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Ratelimit", `{"name":"UNAUTHED_API_GET_LIMIT","period":300,"limit":1000,"remaining":980,"until":"2026-11-06T18:00:00Z"}`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL), WithToken("test-token"))

	var res map[string]string
	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, &res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", res["status"])
	}
	if resp.RateLimit.Limit != 1000 || resp.RateLimit.Remaining != 980 {
		t.Errorf("unexpected rate limit info: %+v", resp.RateLimit)
	}
}

func TestTransport_RateLimitError_Seconds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "12")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message":"Rate limit exceeded"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsRateLimited(err) {
		t.Errorf("expected IsRateLimited to be true for %v", err)
	}
	if resp.RateLimit.RetryAfter != 12*time.Second {
		t.Errorf("expected RetryAfter 12s, got %v", resp.RateLimit.RetryAfter)
	}
}

func TestTransport_RateLimitError_HTTPDate(t *testing.T) {
	targetTime := time.Now().Add(3 * time.Minute).Truncate(time.Second)
	dateStr := targetTime.Format(http.TimeFormat)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", dateStr)
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message":"Too many requests"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.RateLimit.RetryAfter <= 0 {
		t.Errorf("expected positive RetryAfter from HTTP-Date, got %v", resp.RateLimit.RetryAfter)
	}
}

func TestTransport_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"Invalid token"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	_, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsUnauthorized(err) {
		t.Errorf("expected IsUnauthorized to be true for %v", err)
	}
}

func TestTransport_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Item not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	_, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound to be true for %v", err)
	}
}

func TestTransport_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.doRequest(ctx, http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected context cancelled error, got nil")
	}
	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context cancellation in error, got %v", err)
	}
}
