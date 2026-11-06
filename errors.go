package trakt

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrRateLimited is returned when Trakt responds with HTTP 429 Too Many Requests.
	ErrRateLimited = errors.New("trakt: rate limit exceeded")

	// ErrUnauthorized is returned when Trakt responds with HTTP 401 Unauthorized.
	ErrUnauthorized = errors.New("trakt: unauthorized or invalid credentials")

	// ErrNotFound is returned when Trakt responds with HTTP 404 Not Found.
	ErrNotFound = errors.New("trakt: resource not found")

	// ErrConflict is returned when Trakt responds with HTTP 409 Conflict.
	ErrConflict = errors.New("trakt: resource conflict or already exists")
)

// ErrorResponse represents an error returned by the Trakt API.
type ErrorResponse struct {
	Response         *http.Response `json:"-"`
	Message          string         `json:"message,omitempty"`
	ErrorDescription string         `json:"error_description,omitempty"`
	RateLimit        RateLimitInfo  `json:"-"`
}

func (r *ErrorResponse) Error() string {
	var statusCode int
	if r.Response != nil {
		statusCode = r.Response.StatusCode
	}

	msg := r.Message
	if msg == "" && r.ErrorDescription != "" {
		msg = r.ErrorDescription
	}

	if msg != "" {
		return fmt.Sprintf("trakt: HTTP %d: %s", statusCode, msg)
	}
	return fmt.Sprintf("trakt: HTTP %d", statusCode)
}

// IsRateLimited reports whether err was caused by an HTTP 429 rate limit response.
func IsRateLimited(err error) bool {
	if errors.Is(err, ErrRateLimited) {
		return true
	}
	var errResp *ErrorResponse
	return errors.As(err, &errResp) && errResp.Response != nil && errResp.Response.StatusCode == http.StatusTooManyRequests
}

// IsNotFound reports whether err was caused by an HTTP 404 Not Found response.
func IsNotFound(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	var errResp *ErrorResponse
	return errors.As(err, &errResp) && errResp.Response != nil && errResp.Response.StatusCode == http.StatusNotFound
}

// IsUnauthorized reports whether err was caused by an HTTP 401 Unauthorized response.
func IsUnauthorized(err error) bool {
	if errors.Is(err, ErrUnauthorized) {
		return true
	}
	var errResp *ErrorResponse
	return errors.As(err, &errResp) && errResp.Response != nil && errResp.Response.StatusCode == http.StatusUnauthorized
}
