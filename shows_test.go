package trakt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestShowsService_GetByID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/show_by_id_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/steins-gate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("extended") != "full" {
			t.Errorf("expected extended=full, got %s", r.URL.Query().Get("extended"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	show, _, err := client.Shows.GetByID(context.Background(), "steins-gate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if show.Title != "Steins;Gate" {
		t.Errorf("expected title 'Steins;Gate', got %q", show.Title)
	}
	if show.AiredEpisodes != 24 {
		t.Errorf("expected 24 episodes, got %d", show.AiredEpisodes)
	}
	if show.IDs.Trakt != 38165 {
		t.Errorf("expected Trakt ID 38165, got %d", show.IDs.Trakt)
	}
}

func TestShowsService_GetSeasons(t *testing.T) {
	fixture, err := os.ReadFile("testdata/show_seasons_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/38165/seasons" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	seasons, _, err := client.Shows.GetSeasons(context.Background(), "38165")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(seasons) != 1 {
		t.Fatalf("expected 1 season, got %d", len(seasons))
	}
	season := seasons[0]
	if len(season.Episodes) != 2 {
		t.Fatalf("expected 2 episodes, got %d", len(season.Episodes))
	}
	if season.Episodes[0].Title != "Turning Point" {
		t.Errorf("expected episode 1 title 'Turning Point', got %q", season.Episodes[0].Title)
	}
}
