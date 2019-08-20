package players

import (
	"context"
	"errors"
	"github.com/thethan/fantasydraftroom/internal/users"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Set struct {
	PlayersOrder     endpoint.Endpoint
	PlayerPreference endpoint.Endpoint
}

type DraftPlayerRankingsRequest struct {
	DraftID int
	UserID  int
}

type UserPlayerPreferenceRequest struct {
	DraftID int
	UserID  int
	Body    struct {
		PlayerIDs []PlayerID `json:"players"`
	} `json:"data"`
}

type GenericData struct {
	Data interface{} `json:"data"`
}
type DraftPlayerLists struct {
	DefaultPlayersOrder []PlayerID    `json:"default_players_order"`
	Results             playerToIndex `json:"results"`
	UserPlayerOrder     []PlayerID    `json:"user_players_order"`
}

// New endpoints for players order stuff
func New(logger log.Logger, usersMiddleware users.UserMiddleware, svc Service) Set {
	var playersOrder endpoint.Endpoint
	{
		playersOrder = MakePlayersOrderEndpoint(svc)
		playersOrder = users.LoggingMiddleware(log.With(logger, "method", "GetPlayersByList"))(playersOrder)

		playersOrder = usersMiddleware.GetAPITokenFromEndpoint(playersOrder)
	}

	var playerPreference endpoint.Endpoint
	{
		playerPreference = MakeUserPlayerPreference(svc)
		playerPreference = users.LoggingMiddleware(log.With(logger, "method", "UserPlayerPreference"))(playerPreference)

		playerPreference = usersMiddleware.GetAPITokenFromEndpoint(playerPreference)
	}

	return Set{PlayersOrder: playersOrder, PlayerPreference: playerPreference}
}

// MakePlayersOrderEndpoint constructs a Sum endpoint wrapping the service.
func MakePlayersOrderEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DraftPlayerRankingsRequest)
		ctxUser := ctx.Value(users.USER)
		user, assertion := ctxUser.(*users.User)
		if assertion == false {
			return nil, errors.New("could not get user")
		}
		draftPlayerLists, err := svc.GetPlayersByList(ctx, req.DraftID, user.ID)
		return GenericData{
			Data: draftPlayerLists,
		}, err
	}
}

// MakePlayersOrderEndpoint constructs a Sum endpoint wrapping the service.
func MakeUserPlayerPreference(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserPlayerPreferenceRequest)
		ctxUser := ctx.Value(users.USER)
		user, assertion := ctxUser.(*users.User)
		if assertion == false {
			return nil, errors.New("could not get user")
		}
		err = svc.SaveUsersPlayerList(ctx, req.DraftID, user.ID, req.Body.PlayerIDs)
		return GenericData{
			Data: nil,
		}, err
	}
}
