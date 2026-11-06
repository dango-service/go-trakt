# Library Scope and Boundaries  -  `go-trakt`

This document defines the functional perimeter of `go-trakt`.
The library does **not** attempt to cover the entirety of Trakt's API,
focusing strictly on the needs of **Dango** and media tracking tools.

---

## 1. Boundary Principle

The official Trakt API exposes hundreds of endpoints, including user comments, social networks, people/actors, certifications, countries, and custom recommendation engines.
Building an exhaustive client would:
- Multiply compiled binary size and build times.
- Increase runtime allocations with zero added value for scrobblers.
- Introduce an unmanageable maintenance burden.

Therefore, `go-trakt` follows a **strict minimum viable scope**, focusing on
searching, querying, and updating show, episode, and movie watch progress.

---

## 2. In-Scope Modules

1. **Search & Metadata**:
   - Title search across movies, shows, and episodes.
   - Cross-catalog external ID resolution (`/search/{id_type}/{id}`) supporting IMDb, TMDB, TVDB, and Trakt IDs.
   - Detailed show metadata, seasons, and episode lists.
   - Detailed movie metadata.
2. **Sync & Tracking**:
   - Activity checking (`/sync/last_activities`) for efficient cache validation.
   - Watchlist queries and mutations (`/sync/watchlist`).
   - Watch history queries and mutations (`/sync/history`).
   - Rating queries and mutations (`/sync/ratings`).
3. **Scrobble**:
   - Real-time media playback session scrobbling (`/scrobble/start`, `/scrobble/stop`, `/scrobble/pause`).
4. **User & Auth**:
   - User profile settings (`/users/settings`).
   - OAuth2 Authorization Code grant and token refresh (`/oauth/token`).
   - Device & PIN authorization flow (`/oauth/device/code`, `/oauth/device/token`).
   - Token revocation (`/oauth/revoke`).

---

## 3. Out-of-Scope Modules

- Comments, notes, and forum messaging.
- Cast, crew, and people biographies.
- Country, language, and network catalogs.
- Checkin service (superseded by `/scrobble`).
