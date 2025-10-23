package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/giakiet05/lkforum/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserInfo holds the essential user information returned from Google.
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

var googleOauthConfig *oauth2.Config

// InitGoogleOAuthConfig initializes the Google OAuth2 configuration.
// This should be called once at application startup.
func InitGoogleOAuthConfig() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.Cfg.Google.RedirectURL,
		ClientID:     config.Cfg.Google.ClientID,
		ClientSecret: config.Cfg.Google.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// GetGoogleLoginURL generates the URL for the user to log in with Google.
func GetGoogleLoginURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state)
}

// GetGoogleUserInfo exchanges the authorization code for user info.
func GetGoogleUserInfo(code string) (*GoogleUserInfo, error) {
	// Exchange the authorization code for an access token.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange code for token: " + err.Error())
	}

	// Use the access token to get the user's information.
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.New("failed to get user info: " + err.Error())
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("failed to read user info response: " + err.Error())
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		return nil, errors.New("failed to unmarshal user info: " + err.Error())
	}

	if userInfo.Email == "" {
		return nil, errors.New("email not found in Google user info")
	}

	return &userInfo, nil
}
