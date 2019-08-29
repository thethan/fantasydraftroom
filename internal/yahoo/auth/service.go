package auth

import (
	"database/sql"
	"errors"
	"github.com/thethan/fantasydraftroom/internal/yahoo/fantasy"
	"golang.org/x/oauth2"
	"net/http"
)

const ENVVAR_CONSUMER_KEY = "CONSUMER_KEY"
const ENVVAR_CONSUMER_SECRET = "CONSUMER_SECRET"

func NewAuthService(clientID, clientSecret string) AuthService {
	return AuthService{ClientID: clientID, ClientSecret: clientSecret}
}

type AuthService struct {
	ClientID     string
	ClientSecret string
	mysql        sql.Conn

	config *oauth2.Config
	client *fantasy.Client
}

func (as *AuthService) SaveClient(c *http.Client) error {
	as.client = fantasy.NewClient(c)
	return nil

}
// ReturnGoff is returning the version of the client
func (as *AuthService) ReturnGoff() (*fantasy.Client, error) {
	if as.client != nil {
		return as.client, nil
	}

	return nil, errors.New("could not get client. Please try initializing it")
}

// GetConfig returns the config for the GetOauth2Config
func (as AuthService) GetConfig() *oauth2.Config {
	return fantasy.GetOAuth2Config(as.ClientID, as.ClientSecret, "https://fantasydraftroom.com/go/yahoo/callback")
}
