package stats

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/thethan/fantasydraftroom/internal/yahoo/auth"
)

// Service Import stats for a leagues ID
type Service interface {
	GetLeagueStats(leagueID int) error
}

func NewStatService(log log.Logger, auth *auth.AuthService) StatService {
	return StatService{log: log, auth: auth}
}

type StatService struct {
	auth *auth.AuthService
	log  log.Logger
}

func (ss StatService) GetLeagueStats(leagueID int) error {
	client, err := ss.auth.ReturnGoff()
	if err != nil {
		level.Error(ss.log).Log("msg", "error getting Goff client", "err", err)
		return err
	}
	league, err := client.GetLeagueMetadata("390.l.705710")

	return nil
}
