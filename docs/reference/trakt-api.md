# Trakt.tv API v2 Reference  -  `go-trakt`

Technical reference for interacting with the Trakt.tv API v2 (`https://api.trakt.tv`).

---

## 1. Endpoints Overview

| Category | HTTP Method | Path | Description |
| :--- | :--- | :--- | :--- |
| **Search** | `GET` | `/search/{type}` | Search items by title query |
| **Search** | `GET` | `/search/{id_type}/{id}` | Find items by external ID (IMDb, TMDB, TVDB, Trakt) |
| **Shows** | `GET` | `/shows/{id}` | Detailed TV show metadata |
| **Shows** | `GET` | `/shows/{id}/seasons` | Seasons and episodes for a TV show |
| **Movies** | `GET` | `/movies/{id}` | Detailed movie metadata |
| **Sync** | `GET` | `/sync/last_activities` | Last activity timestamps (crucial for polling) |
| **Sync** | `GET` | `/sync/watchlist/{type}` | Get user's watchlist entries |
| **Sync** | `GET` | `/sync/history/{type}` | Get user's watch history |
| **Sync** | `GET` | `/sync/ratings/{type}` | Get user's ratings |
| **Sync** | `POST` | `/sync/history` | Add items or episodes to watch history |
| **Sync** | `POST` | `/sync/history/remove` | Remove items from watch history |
| **Sync** | `POST` | `/sync/watchlist` | Add items to watchlist |
| **Sync** | `POST` | `/sync/watchlist/remove` | Remove items from watchlist |
| **Sync** | `POST` | `/sync/ratings` | Add ratings |
| **Sync** | `POST` | `/sync/ratings/remove` | Remove ratings |
| **Scrobble** | `POST` | `/scrobble/start` | Start scrobbling playback session |
| **Scrobble** | `POST` | `/scrobble/stop` | Stop scrobbling playback session |
| **Scrobble** | `POST` | `/scrobble/pause` | Pause scrobbling playback session |
| **Users** | `GET` | `/users/settings` | Authenticated user settings and profile |
| **OAuth2** | `GET` | `/oauth/authorize` | Authorization code consent URL |
| **OAuth2** | `POST` | `/oauth/token` | Exchange code or refresh token |
| **OAuth2** | `POST` | `/oauth/device/code` | Request device code & user PIN |
| **OAuth2** | `POST` | `/oauth/device/token` | Poll device authorization status |
| **OAuth2** | `POST` | `/oauth/revoke` | Revoke active access token |

---

## 2. Required Request Headers

Every request to `https://api.trakt.tv` must provide:
- `Content-Type`: `application/json`
- `trakt-api-version`: `2`
- `trakt-api-key`: Application Client ID
- `User-Agent`: Descriptive client identifier (`go-trakt/0.1.0`)
- `Authorization`: `Bearer <token>` on user-authenticated endpoints

---

## 3. Rate Limits & Headers

- **Default Limits**: 1,000 calls / 5 minutes for GET requests.
- **Throttling Status**: Returns HTTP `429 Too Many Requests`.
- **Headers**:
  - `Retry-After`: Integer seconds to pause.
  - `X-Ratelimit`: JSON object containing period, limit, remaining calls, and reset timestamp.
