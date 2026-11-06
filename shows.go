package trakt

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ShowsService handles show, season, and episode metadata.
type ShowsService struct {
	client *Client
}

// Show represents a TV show in Trakt.
type Show struct {
	Title         string   `json:"title"`
	Year          int      `json:"year,omitempty"`
	IDs           IDs      `json:"ids"`
	Overview      string   `json:"overview,omitempty"`
	FirstAired    string   `json:"first_aired,omitempty"`
	Status        string   `json:"status,omitempty"`
	Rating        float64  `json:"rating,omitempty"`
	Votes         int      `json:"votes,omitempty"`
	AiredEpisodes int      `json:"aired_episodes,omitempty"`
	Genres        []string `json:"genres,omitempty"`
}

// Season represents a season of a TV show with optional episode list.
type Season struct {
	Number        int        `json:"number"`
	IDs           IDs        `json:"ids"`
	Rating        float64    `json:"rating,omitempty"`
	EpisodeCount  int        `json:"episode_count,omitempty"`
	AiredEpisodes int        `json:"aired_episodes,omitempty"`
	Overview      string     `json:"overview,omitempty"`
	Episodes      []*Episode `json:"episodes,omitempty"`
}

// Episode represents an individual episode of a TV show.
type Episode struct {
	Season     int     `json:"season"`
	Number     int     `json:"number"`
	Title      string  `json:"title,omitempty"`
	IDs        IDs     `json:"ids"`
	Overview   string  `json:"overview,omitempty"`
	Rating     float64 `json:"rating,omitempty"`
	FirstAired string  `json:"first_aired,omitempty"`
}

// GetByID fetches full metadata for a show by Trakt ID, slug, or external ID.
func (s *ShowsService) GetByID(ctx context.Context, id string) (*Show, *Response, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, nil, fmt.Errorf("trakt: show ID cannot be empty")
	}

	path := fmt.Sprintf("/shows/%s?extended=full", url.PathEscape(id))

	var show Show
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &show)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get show %q: %w", id, err)
	}
	return &show, resp, nil
}

// GetSeasons fetches the seasons and episode list for a show.
func (s *ShowsService) GetSeasons(ctx context.Context, id string) ([]*Season, *Response, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, nil, fmt.Errorf("trakt: show ID cannot be empty")
	}

	path := fmt.Sprintf("/shows/%s/seasons?extended=episodes", url.PathEscape(id))

	var seasons []*Season
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &seasons)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get seasons for show %q: %w", id, err)
	}
	return seasons, resp, nil
}
