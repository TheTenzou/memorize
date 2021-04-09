package mocks

import (
	"context"

	"memorize/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, uid)

	var r0 *model.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*model.User)
	}

	var r1 error

	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}

	return r0, r1
}
