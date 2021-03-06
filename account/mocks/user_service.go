package mocks

import (
	"context"
	"memorize/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, uid)

	var r0 *models.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*models.User)
	}

	var r1 error

	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) Signup(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}

func (m *MockUserService) Signin(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)

	var r0 *models.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*models.User)
	}

	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) UpdateDetails(ctx context.Context, u *models.User) error {
	args := m.Called(ctx, u)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}
