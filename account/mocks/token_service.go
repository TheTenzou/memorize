package mocks

import (
	"context"
	"memorize/models"

	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) NewPairFromUser(
	ctx context.Context,
	user *models.User,
	prevTokenID string,
) (*models.TokenPair, error) {

	args := m.Called(ctx, user, prevTokenID)

	var r0 *models.TokenPair
	if args.Get(0) != nil {
		r0 = args.Get(0).(*models.TokenPair)
	}

	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}

	return r0, r1
}
