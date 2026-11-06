package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SyncService handles watchlist queries, watch history, user ratings, and activity polling.
type SyncService struct {
	client *Client
}

// LastActivities contains timestamps of the user's most recent activities across all media types.
type LastActivities struct {
	All      string                    `json:"all"`
	Movies   MovieActivityTimestamps   `json:"movies"`
	Episodes EpisodeActivityTimestamps `json:"episodes"`
	Shows    ShowActivityTimestamps    `json:"shows"`
	Seasons  SeasonActivityTimestamps  `json:"seasons"`
}

// MovieActivityTimestamps holds activity timestamps for movies.
type MovieActivityTimestamps struct {
	WatchedAt     string `json:"watched_at"`
	CollectedAt   string `json:"collected_at"`
	RatedAt       string `json:"rated_at"`
	WatchlistedAt string `json:"watchlisted_at"`
	PausedAt      string `json:"paused_at"`
}

// EpisodeActivityTimestamps holds activity timestamps for episodes.
type EpisodeActivityTimestamps struct {
	WatchedAt     string `json:"watched_at"`
	CollectedAt   string `json:"collected_at"`
	RatedAt       string `json:"rated_at"`
	WatchlistedAt string `json:"watchlisted_at"`
	PausedAt      string `json:"paused_at"`
}

// ShowActivityTimestamps holds activity timestamps for shows.
type ShowActivityTimestamps struct {
	RatedAt       string `json:"rated_at"`
	WatchlistedAt string `json:"watchlisted_at"`
}

// SeasonActivityTimestamps holds activity timestamps for seasons.
type SeasonActivityTimestamps struct {
	RatedAt       string `json:"rated_at"`
	WatchlistedAt string `json:"watchlisted_at"`
}

// WatchlistItem represents an item on a user's watchlist.
type WatchlistItem struct {
	Rank     int      `json:"rank"`
	ListedAt string   `json:"listed_at"`
	Type     string   `json:"type"`
	Movie    *Movie   `json:"movie,omitempty"`
	Show     *Show    `json:"show,omitempty"`
	Episode  *Episode `json:"episode,omitempty"`
}

// HistoryItem represents an entry in a user's watch history.
type HistoryItem struct {
	ID        int64    `json:"id"`
	WatchedAt string   `json:"watched_at"`
	Action    string   `json:"action"` // "watch" or "scrobble"
	Type      string   `json:"type"`
	Movie     *Movie   `json:"movie,omitempty"`
	Show      *Show    `json:"show,omitempty"`
	Episode   *Episode `json:"episode,omitempty"`
}

// RatingItem represents a user's rating for a media item.
type RatingItem struct {
	RatedAt string   `json:"rated_at"`
	Rating  int      `json:"rating"` // 1-10
	Type    string   `json:"type"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
}

// SyncItems contains items to be added or removed in a sync mutation.
type SyncItems struct {
	Movies   []*SyncMovie   `json:"movies,omitempty"`
	Shows    []*SyncShow    `json:"shows,omitempty"`
	Episodes []*SyncEpisode `json:"episodes,omitempty"`
}

// SyncMovie represents a movie in a sync mutation.
type SyncMovie struct {
	Title     string `json:"title,omitempty"`
	Year      int    `json:"year,omitempty"`
	IDs       IDs    `json:"ids"`
	WatchedAt string `json:"watched_at,omitempty"`
	Rating    int    `json:"rating,omitempty"`
}

// SyncShow represents a show in a sync mutation.
type SyncShow struct {
	Title   string        `json:"title,omitempty"`
	Year    int           `json:"year,omitempty"`
	IDs     IDs           `json:"ids"`
	Seasons []*SyncSeason `json:"seasons,omitempty"`
	Rating  int           `json:"rating,omitempty"`
}

// SyncSeason represents a season in a sync mutation.
type SyncSeason struct {
	Number   int            `json:"number"`
	Episodes []*SyncEpisode `json:"episodes,omitempty"`
	Rating   int            `json:"rating,omitempty"`
}

// SyncEpisode represents an episode in a sync mutation.
type SyncEpisode struct {
	Number    int    `json:"number,omitempty"`
	Season    int    `json:"season,omitempty"`
	IDs       IDs    `json:"ids,omitempty"`
	WatchedAt string `json:"watched_at,omitempty"`
	Rating    int    `json:"rating,omitempty"`
}

// SyncResult represents the summary response from sync mutation endpoints.
type SyncResult struct {
	Added    map[string]int `json:"added,omitempty"`
	Deleted  map[string]int `json:"deleted,omitempty"`
	NotFound map[string]any `json:"not_found,omitempty"`
}

// GetLastActivities retrieves timestamps of the user's latest activities for cache optimization.
func (s *SyncService) GetLastActivities(ctx context.Context) (*LastActivities, *Response, error) {
	var activities LastActivities
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/sync/last_activities", nil, &activities)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get last activities: %w", err)
	}
	return &activities, resp, nil
}

// GetWatchlist retrieves the user's watchlist entries for a mediaType ("movies", "shows", "seasons", "episodes").
func (s *SyncService) GetWatchlist(ctx context.Context, mediaType string, opts ...SyncOption) ([]*WatchlistItem, *Response, error) {
	path := "/sync/watchlist"
	if mediaType != "" {
		path = fmt.Sprintf("/sync/watchlist/%s", url.PathEscape(mediaType))
	}
	path = applySyncOptions(path, opts...)

	var items []*WatchlistItem
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &items)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get watchlist %s: %w", mediaType, err)
	}
	return items, resp, nil
}

// GetHistory retrieves the user's watch history for a mediaType ("movies", "shows", "seasons", "episodes").
func (s *SyncService) GetHistory(ctx context.Context, mediaType string, opts ...SyncOption) ([]*HistoryItem, *Response, error) {
	path := "/sync/history"
	if mediaType != "" {
		path = fmt.Sprintf("/sync/history/%s", url.PathEscape(mediaType))
	}
	path = applySyncOptions(path, opts...)

	var items []*HistoryItem
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &items)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get history %s: %w", mediaType, err)
	}
	return items, resp, nil
}

// GetRatings retrieves the user's ratings for a mediaType ("movies", "shows", "seasons", "episodes").
func (s *SyncService) GetRatings(ctx context.Context, mediaType string, opts ...SyncOption) ([]*RatingItem, *Response, error) {
	path := "/sync/ratings"
	if mediaType != "" {
		path = fmt.Sprintf("/sync/ratings/%s", url.PathEscape(mediaType))
	}
	path = applySyncOptions(path, opts...)

	var items []*RatingItem
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &items)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get ratings %s: %w", mediaType, err)
	}
	return items, resp, nil
}

// AddHistory adds movies, shows, or episodes to the user's watch history.
func (s *SyncService) AddHistory(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error) {
	return s.doSyncMutation(ctx, "/sync/history", input)
}

// RemoveHistory removes movies, shows, or episodes from the user's watch history.
func (s *SyncService) RemoveHistory(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error) {
	return s.doSyncMutation(ctx, "/sync/history/remove", input)
}

// AddWatchlist adds movies, shows, or episodes to the user's watchlist.
func (s *SyncService) AddWatchlist(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error) {
	return s.doSyncMutation(ctx, "/sync/watchlist", input)
}

// RemoveWatchlist removes movies, shows, or episodes from the user's watchlist.
func (s *SyncService) RemoveWatchlist(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error) {
	return s.doSyncMutation(ctx, "/sync/watchlist/remove", input)
}

// AddRatings adds or updates ratings for movies, shows, or episodes.
func (s *SyncService) AddRatings(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error) {
	return s.doSyncMutation(ctx, "/sync/ratings", input)
}

// RemoveRatings removes ratings from movies, shows, or episodes.
func (s *SyncService) RemoveRatings(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error) {
	return s.doSyncMutation(ctx, "/sync/ratings/remove", input)
}

func (s *SyncService) doSyncMutation(ctx context.Context, path string, input *SyncItems) (*SyncResult, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("trakt: sync mutation: input is required")
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("trakt: marshal sync input: %w", err)
	}

	var result SyncResult
	resp, err := s.client.doRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), &result)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: sync mutation %s: %w", path, err)
	}
	return &result, resp, nil
}

func applySyncOptions(basePath string, opts ...SyncOption) string {
	params := url.Values{}
	params.Set("extended", "full")
	for _, opt := range opts {
		if opt != nil {
			opt(params)
		}
	}
	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", basePath, params.Encode())
	}
	return basePath
}
