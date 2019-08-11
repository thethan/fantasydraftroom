package users

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func NewUserMiddleware(log log.Logger, svc Service) UserMiddleware {
	return UserMiddleware{log: log, service: svc}
}

type UserMiddleware struct {
	log     log.Logger
	service Service
}

func (u UserMiddleware) GetAPITokenFromEndpoint(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		apiToken := ctx.Value(APIToken).(string)
		user, err := u.service.GetUserByFromRequest(ctx, apiToken)
		if err != nil {
			return nil, err
		}
		// @todo i do not like this
		ctx = context.WithValue(ctx, USER, user)

		return next(ctx, request)

	}
}
