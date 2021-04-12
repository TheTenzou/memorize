package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
type UserService interface {
	// fetch user by uid
	GetUser(context.Context, uuid.UUID) (*User, error)
	// signup user if login avaliable
	Signup(context.Context, *User) error
}

// TokenService defines methods the handler layer expects
type TokenService interface {
	// create new pair of tokens
	NewPairFromUser(context.Context, *User, string) (*TokenPair, error)
}

// UserRepository defines methods the service layer expects
type UserRepository interface {
	// fetch user by id from database
	FindByID(context.Context, uuid.UUID) (*User, error)
	// create user record in database
	Create(context.Context, *User) error
}

// TokenRepository defines methods it expects a repository
type TokenRepository interface {
	// stores a refresh token with an expiry time
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	// delete old refresh tokens
	DeleteRefreshToken(ctx context.Context, userID string, previousTokenID string) error
}
