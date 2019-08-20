package players

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"sync"
)

const RankNotAvailable = 9000

type PlayerID int
type PlayerIndex int

type Results struct {
	PlayerIndex PlayerIndex
	PlayerID    PlayerID
}

type playerToIndex map[PlayerID]PlayerIndex

type Service interface {
	GetPlayersByList(ctx context.Context, draftID, userID int) (*DraftPlayerLists, error)
	SaveUsersPlayerList(ctx context.Context, draftID, userID int, playersList []PlayerID) ( error)
}

type DraftResultRepository interface {
	GetDraftResults(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan<- Results) error
}
type DefaultPlayerRankRepository interface {
	GetDefaultPlayerRank(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan<- Results) error
}
type UserPlayerRankRepository interface {
	GetUserPlayerRank(ctx context.Context, wg *sync.WaitGroup, draftID int, userID int, resultsChan chan<- Results) error
}

type SaveUserPlayerPreference interface {
	SaveUserPlayerPreference(ctx context.Context, wg *sync.WaitGroup, draftID int, userID int, playerID int, preferenceOrder int) error
	RemoveFromListIfNotIn(ctx context.Context, draftID int, userID int, playerList []int) error
}

func NewService(log log.Logger, draftResultRepository DraftResultRepository, defaultRankRepository DefaultPlayerRankRepository, userRankRepository UserPlayerRankRepository, preference SaveUserPlayerPreference) Service {
	return &PlayerService{log: log, draftResultRepository: draftResultRepository, defaultRankRepository: defaultRankRepository, userRankRepository: userRankRepository, saveUserPlayerPreferenceRepository: preference}
}

type PlayerService struct {
	log                                log.Logger
	defaultRankRepository              DefaultPlayerRankRepository
	userRankRepository                 UserPlayerRankRepository
	draftResultRepository              DraftResultRepository
	saveUserPlayerPreferenceRepository SaveUserPlayerPreference
}

func (s PlayerService) GetPlayersByList(ctx context.Context, draftID int, userID int) (*DraftPlayerLists, error) {

	var results playerToIndex
	results = make(map[PlayerID]PlayerIndex)
	userRank := make([]PlayerID, 0)
	defaultRank := make([]PlayerID, 0)

	resultsChan := make(chan Results, 1)
	defaultOrderChan := make(chan Results, 1)
	userRankChan := make(chan Results, 1)
	errorChan := make(chan error, 3)

	_ = level.Info(s.log).Log("msg", "starting go routines")

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		errorChan <- s.draftResultRepository.GetDraftResults(ctx, wg, draftID, resultsChan)
	}()

	go func() {
		errorChan <- s.defaultRankRepository.GetDefaultPlayerRank(ctx, wg, draftID, defaultOrderChan)
	}()

	go func() {
		errorChan <- s.userRankRepository.GetUserPlayerRank(ctx, wg, draftID, userID, userRankChan)
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
			defaultRank = append(defaultRank, res.PlayerID)
		}
	}()

	// map the results
	go func() {
		for res := range userRankChan {
			level.Info(s.log).Log("msg", "playerID to user rank")
			userRank = append(userRank, res.PlayerID)
		}
	}()

	go func() {
		for err := range errorChan {
			if err != nil {
				_ = level.Error(s.log).Log("msg", "error in getting player list", "error", err)
			}
		}
	}()

	wg.Wait()
	_ = level.Info(s.log).Log("msg", "wait groups closed, closing channels")





	close(resultsChan)
	close(defaultOrderChan)
	close(userRankChan)
	close(errorChan)

	draftPlayerList := &DraftPlayerLists{
		Results:             results,
		DefaultPlayersOrder: defaultRank,
		UserPlayerOrder: userRank,
	}

	return draftPlayerList, nil

}

func (s PlayerService) SaveUsersPlayerList(ctx context.Context, draftID, userID int, playersList []PlayerID) (error) {
	errorChan := make(chan error, len(playersList))

	intIDs := make([]int, len(playersList))
	for idx := range playersList {
		intIDs[idx] = int(playersList[idx])
	}
	err := s.saveUserPlayerPreferenceRepository.RemoveFromListIfNotIn(ctx,  draftID, userID, intIDs)
	if err != nil {
		return err
	}
	_ =level.Info(s.log).Log("method", "SaveUsersPlayerList", "number_of_preferences", len(playersList))

	wg := &sync.WaitGroup{}
	wg.Add(len(playersList))
	for idx := range playersList {
		_ =level.Info(s.log).Log("Saving Player ID into repo", playersList[idx])
		err := s.saveUserPlayerPreferenceRepository.SaveUserPlayerPreference(ctx, wg,  draftID, userID, intIDs[idx], idx+1)
		if err != nil {
			errorChan <- err
		}
	}
	wg.Wait()
	return nil
}
