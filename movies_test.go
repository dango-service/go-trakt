package trakt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMoviesService_GetByID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/movie_by_id_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/movies/inception-2010" {
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

	movie, _, err := client.Movies.GetByID(context.Background(), "inception-2010")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movie.Title != "Inception" {
		t.Errorf("expected title 'Inception', got %q", movie.Title)
	}
	if movie.Runtime != 148 {
		t.Errorf("expected runtime 148, got %d", movie.Runtime)
	}
	if movie.IDs.Trakt != 16 {
		t.Errorf("expected Trakt ID 16, got %d", movie.IDs.Trakt)
	}
}

func TestMoviesService_GetByID_EmptyID(t *testing.T) {
	client := NewClient("test-id")

	if _, _, err := client.Movies.GetByID(context.Background(), "   "); err == nil {
		t.Error("expected error for empty movie ID")
	}
	if _, _, err := client.Shows.GetByID(context.Background(), "   "); err == nil {
		t.Error("expected error for empty show ID")
	}
}
