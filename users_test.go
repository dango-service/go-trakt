package trakt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUsersService_GetSettings(t *testing.T) {
	fixture, err := os.ReadFile("testdata/user_settings_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL), WithToken("test-token"))

	settings, _, err := client.Users.GetSettings(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings.User.Username != "dango_watcher" {
		t.Errorf("expected username 'dango_watcher', got %q", settings.User.Username)
	}
	if !settings.User.VIP {
		t.Error("expected VIP to be true")
	}
	if settings.Account.Timezone != "America/New_York" {
		t.Errorf("expected timezone 'America/New_York', got %q", settings.Account.Timezone)
	}
}
