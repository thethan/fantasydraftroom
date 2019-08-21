package auth

import (
	"fmt"
	"github.com/Forestmb/goff"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSomething(t *testing.T) {
	svc := AuthService{}
	client := svc.GetClient()
	goffClient := goff.NewClient(client)

	teams, _ := goffClient.GetAllTeams("390.l.705710")
	for idx := range teams {
		fmt.Printf("%+v\n",teams[idx])
	}
	assert.True(t, true)
}
