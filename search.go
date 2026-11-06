package trakt

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SearchService handles searching catalog metadata and cross-catalog external ID resolution.
type SearchService struct {
	client *Client
}

// IDs contains cross-catalog identifiers for Trakt items.
type IDs struct {
	Trakt int    `json:"trakt,omitempty"`
	Slug  string `json:"slug,omitempty"`
	TVDB  int    `json:"tvdb,omitempty"`
	IMDb  string `json:"imdb,omitempty"`
	TMDB  int    `json:"tmdb,omitempty"`
}

// SearchResult represents an item returned by Trakt search endpoints.
type SearchResult struct {
	Type    string   `json:"type"` // "movie", "show", "episode", etc.
	Score   float64  `json:"score,omitempty"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
}

// Search queries items by text. mediaType must be "movie", "show", "episode", or comma-separated list.
func (s *SearchService) Search(ctx context.Context, mediaType, query string, opts ...SearchOption) ([]*SearchResult, *Response, error) {
	if mediaType == "" {
		mediaType = "movie,show"
	}
	params := url.Values{}
	params.Set("query", query)

	for _, opt := range opts {
		if opt != nil {
			opt(params)
		}
	}

	path := fmt.Sprintf("/search/%s?%s", url.PathEscape(mediaType), params.Encode())

	var results []*SearchResult
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &results)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: search %s %q: %w", mediaType, query, err)
	}
	return results, resp, nil
}

// FindByExternalID queries items by external identifier.
// idType must be "trakt", "imdb", "tmdb", or "tvdb".
func (s *SearchService) FindByExternalID(ctx context.Context, idType, id string, opts ...SearchOption) ([]*SearchResult, *Response, error) {
	idType = strings.ToLower(strings.TrimSpace(idType))
	id = strings.TrimSpace(id)

	params := url.Values{}
	for _, opt := range opts {
		if opt != nil {
			opt(params)
		}
	}

	path := fmt.Sprintf("/search/%s/%s", url.PathEscape(idType), url.PathEscape(id))
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	var results []*SearchResult
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &results)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: find by %s ID %q: %w", idType, id, err)
	}
	return results, resp, nil
}
