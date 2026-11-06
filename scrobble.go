package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ScrobbleService handles real-time playback scrobbling sessions.
type ScrobbleService struct {
	client *Client
}

// ScrobbleInput represents a payload for sending playback events to Trakt.
type ScrobbleInput struct {
	Movie      *SyncMovie   `json:"movie,omitempty"`
	Show       *SyncShow    `json:"show,omitempty"`
	Episode    *SyncEpisode `json:"episode,omitempty"`
	Progress   float64      `json:"progress"` // percentage (0.0 to 100.0)
	AppVersion string       `json:"app_version,omitempty"`
	AppDate    string       `json:"app_date,omitempty"`
}

// ScrobbleResponse represents the confirmation from a scrobble action.
type ScrobbleResponse struct {
	ID       int64    `json:"id"`
	Action   string   `json:"action"` // "start", "pause", "scrobble"
	Progress float64  `json:"progress"`
	Movie    *Movie   `json:"movie,omitempty"`
	Show     *Show    `json:"show,omitempty"`
	Episode  *Episode `json:"episode,omitempty"`
}

// Start notifies Trakt that playback has started.
func (s *ScrobbleService) Start(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	return s.doScrobble(ctx, "/scrobble/start", input)
}

// Pause notifies Trakt that playback was paused.
func (s *ScrobbleService) Pause(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	return s.doScrobble(ctx, "/scrobble/pause", input)
}

// Stop notifies Trakt that playback stopped or finished.
func (s *ScrobbleService) Stop(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	return s.doScrobble(ctx, "/scrobble/stop", input)
}

func (s *ScrobbleService) doScrobble(ctx context.Context, path string, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("trakt: scrobble: input is required")
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("trakt: marshal scrobble input: %w", err)
	}

	var res ScrobbleResponse
	resp, err := s.client.doRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), &res)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: scrobble %s: %w", path, err)
	}
	return &res, resp, nil
}
