package trakt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSearchService_Search(t *testing.T) {
	fixture, err := os.ReadFile("testdata/search_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/show" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "Breaking Bad" {
			t.Errorf("unexpected query param: %s", r.URL.Query().Get("query"))
		}
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("unexpected limit param: %s", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	results, _, err := client.Search.Search(context.Background(), "show", "Breaking Bad", WithSearchLimit(5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	item := results[0]
	if item.Type != "show" || item.Show == nil {
		t.Fatalf("expected show result, got %+v", item)
	}
	if item.Show.Title != "Breaking Bad" {
		t.Errorf("expected title 'Breaking Bad', got %q", item.Show.Title)
	}
	if item.Show.IDs.Trakt != 1388 {
		t.Errorf("expected Trakt ID 1388, got %d", item.Show.IDs.Trakt)
	}
}

func TestSearchService_FindByExternalID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/search_id_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/tvdb/267440" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	results, _, err := client.Search.FindByExternalID(context.Background(), "tvdb", "267440")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Show.IDs.Trakt != 46263 {
		t.Errorf("expected Trakt ID 46263, got %d", results[0].Show.IDs.Trakt)
	}
}
