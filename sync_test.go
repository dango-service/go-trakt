package trakt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSyncService_GetLastActivities(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_last_activities_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync/last_activities" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	activities, _, err := client.Sync.GetLastActivities(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if activities.All != "2026-11-06T20:00:00.000Z" {
		t.Errorf("expected all '2026-11-06T20:00:00.000Z', got %q", activities.All)
	}
	if activities.Episodes.WatchedAt != "2026-11-06T20:00:00.000Z" {
		t.Errorf("unexpected episodes watched_at: %s", activities.Episodes.WatchedAt)
	}
}

func TestSyncService_GetWatchlist(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_watchlist_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync/watchlist/shows" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	items, _, err := client.Sync.GetWatchlist(context.Background(), "shows")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Show.Title != "Steins;Gate" {
		t.Errorf("expected title 'Steins;Gate', got %q", items[0].Show.Title)
	}
}

func TestSyncService_GetHistory(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_history_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync/history/episodes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	items, _, err := client.Sync.GetHistory(context.Background(), "episodes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Episode.Title != "Turning Point" {
		t.Errorf("expected episode title 'Turning Point', got %q", items[0].Episode.Title)
	}
}

func TestSyncService_Mutations(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_mutation_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	lastMethod := ""
	lastPath := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastMethod = r.Method
		lastPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	input := &SyncItems{
		Movies: []*SyncMovie{
			{
				Title: "Inception",
				IDs:   IDs{Trakt: 16},
			},
		},
	}

	// AddHistory
	res, _, err := client.Sync.AddHistory(context.Background(), input)
	if err != nil || lastPath != "/sync/history" || lastMethod != http.MethodPost {
		t.Fatalf("unexpected AddHistory: %v, path: %s", err, lastPath)
	}
	if res.Added["movies"] != 1 {
		t.Errorf("expected 1 movie added, got %d", res.Added["movies"])
	}

	// RemoveHistory
	_, _, err = client.Sync.RemoveHistory(context.Background(), input)
	if err != nil || lastPath != "/sync/history/remove" {
		t.Fatalf("unexpected RemoveHistory: %v, path: %s", err, lastPath)
	}

	// AddWatchlist
	_, _, err = client.Sync.AddWatchlist(context.Background(), input)
	if err != nil || lastPath != "/sync/watchlist" {
		t.Fatalf("unexpected AddWatchlist: %v, path: %s", err, lastPath)
	}

	// RemoveWatchlist
	_, _, err = client.Sync.RemoveWatchlist(context.Background(), input)
	if err != nil || lastPath != "/sync/watchlist/remove" {
		t.Fatalf("unexpected RemoveWatchlist: %v, path: %s", err, lastPath)
	}

	// AddRatings
	_, _, err = client.Sync.AddRatings(context.Background(), input)
	if err != nil || lastPath != "/sync/ratings" {
		t.Fatalf("unexpected AddRatings: %v, path: %s", err, lastPath)
	}

	// RemoveRatings
	_, _, err = client.Sync.RemoveRatings(context.Background(), input)
	if err != nil || lastPath != "/sync/ratings/remove" {
		t.Fatalf("unexpected RemoveRatings: %v, path: %s", err, lastPath)
	}
}

func TestSyncService_NilInputs(t *testing.T) {
	client := NewClient("test-id")

	if _, _, err := client.Sync.AddHistory(context.Background(), nil); err == nil {
		t.Error("expected error for nil AddHistory")
	}
	if _, _, err := client.Sync.RemoveHistory(context.Background(), nil); err == nil {
		t.Error("expected error for nil RemoveHistory")
	}
	if _, _, err := client.Sync.AddWatchlist(context.Background(), nil); err == nil {
		t.Error("expected error for nil AddWatchlist")
	}
	if _, _, err := client.Sync.RemoveWatchlist(context.Background(), nil); err == nil {
		t.Error("expected error for nil RemoveWatchlist")
	}
}
