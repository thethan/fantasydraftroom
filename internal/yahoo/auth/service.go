package auth

import (
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/oauth2"
	"log"

)
import "github.com/Forestmb/goff"


func NewAuthService(clientID, clientSecret string) AuthService {
	return AuthService{ClientID:clientID, ClientSecret:clientSecret}
}
type AuthService struct {
	ClientID     string
	ClientSecret string
	mysql        sql.Conn

	config *oauth2.Config
}

func (as *AuthService) GetClient() {
	ctx := context.Background()
	as.config = goff.GetOAuth2Config(as.ClientID, as.ClientSecret, "https://fantasydraftroom.com/go/yahoo/callback")

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := as.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)


	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := as.config.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := as.config.Client(ctx, tok)

	fmt.Printf("%v", client)

}

