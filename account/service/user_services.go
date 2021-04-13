package service

import (
	"context"
	"log"
	"memorize/models"
	"memorize/models/apperrors"

	"github.com/google/uuid"
)

type userService struct {
	UserRespository models.UserRepository
}

// config hold repositories that will eventually be injected into this this service layer
type UserServiceConfig struct {
	UserRepository models.UserRepository
}

// factory function for initializing a UserService with its repository layer dependencies
func NewUserService(config *UserServiceConfig) models.UserService {
	return &userService{
		UserRespository: config.UserRepository,
	}
}

// fetch user by uid
func (u *userService) GetUser(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user, err := u.UserRespository.FindByID(ctx, uid)

	return user, err
}

// signup user if login avaliable
func (u *userService) Signup(ctx context.Context, user *models.User) error {
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

func (u *userService) Signin(ctx context.Context, user *models.User) error {
	fetchedUser, err := u.UserRespository.FindByLogin(ctx, user.Login)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return apperrors.NewAuthorization("User with this login dont exist")
	}

	// verify password
	match, err := comparePasswords(fetchedUser.Password, user.Password)

	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid login and password combination")
	}

	user = fetchedUser
	return nil
}
