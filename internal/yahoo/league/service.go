package league

import (
	"github.com/go-kit/kit/log"
	"github.com/thethan/fantasydraftroom/internal/yahoo/auth"
)

type LeagueService struct {
	log    log.Logger
	client *auth.AuthService
}

//func (svc LeagueService) GetLeagueSettings() {
//	client, err := svc.client.ReturnGoff()
//	if err != nil {
//		level.Error(svc.log).Log("msg", "error getting client from auth service", "err", err)
//	}
//	content, err := client.Provider.Get(
//		fmt.Sprintf("%s/league/%s/settings",
//			goff.YahooBaseURL,
//			"390.l.705710", ))
//
//	content
//}
