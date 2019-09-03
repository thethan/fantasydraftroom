package auth

import (
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/fantasy"
	"golang.org/x/oauth2"
	"net/http"
)

const ENVVAR_CONSUMER_KEY = "CONSUMER_KEY"
const ENVVAR_CONSUMER_SECRET = "CONSUMER_SECRET"

func NewAuthService(log log.Logger, clientID, clientSecret string) AuthService {
	return AuthService{log: log, ClientID: clientID, ClientSecret: clientSecret}
}

type AuthService struct {
	log          log.Logger
	ClientID     string
	ClientSecret string
	mysql        sql.Conn

	config *oauth2.Config
	client *fantasy.Client
}

func (as *AuthService) SaveClient(c *http.Client) error {
	as.client = fantasy.NewClient(c)
	as.log.Log("msg", "saving client to authservice")
	level.Debug(as.log).Log("client", c)
	level.Debug(as.log).Log("as.client", as.client)
	return nil

}

// ReturnGoff is returning the version of the client
func (as *AuthService) ReturnGoff() (*fantasy.Client, error) {
	level.Debug(as.log).Log("msg", "ReturnGoff")

	if as.client != nil {
		return as.client, nil
	}
	level.Debug(as.log).Log("msg", "goff == nil", "as.client", as.client)

	return nil, errors.New("could not get client. Please try initializing it")
}

// GetConfig returns the config for the GetOauth2Config
func (as AuthService) GetConfig() *oauth2.Config {
	return fantasy.GetOAuth2Config(as.ClientID, as.ClientSecret, "https://fantasydraftroom.com/go/yahoo/callback")
}
