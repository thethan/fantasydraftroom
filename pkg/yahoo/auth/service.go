package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/thethan/fantasydraftroom/internal/fdr/php/users"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/fantasy"
	"golang.org/x/oauth2"
)


const ENVVAR_CONSUMER_KEY = "CONSUMER_KEY"
const ENVVAR_CONSUMER_SECRET = "CONSUMER_SECRET"
const USERID = 1

func NewAuthService(log log.Logger, clientID, clientSecret string) AuthService {
	return AuthService{log: log, ClientID: clientID, ClientSecret: clientSecret}
}

type AuthService struct {
	log          log.Logger
	ClientID     string
	ClientSecret string

	userRepo     *users.Repository
	mysql        sql.Conn

	config *oauth2.Config
	client *fantasy.Client

	userToClients map[int]*fantasy.Client
}

func (as *AuthService) AuthenticateUser(ctx context.Context, yahooTokenCode string) error {
	config := as.GetConfig()

	tok, err := config.Exchange(ctx, yahooTokenCode)
	if err != nil {
		level.Error(as.log).Log("msg", "error in logging user information to yahoo")
		return err
	}

	level.Info(as.log).Log("msg", "saving client to authservice")
	err = as.userRepo.SaveYahooToken(USERID, tok)
	if err != nil {
		level.Error(as.log).Log("msg", "error saving yahoo token")
		return err
	}

	client := config.Client(ctx, tok)
	as.userToClients[USERID] = fantasy.NewClient(client)

	return nil

}

// ReturnGoff is returning the version of the client
func (as *AuthService) ReturnGoff(userID int) (*fantasy.Client, error) {
	level.Debug(as.log).Log("msg", "ReturnGoff")

	if client, ok :=  as.userToClients[userID]; ok {
		return client, nil
	}
	level.Debug(as.log).Log("msg", "goff == nil", "as.client", as.client)

	return nil, errors.New("could not get client. Please try initializing it")
}

// GetConfig returns the config for the GetOauth2Config
func (as AuthService) GetConfig() *oauth2.Config {
	return fantasy.GetOAuth2Config(as.ClientID, as.ClientSecret, "https://fantasydraftroom.com/go/yahoo/callback")
}