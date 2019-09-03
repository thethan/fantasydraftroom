package leagues

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/thethan/fantasydraftroom/internal/users"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/leagues"
)

type LeagueRequest struct {
	LeagueID string
}

type Set struct {
	LeagueEndpoint endpoint.Endpoint
}

// New endpoints for players order stuff
func New(logger log.Logger, usersMiddleware users.UserMiddleware, league *leagues.Service) Set {


	var leagueEndpoint endpoint.Endpoint
	{
		leagueEndpoint = makeLeagueEndpoint(logger, league)
		leagueEndpoint = users.LoggingMiddleware(log.With(logger, "method", "LeagueEndpoint"))(leagueEndpoint)
	}

	return Set{LeagueEndpoint: leagueEndpoint}
}

// MakeLeagueEndpoint constructs a Sum endpoint wrapping the service.
func makeLeagueEndpoint(logger log.Logger, svc *leagues.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(LeagueRequest)
		if !ok {
			return nil, errors.New("could not decode")
		}
		league, err := svc.GetLeague(ctx, req.LeagueID)
		if err != nil {
			level.Error(logger).Log("msg", "could not exchange token", "err", err)
			return nil, err
		}
		return league, nil

	}
}
