package users

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/thethan/fantasydraftroom/internal/mysql"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/auth"
	"time"
)

type MysqlRepository struct {
	connector *mysql.Connector
	log       log.Logger
}

func NewMysqlRepository(connector *mysql.Connector, logger log.Logger) MysqlRepository{
	return MysqlRepository{
		connector: connector,
		log : logger,
	}
}

func (r MysqlRepository) GetUserByApiToken(ctx context.Context, apiToken string) (*User, error) {
	db := r.connector.Connect()
	defer db.Close()
	stmt, err := db.Prepare("SELECT id, name, email, api_token FROM fdr_users where api_token = ? ORDER BY id DESC LIMIT 1 ")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(apiToken)

	var user User
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.APIToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}


// GetYahooToken retrieves the yahoo token from the database
func (r MysqlRepository) GetYahooToken(ctx context.Context) (*auth.YahooAuth, error) {
	db := r.connector.Connect()
	defer db.Close()

	stmt, err := db.Prepare("SELECT access_token, token_type, expires_in, refresh_token, xoauth_yahoo_guid, created_at, updated_at FROM fdr_yahoo_tokens WHERE user_id = ?")
	if err != nil {
		return nil, err
	}
	userID := ctx.Value(USER)
	row := stmt.QueryRow(userID)
	var yahooAuth auth.YahooAuth

	var createdAt, updatedAt mysql.NullTime
	err = row.Scan(&yahooAuth.AccessToken, &yahooAuth.TokenType, &yahooAuth.ExpiresIn, &yahooAuth.RefreshToken, &yahooAuth, &yahooAuth.XoauthYahooGuid, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	yahooAuth.UpdateToken(yahooAuth.AccessToken, yahooAuth.TokenType, yahooAuth.RefreshToken, yahooAuth.XoauthYahooGuid, &updatedAt.Time, time.Duration(time.Second*yahooAuth.ExpiresIn))
	return &yahooAuth, nil

}
