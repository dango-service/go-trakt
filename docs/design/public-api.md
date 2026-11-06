# Go Public API Specification  -  `go-trakt`

This document defines the Go public API signatures, models, options, and contracts
for `go-trakt`. The design mirrors the clean architecture of `go-kitsu`, `go-anilist`, and `go-simkl`.

---

## 1. Package and Core Client

```go
package trakt

import (
    "context"
    "net/http"
    "net/url"
)

// Client coordinates communication with the Trakt.tv API.
type Client struct {
    Search   *SearchService
    Shows    *ShowsService
    Movies   *MoviesService
    Sync     *SyncService
    Scrobble *ScrobbleService
    Users    *UsersService
    Auth     *AuthService
}

func NewClient(clientID string, opts ...ClientOption) *Client
```

---

## 2. Functional Options

```go
type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption
func WithBaseURL(rawURL string) ClientOption
func WithToken(token string) ClientOption
func WithUserAgent(userAgent string) ClientOption
func WithMaxBodyBytes(limit int64) ClientOption
```

---

## 3. Core Models

```go
type IDs struct {
    Trakt int    `json:"trakt,omitempty"`
    Slug  string `json:"slug,omitempty"`
    TVDB  int    `json:"tvdb,omitempty"`
    IMDb  string `json:"imdb,omitempty"`
    TMDB  int    `json:"tmdb,omitempty"`
}

type Show struct {
    Title    string   `json:"title"`
    Year     int      `json:"year,omitempty"`
    IDs      IDs      `json:"ids"`
    Overview string   `json:"overview,omitempty"`
    Status   string   `json:"status,omitempty"`
    AiredEpisodes int `json:"aired_episodes,omitempty"`
}

type Movie struct {
    Title    string `json:"title"`
    Year     int    `json:"year,omitempty"`
    IDs      IDs    `json:"ids"`
    Overview string `json:"overview,omitempty"`
    Runtime  int    `json:"runtime,omitempty"`
}

type Episode struct {
    Season int    `json:"season"`
    Number int    `json:"number"`
    Title  string `json:"title,omitempty"`
    IDs    IDs    `json:"ids"`
}
```

---

## 4. Service Signatures

### SearchService
```go
func (s *SearchService) Search(ctx context.Context, mediaType, query string, opts ...SearchOption) ([]*SearchResult, *Response, error)
func (s *SearchService) FindByExternalID(ctx context.Context, idType, id string, opts ...SearchOption) ([]*SearchResult, *Response, error)
```

### ShowsService & MoviesService
```go
func (s *ShowsService) GetByID(ctx context.Context, id string) (*Show, *Response, error)
func (s *ShowsService) GetSeasons(ctx context.Context, id string) ([]*Season, *Response, error)
func (s *MoviesService) GetByID(ctx context.Context, id string) (*Movie, *Response, error)
```

### SyncService
```go
func (s *SyncService) GetLastActivities(ctx context.Context) (*LastActivities, *Response, error)
func (s *SyncService) GetWatchlist(ctx context.Context, mediaType string, opts ...SyncOption) ([]*WatchlistItem, *Response, error)
func (s *SyncService) GetHistory(ctx context.Context, mediaType string, opts ...SyncOption) ([]*HistoryItem, *Response, error)
func (s *SyncService) AddHistory(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error)
func (s *SyncService) RemoveHistory(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error)
func (s *SyncService) AddWatchlist(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error)
func (s *SyncService) RemoveWatchlist(ctx context.Context, input *SyncItems) (*SyncResult, *Response, error)
```

### ScrobbleService
```go
func (s *ScrobbleService) Start(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error)
func (s *ScrobbleService) Stop(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error)
func (s *ScrobbleService) Pause(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error)
```

### UsersService
```go
func (s *UsersService) GetSettings(ctx context.Context) (*UserSettings, *Response, error)
```

### AuthService
```go
func (s *AuthService) AuthorizeURL(redirectURI, state string) string
func (s *AuthService) ExchangeCode(ctx context.Context, code, redirectURI, clientSecret string) (*Token, *Response, error)
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, clientSecret, redirectURI string) (*Token, *Response, error)
func (s *AuthService) GenerateDeviceCode(ctx context.Context) (*DeviceCodeResponse, *Response, error)
func (s *AuthService) PollDeviceToken(ctx context.Context, deviceCode, clientSecret string) (*Token, *Response, error)
func (s *AuthService) RevokeToken(ctx context.Context, token, clientSecret string) (*Response, error)
```
