package mysql

import (
	"database/sql"
	_ "database/sql/driver"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	dbdriver   string
	dbUser     string
	dbPassword string
	dbHost     string
	dbDatabase string
	dbPort     string
}

type Connector struct {
	log    log.Logger
	config *Config
	db     *sql.DB
}

const (
	DB_CONNECTION = "DB_CONNECTION"
	DB_USER       = "DB_USERNAME"
	DB_HOST       = "DB_HOST"
	DB_PASSWORD   = "DB_PASSWORD"
	DB_DATABASE   = "DB_DATABASE"
	DB_PORT       = "DB_PORT"
)

func NewConnector(logger log.Logger, envMap map[string]string) Connector {
	return Connector{
		log: logger,
		config: &Config{envMap[DB_CONNECTION],
			envMap[DB_USER],
			envMap[DB_PASSWORD],
			envMap[DB_HOST],
			envMap[DB_DATABASE],
			envMap[DB_PORT]},
	}
}

// Connect returns a sql.DB that can be used to query.
// Not keeping this persistent because for a small application think it is more trouble than it is worth
func (c Connector) Connect() (db *sql.DB) {
	level.Info(c.log).Log("msg", "connecting to mysql")
	level.Info(c.log).Log("msg", "returning connection")
	_ = c.config.dbdriver
	dbUser := c.config.dbUser
	dbHost := c.config.dbHost
	dbPass := c.config.dbPassword
	dbName := c.config.dbDatabase
	dbPort := c.config.dbPort

	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	c.db = db
	return c.db

}
