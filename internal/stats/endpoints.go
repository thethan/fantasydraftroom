package stats

import (
	"context"
	"github.com/thethan/fantasydraftroom/internal/users"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Set struct {
	Import endpoint.Endpoint
}

type UserRequest struct {
	UserUD int `json:"player_id" `
}

type StatImportResponse struct {
	Message string `json:"message"`
}

// New endpoints
func New( logger log.Logger, svc Service, usersMiddleware users.UserMiddleware) Set {
	var statImportEndpoint endpoint.Endpoint
	{
		statImportEndpoint = MakeStatImportEndpoint(svc)
		//statImportEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(statImportEndpoint)
		//statImportEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(statImportEndpoint)
		//statImportEndpoint = LoggingMiddleware(log.With(logger, "method", "Concat"))(statImportEndpoint)
		statImportEndpoint = usersMiddleware.GetAPITokenFromEndpoint(statImportEndpoint)
	}
	return Set{Import: statImportEndpoint}
}

// MakeStatImportEndpoint constructs a Sum endpoint wrapping the service.
func MakeStatImportEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctxUser := ctx.Value(users.USER)
		user := ctxUser.(*users.User)
		return StatImportResponse{Message: "hello user..." + user.Name}, nil
	}
}
