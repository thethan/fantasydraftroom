package players

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/thethan/fantasydraftroom/internal/users"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/auth"
)

type Set struct {
	PlayersOrder     endpoint.Endpoint
	PlayerPreference endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
	Callback         endpoint.Endpoint
	LeagueEndpoint   endpoint.Endpoint
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

type UserYahoo struct {
	Code  string
	State string
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
func New(logger log.Logger, usersMiddleware users.UserMiddleware, svc Service, auth *auth.AuthService) Set {
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

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(logger, auth)
		loginEndpoint = users.LoggingMiddleware(log.With(logger, "method", "LoginEndpoint"))(loginEndpoint)
	}

	var leagueEndpoint endpoint.Endpoint
	{
		leagueEndpoint = MakeLeagueEndpoint(logger, auth)
		leagueEndpoint = users.LoggingMiddleware(log.With(logger, "method", "LeagueEndpoint"))(leagueEndpoint)
	}

	return Set{PlayersOrder: playersOrder, PlayerPreference: playerPreference, LoginEndpoint: loginEndpoint, LeagueEndpoint: leagueEndpoint}
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

// MakePlayersOrderEndpoint constructs a Sum endpoint wrapping the service.
func MakeLoginEndpoint(logger log.Logger, svc *auth.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserYahoo)

		if err != nil {
			level.Error(logger).Log("msg", "could not exchange token", "err", err)
			return nil, err
		}
		// todo save user if not userID
		err = svc.AuthenticateUser(ctx, r.Code)
		if err  != nil {
			return nil, err
		}
		level.Info(logger).Log("msg", "could not", "err", err)

		ff, err := svc.ReturnGoff(auth.USERID)

		leagues, err := ff.GetUserLeagues("2019")
		if err != nil {
			return nil, err
		}

		return leagues, err

	}
}

// MakeLeagueEndpoint constructs a Sum endpoint wrapping the service.
func MakeLeagueEndpoint(logger log.Logger, svc *auth.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ff, err := svc.ReturnGoff()
		if err != nil {
			level.Error(logger).Log("msg", "could not exchange token", "err", err)
			return nil, err
		}
		level.Info(logger).Log("msg", "getting leagues info")
		leagues, err := ff.GetUserLeagues("2019")
		return leagues, nil

	}
}
