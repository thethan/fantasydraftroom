package players

import (
	"context"
	"fmt"
	log2 "github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	mysql2 "github.com/thethan/fantasydraftroom/internal/mysql"
	"sync"
	"testing"
)

func TestMysqlRepository_GetDraftResults(t *testing.T) {
	playerToIndexChan := make(chan Results, 1)

	playerToIdxMap := make(playerToIndex)
	envMap, err := godotenv.Read("../../.env")
	if err != nil {
		panic("could not read .env file")
	}

	log := log2.NewNopLogger()
	mysql := mysql2.NewConnector(log, envMap)

	repo := NewMysqlRepository(&mysql, log)
	ctx := context.Background()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = repo.GetDraftResults(ctx, wg, 1, playerToIndexChan)
	}()

	go func() {
		for res := range playerToIndexChan {
			playerToIdxMap[res.PlayerID] = res.PlayerIndex
		}
	}()

	wg.Wait()
	close(playerToIndexChan)
	assert.Nil(t, err)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	assert.True(t, len(playerToIdxMap) > 12, "playerDraft does not equal 12 results")
}



func TestMysqlRepository_GetDefaultPlayerRank(t *testing.T) {
	resultsChan := make(chan Results, 1)

	playerToIdxMap := make(playerToIndex)
	envMap, err := godotenv.Read("../../.env")
	if err != nil {
		panic("could not read .env file")
	}

	log := log2.NewNopLogger()
	mysql := mysql2.NewConnector(log, envMap)

	repo := NewMysqlRepository(&mysql, log)
	ctx := context.Background()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = repo.GetDefaultPlayerRank(ctx, wg, 1, resultsChan)
	}()

	go func() {
		for res := range resultsChan {
			playerToIdxMap[res.PlayerID] = res.PlayerIndex
		}
	}()

	wg.Wait()
	close(resultsChan)
	assert.Nil(t, err)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	assert.Len(t, playerToIdxMap, 1182, "player count does not equal number of players in the database")
}
