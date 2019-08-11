package players

import (
	"context"
	"github.com/go-kit/kit/log"
	"sync"
)



type PlayerID int
type PlayerIndex int

type Results struct {
	PlayerIndex PlayerIndex
	PlayerID    PlayerID
}

type playerToIndex map[PlayerID]PlayerIndex

type Service interface {
	GetPlayersByList(ctx context.Context, draftID, userID int) (map[string]playerToIndex, error)
}

type DraftResultRepository interface {
	GetDraftResults(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan <- Results) error
}
type DefaultPlayerRankRepository interface {
	GetDefaultPlayerRank(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan <- Results) error
}
type UserPlayerRankRepository interface {
	GetUserPlayerRank(ctx context.Context,  wg *sync.WaitGroup, draftID int, userID int,  resultsChan chan <- Results) error
}

func NewService(log log.Logger, draftResultRepository DraftResultRepository, defaultRankRepository DefaultPlayerRankRepository, userRankRepository UserPlayerRankRepository) Service {
	return &PlayerService{log: log, draftResultRepository: draftResultRepository, defaultRankRepository: defaultRankRepository, userRankRepository: userRankRepository,}
}

type PlayerService struct {
	log                   log.Logger
	defaultRankRepository DefaultPlayerRankRepository
	userRankRepository    UserPlayerRankRepository
	draftResultRepository DraftResultRepository
}

func (s PlayerService) GetPlayersByList(ctx context.Context, draftID int, userID int) (map[string]playerToIndex, error) {

	var results, defaultRank, userRank playerToIndex
	results = make(map[PlayerID]PlayerIndex)
	defaultRank = make(map[PlayerID]PlayerIndex)
	userRank = make(map[PlayerID]PlayerIndex)

	resultsChan := make(chan Results, 1)
	defaultOrderChan := make(chan Results, 1)
	_ = make(chan Results, 1)
	errorChan := make(chan error, 3)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		errorChan <- s.draftResultRepository.GetDraftResults(ctx, wg, draftID, resultsChan)
	}()

	go func() {
		errorChan <- s.defaultRankRepository.GetDefaultPlayerRank(ctx, wg, draftID, defaultOrderChan)
	}()

	// map the results
	go func() {
		for res := range resultsChan {
			results[res.PlayerID] = res.PlayerIndex
		}
	}()

	// map the results
	go func() {
		for res := range defaultOrderChan {
			defaultRank[res.PlayerID] = res.PlayerIndex
		}
	}()

	wg.Wait()
	close(resultsChan)
	close(defaultOrderChan)

	return map[string]playerToIndex{
		"results_to_player":      results,
		"default_rank": defaultRank,
		"user_rank":    userRank,
	}, nil

}


