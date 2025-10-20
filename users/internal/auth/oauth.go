package auth

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func GetJWTSecretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

// --- OAUTH CONFIGURATION ---
// In a real app, these would be loaded from environment variables.
var googleOAuthConfig *oauth2.Config
var facebookOAuthConfig *oauth2.Config

func init() {
	// Google OAuth Config
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/users/oauth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"), // You must set these env vars
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Facebook OAuth Config
	facebookOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/users/oauth/facebook/callback",
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"), // You must set these env vars
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
}

func GetGoogleAuthUrl(state string) string {
	return googleOAuthConfig.AuthCodeURL(state)
}

func GetGoogleOAuthExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return googleOAuthConfig.Exchange(ctx, code)
}

func GetFacebookAuthUrl(state string) string {
	return facebookOAuthConfig.AuthCodeURL(state)
}

func GetFacebookOAuthExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return facebookOAuthConfig.Exchange(ctx, code)
}
