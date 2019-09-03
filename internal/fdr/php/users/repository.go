package users

import (
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/auth"
)

type Repository struct {
	log         log.Logger
	db          *sql.DB
	authService *auth.AuthService
}

type User struct {
	Name     string
	Email    string
	APIToken string
}

func (r Repository) GetUser() {

}

func (r Repository) SaveYahooToken(userID int) {

}

