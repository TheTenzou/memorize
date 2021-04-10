package services

import (
	"context"
	"memorize/models"

	"github.com/google/uuid"
)

type UserService struct {
	UserRespository models.UserRepository
}

type UserServiceConfig struct {
	UserRepository models.UserRepository
}

func NewUserService(c *UserServiceConfig) models.UserService {
	return &UserService{
		UserRespository: c.UserRepository,
	}
}

func (u *UserService) Get(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user, err := u.UserRespository.FindByID(ctx, uid)

	return user, err
}

func (u *UserService) Signup(ctx context.Context, user *models.User) error {
	panic("Method not implemented")
}
