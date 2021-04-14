package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
type UserService interface {
	// fetch user by uid
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	// signup user if login avaliable
	Signup(ctx context.Context, user *User) error
	// signin user if credentials are right
	Signin(ctx context.Context, user *User) (*User, error)
	// update user details
	UpdateDetails(ctx context.Context, user *User) error
}

// TokenService defines methods the handler layer expects
type TokenService interface {
	// create new pair of tokens
	NewPairFromUser(ctx context.Context, user *User, previousToken string) (*TokenPair, error)
	// validate token
	ValidateAccessToken(token string) (*User, error)
	// validate refresh token
	ValidateRefreshToken(token string) (*RefreshToken, error)
	// delete all users tokens
	Signout(ctx context.Context, userID uuid.UUID) error
}

// UserRepository defines methods the service layer expects
type UserRepository interface {
	// fetch user by id from database
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	// fetch user by login from database
	FindByLogin(ctx context.Context, login string) (*User, error)
	// create user record in database
	Create(ctx context.Context, user *User) error
	// update user record in database
	Update(ctx context.Context, user *User) error
}

// TokenRepository defines methods it expects a repository
type TokenRepository interface {
	// stores a refresh token with an expiry time
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	// delete old refresh tokens
	DeleteRefreshToken(ctx context.Context, userID string, previousTokenID string) error
	// delete allUsers refresh tokens
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}
