package service

import (
	"context"
	"memorize/model"

	"github.com/google/uuid"
)

type UserService struct {
	UserRespository model.UserRepository
}

type UserServiceConfig struct {
	UserRepository model.UserRepository
}

func NewUserService(c *UserServiceConfig) model.UserService {
	return &UserService{
		UserRespository: c.UserRepository,
	}
}

func (u *UserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user, err := u.UserRespository.FindByID(ctx, uid)

	return user, err
}
