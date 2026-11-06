# Architectural Invariants  -  go-trakt

Structural rules and invariants that must be upheld at all times in `go-trakt`.

---

## 1. Independence and Decoupling

- **Completely Agnostic of Dango**: `go-trakt` is an autonomous library.
  Under no circumstances may it import `Dango` packages, reference Plex, SQLite,
  media matching algorithms, or webhook handlers.
- **Zero Global State & Thread Safety**: Multiple `Client` instances must be able to run concurrently
  and safely within the same process without interference or shared mutable state.

---

## 2. Layer Separation: Transport vs Typed Public API

1. **Internal Transport Layer (`doRequest`)**:
   - Handles HTTP communication (`GET`, `POST`, `DELETE`).
   - Injects required headers (`Content-Type: application/json`, `trakt-api-version: 2`, `trakt-api-key: <client_id>`, `User-Agent`).
   - Injects `Authorization: Bearer <token>` when tokens are present.
   - Parses Trakt rate limit headers (`X-Ratelimit`, `Retry-After`).
   - Deserializes JSON via streaming `json.NewDecoder` bounded by `io.LimitReader(body, maxBytes)`.
   - Normalizes HTTP status codes and error payloads into Go error types.
2. **Typed Public Services**:
   - `SearchService`: Full-text search and cross-catalog external ID resolution (IMDb, TMDB, TVDB, Trakt ID).
   - `ShowsService`: TV show details, seasons, and episode lists.
   - `MoviesService`: Movie details.
   - `SyncService`: Watchlist, watch history, ratings, and activity polling check.
   - `ScrobbleService`: Real-time media playback scrobbling sessions.
   - `UsersService`: User account settings and profiles.
   - `AuthService`: OAuth2 Authorization Code exchange, token refresh, token revoke, and Device/PIN flow.

---

## 3. Rate Limit Handling

- Trakt enforces quota limits (typically 1,000 calls / 5 minutes) and returns HTTP 429 when throttled.
- The transport layer extracts `Retry-After` (supporting integer seconds and HTTP-Date) and `X-Ratelimit` into `RateLimitInfo`.
- The library does not mutate client state or sleep invisibly; it returns enriched `ErrRateLimited` with duration details, allowing callers to handle backoff deterministically.

---

## 4. Credential Protection

- Access tokens and refresh tokens must never be written to logs in plaintext.
- The `Token` struct implements `fmt.Stringer` and `fmt.GoStringer` returning `[REDACTED]`.
- Tokens are stripped of leading and trailing whitespace upon option injection.
