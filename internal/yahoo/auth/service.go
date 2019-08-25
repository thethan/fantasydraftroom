package auth

import (
	"database/sql"
	"golang.org/x/oauth2"
	"net/http"
)
import "github.com/Forestmb/goff"

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
	client *goff.Client
}

func (as *AuthService) SaveClient(c *http.Client) error {

	goff.NewClient(c)
	return nil

}

func (as AuthService) GetConfig() *oauth2.Config {
	return goff.GetOAuth2Config(as.ClientID, as.ClientSecret, "https://fantasydraftroom.com/go/yahoo/callback")
}