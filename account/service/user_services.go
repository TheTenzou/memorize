package service

import (
	"context"
	"log"
	"memorize/models"
	"memorize/models/apperrors"

	"github.com/google/uuid"
)

type userService struct {
	UserRepository models.UserRepository
}

// config hold repositories that will eventually be injected into this this service layer
type UserServiceConfig struct {
	UserRepository models.UserRepository
}

// factory function for initializing a UserService with its repository layer dependencies
func NewUserService(config *UserServiceConfig) models.UserService {
	return &userService{
		UserRepository: config.UserRepository,
	}
}

// fetch user by uid
func (s *userService) GetUser(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	return s.UserRepository.FindByID(ctx, uid)
}

// signup user if login avaliable
func (s *userService) Signup(ctx context.Context, user *models.User) error {
	password, err := HashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for login: %v\n", user.Login)
		return apperrors.NewInternal()
	}

	user.Password = password

	if err := s.UserRepository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) Signin(ctx context.Context, user *models.User) (*models.User, error) {
	fetchedUser, err := s.UserRepository.FindByLogin(ctx, user.Login)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return nil, apperrors.NewAuthorization("User with this login dont exist")
	}

	// verify password
	match, err := comparePasswords(fetchedUser.Password, user.Password)

	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization("Invalid login and password combination")
	}

	return fetchedUser, nil
}

// update user details
func (s *userService) UpdateDetails(ctx context.Context, user *models.User) error {
	return s.UserRepository.Update(ctx, user)
}
