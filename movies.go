package trakt

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MoviesService handles movie metadata.
type MoviesService struct {
	client *Client
}

// Movie represents a movie in Trakt.
type Movie struct {
	Title    string   `json:"title"`
	Year     int      `json:"year,omitempty"`
	IDs      IDs      `json:"ids"`
	Overview string   `json:"overview,omitempty"`
	Released string   `json:"released,omitempty"`
	Runtime  int      `json:"runtime,omitempty"`
	Rating   float64  `json:"rating,omitempty"`
	Votes    int      `json:"votes,omitempty"`
	Genres   []string `json:"genres,omitempty"`
}

// GetByID fetches full metadata for a movie by Trakt ID, slug, or external ID.
func (s *MoviesService) GetByID(ctx context.Context, id string) (*Movie, *Response, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, nil, fmt.Errorf("trakt: movie ID cannot be empty")
	}

	path := fmt.Sprintf("/movies/%s?extended=full", url.PathEscape(id))

	var movie Movie
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &movie)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get movie %q: %w", id, err)
	}
	return &movie, resp, nil
}
