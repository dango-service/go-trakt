package trakt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// doRequest performs an HTTP request against the Trakt API, enforces memory limits,
// and decodes JSON responses or structured errors.
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, v any) (*Response, error) {
	relURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("trakt: invalid path %q: %w", path, err)
	}

	u := c.baseURL.ResolveReference(relURL)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("trakt: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-version", apiVersion)
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.clientID != "" {
		req.Header.Set("trakt-api-key", c.clientID)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("trakt: request cancelled: %w", ctx.Err())
		default:
			return nil, fmt.Errorf("trakt: request failed: %w", err)
		}
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	limitedReader := io.LimitReader(resp.Body, c.maxBytes)
	rateLimit := parseRateLimit(resp.Header)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		response := &Response{
			Response:  resp,
			RateLimit: rateLimit,
		}

		if v != nil && resp.StatusCode != http.StatusNoContent {
			if err := json.NewDecoder(limitedReader).Decode(v); err != nil {
				return response, fmt.Errorf("trakt: failed to decode response: %w", err)
			}
		}
		return response, nil
	}

	// Handle error payloads
	var errResp ErrorResponse
	_ = json.NewDecoder(limitedReader).Decode(&errResp)
	errResp.Response = resp
	errResp.RateLimit = rateLimit

	response := &Response{
		Response:  resp,
		RateLimit: rateLimit,
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return response, fmt.Errorf("%w: %s", ErrRateLimited, errResp.Error())
	case http.StatusUnauthorized:
		return response, fmt.Errorf("%w: %s", ErrUnauthorized, errResp.Error())
	case http.StatusNotFound:
		return response, fmt.Errorf("%w: %s", ErrNotFound, errResp.Error())
	case http.StatusConflict:
		return response, fmt.Errorf("%w: %s", ErrConflict, errResp.Error())
	default:
		return response, &errResp
	}
}
