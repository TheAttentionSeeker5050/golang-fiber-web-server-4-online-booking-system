package config

import (
	"log"
	"os"

	"github.com/markbates/goth/providers/google"
	"golang.org/x/oauth2"
)

type OauthBaseConfig struct {
	GoogleLoginConfig oauth2.Config
}

var OauthConfig OauthBaseConfig

func InitGoogleConfig() oauth2.Config {
	clientID := os.Getenv("OAUTH2_GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET")
	clientCallbackURL := os.Getenv("OAUTH2_GOOGLE_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	// This is a package that simplifies the OAuth process for Go applications
	OauthConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  clientCallbackURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	return OauthConfig.GoogleLoginConfig
}
