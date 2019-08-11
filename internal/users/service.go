package users

import (
	"context"
	"github.com/go-kit/kit/log"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	APIToken string `json:"api_token"`
}

type Stat struct {
}

type Service interface {
	GetUserByFromRequest(ctx context.Context, apiToken string) (*User, error)
}

type Repository interface {
	GetUserByApiToken(ctx context.Context, apiToken string) (*User, error)
}

func NewService(log log.Logger, repository Repository) Service {
	return &UserService{log: log, repo:repository}
}

type UserService struct {
	log log.Logger
	repo Repository

}

func (s UserService) GetUserByFromRequest(ctx context.Context, apiToken string) (*User, error) {
	return s.repo.GetUserByApiToken(ctx, apiToken)
}

func (s UserService) Close() error {
	return nil
}
