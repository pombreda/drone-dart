package googlestorage

import (
	"time"

	"code.google.com/p/goauth2/oauth"
)

const (
	Scope       = "https://www.googleapis.com/auth/devstorage.read_write"
	AuthURL     = "https://accounts.google.com/o/oauth2/auth"
	TokenURL    = "https://accounts.google.com/o/oauth2/token"
	RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
)

func MakeOauthTransport(clientId string, clientSecret string, refreshToken string) *oauth.Transport {
	return &oauth.Transport{
		&oauth.Config{
			ClientId:     clientId,
			ClientSecret: clientSecret,
			Scope:        Scope,
			AuthURL:      AuthURL,
			TokenURL:     TokenURL,
			RedirectURL:  RedirectURL,
		},
		&oauth.Token{
			AccessToken:  "",
			RefreshToken: refreshToken,
			Expiry:       time.Time{}, // no expiry
		},
		nil,
	}
}
