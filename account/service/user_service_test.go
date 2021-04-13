package service

import (
	"context"
	"fmt"
	"memorize/mocks"
	"memorize/models"
	"memorize/models/apperrors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResponce := &models.User{
			UID:   uid,
			Email: "Alice@alice.com",
			Name:  "Alice",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		userService := NewUserService(&UserServiceConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResponce, nil)

		ctx := context.TODO()
		user, err := userService.GetUser(ctx, uid)

		assert.NoError(test, err)
		assert.Equal(test, user, mockUserResponce)
		mockUserRepository.AssertExpectations(test)
	})

	test.Run("Error", func(test *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		userService := NewUserService(&UserServiceConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("some error down the call chain"))

		ctx := context.TODO()

		user, err := userService.GetUser(ctx, uid)

		assert.Nil(test, user)
		assert.Error(test, err)
		mockUserRepository.AssertExpectations(test)
	})
}

func TestSignup(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &models.User{
			Login:    "alice",
			Password: "alicepassword",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		userservice := NewUserService(&UserServiceConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.
			On(
				"Create",
				mock.AnythingOfType("*context.emptyCtx"),
				mockUser,
			).
			Run(
				func(args mock.Arguments) {
					// arg 0 is context, arg 1 is *User
					userArg := args.Get(1).(*models.User)
					userArg.UID = uid
				},
			).
			Return(nil)

		ctx := context.TODO()
		err := userservice.Signup(ctx, mockUser)

		assert.NoError(test, err)

		assert.Equal(test, uid, mockUser.UID)

		mockUserRepository.AssertExpectations(test)
	})

	test.Run("Error", func(test *testing.T) {
		mockUser := &models.User{
			Login:    "alice",
			Password: "alicepassword",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&UserServiceConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("login", mockUser.Login)

		mockUserRepository.
			On(
				"Create",
				mock.AnythingOfType("*context.emptyCtx"),
				mockUser,
			).
			Return(mockErr)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.EqualError(test, err, mockErr.Error())

		mockUserRepository.AssertExpectations(test)
	})
}

func TestSignin(test *testing.T) {
	login := "login"
	validPassword := "testPasswrod"
	hashedValidPassword, _ := HashPassword(validPassword)
	invalidPasswrod := "invalid"

	mockUserRepository := new(mocks.MockUserRepository)
	userService := NewUserService(&UserServiceConfig{
		UserRepository: mockUserRepository,
	})

	test.Run("Success", func(test *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &models.User{
			Login:    login,
			Password: validPassword,
		}

		mockUserResponse := &models.User{
			UID:      uid,
			Login:    login,
			Password: hashedValidPassword,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			login,
		}

		mockUserRepository.
			On("FindByLogin", mockArgs...).Return(mockUserResponse, nil)

		ctx := context.TODO()
		loginUser, err := userService.Signin(ctx, mockUser)
		assert.Equal(test, mockUserResponse, loginUser)

		assert.NoError(test, err)
		mockUserRepository.AssertCalled(test, "FindByLogin", mockArgs...)
	})

	test.Run("Invalid login/password combination", func(test *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &models.User{
			Login:    login,
			Password: invalidPasswrod,
		}

		mockUserResponse := &models.User{
			UID:      uid,
			Login:    login,
			Password: hashedValidPassword,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			login,
		}

		mockUserRepository.
			On("FindByLogin", mockArgs...).Return(mockUserResponse, nil)

		ctx := context.TODO()
		loginUser, err := userService.Signin(ctx, mockUser)
		assert.Nil(test, loginUser)

		assert.Error(test, err)
		assert.EqualError(test, err, "Invalid login and password combination")
		mockUserRepository.AssertCalled(test, "FindByLogin", mockArgs...)
	})
}
