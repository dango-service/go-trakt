package trakt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScrobbleService_StartPauseStop(t *testing.T) {
	calledPath := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":12345,"action":"scrobble","progress":55.0}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	input := &ScrobbleInput{
		Movie: &SyncMovie{
			Title: "Inception",
			IDs:   IDs{Trakt: 16},
		},
		Progress: 55.0,
	}

	// Start
	res, _, err := client.Scrobble.Start(context.Background(), input)
	if err != nil || calledPath != "/scrobble/start" {
		t.Fatalf("unexpected Start error: %v, path: %s", err, calledPath)
	}
	if res.Progress != 55.0 {
		t.Errorf("expected progress 55.0, got %f", res.Progress)
	}

	// Pause
	_, _, err = client.Scrobble.Pause(context.Background(), input)
	if err != nil || calledPath != "/scrobble/pause" {
		t.Fatalf("unexpected Pause error: %v, path: %s", err, calledPath)
	}

	// Stop
	_, _, err = client.Scrobble.Stop(context.Background(), input)
	if err != nil || calledPath != "/scrobble/stop" {
		t.Fatalf("unexpected Stop error: %v, path: %s", err, calledPath)
	}

	// Nil input
	if _, _, err := client.Scrobble.Start(context.Background(), nil); err == nil {
		t.Error("expected error for nil input")
	}
}
