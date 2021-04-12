package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Signup(context.Context, *User) error
}

type TokenService interface {
	NewPairFromUser(context.Context, *User, string) (*TokenPair, error)
}

type UserRepository interface {
	FindByID(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, *User) error
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, previousTokenID string) error
}
