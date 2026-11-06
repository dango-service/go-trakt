# go-trakt

[![Go Reference](https://pkg.go.dev/badge/github.com/dango-service/go-trakt.svg)](https://pkg.go.dev/github.com/dango-service/go-trakt)
[![CI](https://github.com/dango-service/go-trakt/actions/workflows/ci.yml/badge.svg)](https://github.com/dango-service/go-trakt/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A lightweight, standalone Go client library for the [Trakt.tv](https://trakt.tv) API v2 (`https://api.trakt.tv`).
Designed with **zero external dependencies** (Go standard library only), direct streaming deserialization bounded by memory limits, and strict efficiency for background sync daemons like [Dango](https://github.com/dango-service/dango).

---

## Features

- **Zero External Dependencies**: Built entirely upon the Go standard library (`net/http`, `encoding/json`, `context`, `time`, `net/url`).
- **Streaming JSON**: Direct streaming decoding with `json.NewDecoder` bounded by `io.LimitReader` (10 MiB ceiling) to prevent unbounded memory growth and memory exhaustion attacks.
- **Robust Rate Limit Handling**: Automatic parsing of Trakt quota headers (`Retry-After` in seconds or RFC1123 HTTP-Date, plus `X-Ratelimit` JSON header) with sentinel `ErrRateLimited` error wrapping.
- **Credential Protection**: Redacted string formatting (`[REDACTED]`) on `Token` to prevent secret leakage in logs, traces, or debug outputs.
- **RFC 8628 Device Authorization**: Built-in support for Device Code / PIN flow (`/oauth/device/code`, `/oauth/device/token`), ideal for headless appliances, Docker containers, and NAS systems without a local web server.
- **Modular Services**:
  - `Search`: Full-text search across shows, movies, and episodes, as well as external ID lookup (IMDb, TMDB, TVDB, Trakt).
  - `Shows`: TV show metadata (`/shows/{id}`) and season/episode trees (`/shows/{id}/seasons`).
  - `Movies`: Movie metadata (`/movies/{id}`).
  - `Sync`: Activity polling checks (`/sync/last_activities`), watchlist retrieval/mutations, playback history, and ratings.
  - `Scrobble`: Real-time media playback sessions (`/scrobble/start`, `/scrobble/pause`, `/scrobble/stop`).
  - `Users`: User settings and account profile retrieval (`/users/settings`).
  - `Auth`: Full OAuth2 Authorization Code exchange, token refresh, token revocation, and RFC 8628 Device Code polling.
- **100% Deterministic Testing**: Complete test suite using `net/http/httptest` with zero real-network requests.

---

## Installation

```bash
go get github.com/dango-service/go-trakt
```

---

## Usage Examples

### 1. Client Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/dango-service/go-trakt"
)

func main() {
    client := trakt.NewClient("YOUR_TRAKT_CLIENT_ID",
        trakt.WithToken("USER_ACCESS_TOKEN"),
    )

    // Check last sync activities
    activities, _, err := client.Sync.GetLastActivities(context.Background())
    if err != nil {
        log.Fatalf("failed to fetch activities: %v", err)
    }

    fmt.Printf("Movies Watched At: %s\n", activities.Movies.WatchedAt)
    fmt.Printf("Episodes Watched At: %s\n", activities.Episodes.WatchedAt)
}
```

### 2. Search & External ID Resolution

```go
// Search shows by query text
results, _, err := client.Search.Search(ctx, "breaking bad", trakt.WithSearchType("show"), trakt.WithSearchLimit(5))

// Find media by external ID (e.g. IMDb ID tt0903747)
matches, _, err := client.Search.FindByExternalID(ctx, "imdb", "tt0903747", trakt.WithSearchType("show"))
if len(matches) > 0 && matches[0].Show != nil {
    fmt.Printf("Trakt ID: %d, Title: %s\n", matches[0].Show.IDs.Trakt, matches[0].Show.Title)
}
```

### 3. Shows & Movies Metadata

```go
// Get show details with extended information
show, _, err := client.Shows.GetByID(ctx, "breaking-bad", trakt.WithExtended("full"))

// Get show seasons with episodes
seasons, _, err := client.Shows.GetSeasons(ctx, "breaking-bad", trakt.WithExtended("episodes"))

// Get movie details
movie, _, err := client.Movies.GetByID(ctx, "tron-legacy-2010")
```

### 4. Watchlist & History Sync

```go
// Retrieve user watchlist
watchlist, _, err := client.Sync.GetWatchlist(ctx, trakt.WithSyncType("shows"))

// Add watched episode to history
result, _, err := client.Sync.AddHistory(ctx, &trakt.SyncItems{
    Shows: []trakt.SyncShow{
        {
            IDs: trakt.IDs{Trakt: 1388},
            Seasons: []trakt.SyncSeason{
                {
                    Number: 1,
                    Episodes: []trakt.SyncEpisode{
                        {Number: 1},
                    },
                },
            },
        },
    },
})
```

### 5. Real-Time Scrobbling

```go
// Start watching playback session
_, _, err := client.Scrobble.Start(ctx, &trakt.ScrobbleBody{
    Movie: &trakt.Movie{
        Title: "TRON: Legacy",
        Year:  2010,
        IDs:   trakt.IDs{Trakt: 343},
    },
    Progress: 25.0,
})

// Stop / complete playback session
_, _, err = client.Scrobble.Stop(ctx, &trakt.ScrobbleBody{
    Movie: &trakt.Movie{
        IDs: trakt.IDs{Trakt: 343},
    },
    Progress: 100.0,
})
```

### 6. Headless Authentication (RFC 8628 Device Code / PIN)

```go
// 1. Generate device authorization codes
codeResp, _, err := client.Auth.GenerateDeviceCode(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Go to %s and enter code: %s\n", codeResp.VerificationURL, codeResp.UserCode)

// 2. Poll until user approves
token, _, err := client.Auth.PollDeviceToken(ctx, codeResp.DeviceCode, clientSecret)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Access token successfully acquired!")
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
