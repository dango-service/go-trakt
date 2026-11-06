package trakt

import (
	"context"
	"fmt"
	"net/http"
)

// UsersService handles user settings and profile details.
type UsersService struct {
	client *Client
}

// UserSettings represents the authenticated user's profile and account settings.
type UserSettings struct {
	User    UserProfile `json:"user"`
	Account UserAccount `json:"account"`
}

// UserProfile represents the public profile fields for a user.
type UserProfile struct {
	Username string  `json:"username"`
	Name     string  `json:"name"`
	VIP      bool    `json:"vip"`
	JoinedAt string  `json:"joined_at"`
	Location string  `json:"location"`
	About    string  `json:"about"`
	Gender   string  `json:"gender"`
	Age      int     `json:"age"`
	IDs      UserIDs `json:"ids"`
}

// UserIDs contains unique identifiers for a user.
type UserIDs struct {
	Slug string `json:"slug"`
	UUID string `json:"uuid"`
}

// UserAccount represents account configuration such as timezone.
type UserAccount struct {
	Timezone   string `json:"timezone"`
	DateFormat string `json:"date_format"`
	Time24hr   bool   `json:"time_24hr"`
	CoverImage string `json:"cover_image"`
}

// GetSettings retrieves the profile and settings of the authenticated user.
func (s *UsersService) GetSettings(ctx context.Context) (*UserSettings, *Response, error) {
	var settings UserSettings
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/users/settings", nil, &settings)
	if err != nil {
		return nil, resp, fmt.Errorf("trakt: get user settings: %w", err)
	}
	return &settings, resp, nil
}
