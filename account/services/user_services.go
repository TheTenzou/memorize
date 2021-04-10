package services

import (
	"context"
	"log"
	"memorize/models"
	"memorize/models/apperrors"

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
	password, err := hashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for login: %v\n", user.Login)
		return apperrors.NewInternal()
	}

	user.Password = password

	if err := u.UserRespository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
