package stats

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Stat struct {
}

type Service interface {
	ImportPlayer(ctx context.Context, id int) (Stat, error)
}

func NewService(log log.Logger) Service {
	return &StatService{log: log}
}

type StatService struct {
	log log.Logger
}

func (s StatService) ImportPlayer(ctx context.Context, id int) (Stat, error) {
	return Stat{}, nil
}

func (s StatService) Close() error {
	return nil
}
