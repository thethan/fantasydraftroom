package auth

import (
	"context"
	"fmt"
	gokitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/thethan/fantasydraftroom/internal/fdr/php/users"
	"github.com/thethan/fantasydraftroom/internal/mysql"
	"testing"
)

func TestSomething(t *testing.T) {
	assert.True(t, true)
}


func TestNewAuthService(t *testing.T) {
	log := gokitlog.NewNopLogger()
	mysqlConnector := mysql.NewConnector(log, map[string]string{
		"DB_USERNAME": "NeilDiamond",
		"DB_PASSWORD": "JazzSinger" ,
		"DB_HOST":"192.168.99.109",
		"DB_DATABASE":"minikube",
		"DB_PORT": "32000",
	})
	repo := users.NewRepository(log, &mysqlConnector)
	svc := NewAuthService(log, "", "", &repo)
	client, err := svc.ReturnGoff(context.Background(), USERID)

	assert.Nil(t, err)
	assert.NotNil(t, client)

	leagues, err := client.GetUserLeagues("2019")
	assert.Nil(t, err)
	fmt.Printf("%v\n", err)
	assert.NotNil(t, leagues)



}