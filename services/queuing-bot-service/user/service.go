package user

import (
	"context"
	"queuingbot/domain"
	"queuingbot/internal/repository/postgresql"
)

type Service struct {
	repo *postgresql.UserRepository
}

func NewUserService(repo *postgresql.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s Service) StoreUser(user *domain.User) error {
	err := s.repo.Store(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
