package user

import (
	"context"
	"files-bot-service/domain"
	"files-bot-service/internal/repository/postgresql"
)

type Service struct {
	Repository *postgresql.UserRepository
}

func NewUserService(repo *postgresql.UserRepository) *Service {
	return &Service{
		Repository: repo,
	}
}

func (s Service) CreateUser(id int64) (user domain.User, err error) {
	user = domain.NewUser(id)
	err = s.Repository.Store(context.Background(), &user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s Service) Save(usr domain.User) (user domain.User, err error) {
	err = s.Repository.Store(context.Background(), &usr)
	if err != nil {
		return usr, err
	}
	return user, nil
}

func (s Service) Update(usr domain.User) (domain.User, error) {
	err := s.Repository.Update(context.Background(), &usr)
	return usr, err
}
