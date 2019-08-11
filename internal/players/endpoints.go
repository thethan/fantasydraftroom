package players

import (
	"context"
	"errors"
	"github.com/thethan/fantasydraftroom/internal/users"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Set struct {
	PlayersOrder endpoint.Endpoint
}

type DraftPlayerRankingsRequest struct {
	DraftID int
	UserID  int
}

type DraftPlayerRankingsResponse struct {
	MapOfPlayerIDToIdx map[string]playerToIndex `json:"data"`
}

// New endpoints
func New(logger log.Logger, usersMiddleware users.UserMiddleware,  svc Service) Set {
	var playersOrder endpoint.Endpoint
	{
		playersOrder = MakePlayersOrderEndpoint(svc)
		playersOrder = users.LoggingMiddleware(log.With(logger, "method", "GetPlayersByList"))(playersOrder)

		playersOrder = usersMiddleware.GetAPITokenFromEndpoint(playersOrder)
	}
	return Set{PlayersOrder: playersOrder}
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
		mapOfPlayerIDToIdx, err := svc.GetPlayersByList(ctx, req.DraftID, user.ID)
		return DraftPlayerRankingsResponse{MapOfPlayerIDToIdx: mapOfPlayerIDToIdx}, err
	}
}
