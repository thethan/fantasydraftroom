package users

import (
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/oauth2"
)

const SaveYahooToken  = "SaveYahooToken"

type Repository struct {
	log         log.Logger
	db          *sql.DB
	statements map[string]*sql.Stmt
}

type User struct {
	Name     string
	Email    string
	APIToken string
}

func (r Repository) GetUser() {

}

func (r Repository) prepareSaveYahooToken() (*sql.Stmt, error) {
	stmt, err := r.db.Prepare("REPLACE INTO fdr_yahoo_tokens (access_token, token_type, expires_in, refresh_token, updated_at) VALUES (? , ? , ? , ? , NOW()) WHERE user_id = ?")
	if err != nil {
		level.Error(r.log).Log("msg", "error in prerparing SaveYahooToken sql", "error", err)
		return nil, err
	}
	r.statements[SaveYahooToken] = stmt
	return stmt, nil
}



func (r Repository) SaveYahooToken(userID int, token *oauth2.Token) error {
	var err error
	var stmt *sql.Stmt
	if _, ok := r.statements[SaveYahooToken]; !ok {
		stmt, err = r.prepareSaveYahooToken()
		if err != nil {
			level.Error(r.log).Log("msg", "error in SaveYahooToken: after preparing the stmt", "error", err)
			return nil
		}
	}
	stmt = r.statements[SaveYahooToken]

	_, err = stmt.Exec(token.AccessToken, token.TokenType, token.Expiry, token.RefreshToken, userID)
	if err != nil {
		level.Error(r.log).Log("msg", "error in SaveYahooToken: post exec", "error", err)
		return nil
	}
	return nil
}
