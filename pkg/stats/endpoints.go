package stats

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Set struct {
	Import endpoint.Endpoint
}

type StatImportRequest struct {
	PlayerID int `json:"player_id" `
}

type StatImportResponse struct {
	Message string `json:"message"`
}

// New endpoints
func New( logger log.Logger, svc Service) Set {
	var statImportEndpoint endpoint.Endpoint
	{
		statImportEndpoint = MakeStatImportEndpoint(svc)
		//statImportEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(statImportEndpoint)
		//statImportEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(statImportEndpoint)
		//statImportEndpoint = LoggingMiddleware(log.With(logger, "method", "Concat"))(statImportEndpoint)
		//statImportEndpoint = InstrumentingMiddleware(duration.With("method", "Concat"))(statImportEndpoint)
	}
	return Set{Import: statImportEndpoint}
}

// MakeStatImportEndpoint constructs a Sum endpoint wrapping the service.
func MakeStatImportEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(StatImportRequest)
		_, err = svc.ImportPlayer(ctx, 1)
		return StatImportResponse{Message: "hello ethan..."}, nil
	}
}
