package leagues

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/auth"
	"github.com/thethan/fantasydraftroom/pkg/yahoo/fantasy"
)

func NewService(logger log.Logger, authSvc *auth.AuthService) Service {
	return Service{log: logger, client: authSvc}
}

type Service struct {
	log    log.Logger
	client *auth.AuthService
}

// Get League returns a specific leagues and its information for a user.
func (svc Service) GetLeague(ctx context.Context, leagueID string) (*fantasy.League, error) {
	client, err := svc.client.ReturnGoff(ctx,1)
	if err != nil {
		level.Error(svc.log).Log("msg", "error getting client from auth service", "err", err)
	}
	league, err := client.GetLeagueMetadata(leagueID)
	return league, err
}
