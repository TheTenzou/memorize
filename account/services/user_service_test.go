package services

import (
	"context"
	"fmt"
	"memorize/mocks"
	"memorize/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
		user, err := userService.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, user, mockUserResponce)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		userService := NewUserService(&UserServiceConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("some error down the call chain"))

		ctx := context.TODO()

		user, err := userService.Get(ctx, uid)

		assert.Nil(t, user)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}
