package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
type UserService interface {
	GetUser(context.Context, uuid.UUID) (*User, error)
	Signup(context.Context, *User) error
}

// TokenService defines methods the handler layer expects
type TokenService interface {
	NewPairFromUser(context.Context, *User, string) (*TokenPair, error)
}

// UserRepository defines methods the service layer expects
type UserRepository interface {
	FindByID(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, *User) error
}

// TokenRepository defines methods it expects a repository
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, previousTokenID string) error
}
