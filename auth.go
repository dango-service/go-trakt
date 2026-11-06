package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// AuthService handles OAuth2 Authorization Code, token refresh, token revocation, and Device/PIN flows.
type AuthService struct {
	client *Client
}

// Token represents an OAuth2 token pair for Trakt.
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}

// String implements fmt.Stringer, redacting credentials to prevent secret leakage in logs.
func (t Token) String() string {
	return "[REDACTED]"
}

// GoString implements fmt.GoStringer, redacting credentials to prevent secret leakage in logs.
func (t Token) GoString() string {
	return "[REDACTED]"
}

// DeviceCodeResponse represents parameters returned by Trakt's Device Authorization endpoint.
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// AuthorizeURL returns the user-facing consent URL to begin the OAuth2 Authorization Code grant.
func (s *AuthService) AuthorizeURL(redirectURI, state string) string {
	baseURL := "https://trakt.tv/oauth/authorize"
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", s.client.clientID)
	if redirectURI != "" {
		params.Set("redirect_uri", redirectURI)
	}
	if state != "" {
		params.Set("state", state)
	}
	return baseURL + "?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for an OAuth2 token pair.
func (s *AuthService) ExchangeCode(ctx context.Context, code, redirectURI, clientSecret string) (*Token, *Response, error) {
	reqBody := map[string]string{
		"grant_type":    "authorization_code",
		"code":          strings.TrimSpace(code),
		"client_id":     s.client.clientID,
		"client_secret": strings.TrimSpace(clientSecret),
		"redirect_uri":  redirectURI,
	}
	return s.requestToken(ctx, "/oauth/token", reqBody)
}

// RefreshToken exchanges an active refresh token for a fresh token pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, clientSecret, redirectURI string) (*Token, *Response, error) {
	reqBody := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": strings.TrimSpace(refreshToken),
		"client_id":     s.client.clientID,
		"client_secret": strings.TrimSpace(clientSecret),
		"redirect_uri":  redirectURI,
	}
	return s.requestToken(ctx, "/oauth/token", reqBody)
}

// GenerateDeviceCode initiates the Device / PIN flow for headless devices (CLIs, TVs, NAS).
func (s *AuthService) GenerateDeviceCode(ctx context.Context) (*DeviceCodeResponse, *Response, error) {
	reqBody := map[string]string{
		"client_id": s.client.clientID,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("trakt: marshal device code request: %w", err)
	}

	var res DeviceCodeResponse
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/oauth/device/code", bytes.NewReader(bodyBytes), &res)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: generate device code: %w", err)
	}
	return &res, resp, nil
}

// PollDeviceToken polls the device token endpoint using the device code.
func (s *AuthService) PollDeviceToken(ctx context.Context, deviceCode, clientSecret string) (*Token, *Response, error) {
	reqBody := map[string]string{
		"code":          strings.TrimSpace(deviceCode),
		"client_id":     s.client.clientID,
		"client_secret": strings.TrimSpace(clientSecret),
	}
	return s.requestToken(ctx, "/oauth/device/token", reqBody)
}

// RevokeToken revokes an active access token.
func (s *AuthService) RevokeToken(ctx context.Context, token, clientSecret string) (*Response, error) {
	reqBody := map[string]string{
		"token":         strings.TrimSpace(token),
		"client_id":     s.client.clientID,
		"client_secret": strings.TrimSpace(clientSecret),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("trakt: marshal revoke request: %w", err)
	}

	resp, err := s.client.doRequest(ctx, http.MethodPost, "/oauth/revoke", bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return resp, fmt.Errorf("trakt: revoke token: %w", err)
	}
	return resp, nil
}

func (s *AuthService) requestToken(ctx context.Context, path string, payload map[string]string) (*Token, *Response, error) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("trakt: marshal token request: %w", err)
	}

	var token Token
	resp, err := s.client.doRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), &token)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: token request: %w", err)
	}
	return &token, resp, nil
}
