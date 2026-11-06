# Testing Strategy  -  `go-trakt`

This document defines testing standards for `go-trakt`, ensuring deterministic,
fast, isolated, and reproducible test suites without external network access.

---

## 1. Core Testing Principles

1. **Zero Real Network Traffic**: No unit test may connect to `https://api.trakt.tv`. All network traffic is mocked locally via `httptest.Server`.
2. **Determinism and Speed**: The full suite must complete in less than 2 seconds.
3. **Table-Driven Verification**: Tests verify both what the client sends (URL query parameters, HTTP headers, request bodies) and what it deserializes (domain models, rate limit info, errors).
4. **Error Coverage**: Full testing of HTTP 401, 404, 409, and 429 (`Retry-After` seconds and HTTP-Date format).

---

## 2. Local Test Server (`httptest.Server`)

Example pattern:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("trakt-api-key") != "test-client-id" {
        http.Error(w, `{"message":"invalid api key"}`, http.StatusUnauthorized)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`[{"title":"Breaking Bad"}]`))
}))
defer server.Close()

client := NewClient("test-client-id", WithBaseURL(server.URL))
```
