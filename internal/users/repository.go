package users

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/thethan/fantasydraftroom/internal/mysql"
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
