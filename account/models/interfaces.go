package models

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Signup(context.Context, *User) error
}

type UserRepository interface {
	FindByID(context.Context, uuid.UUID) (*User, error)
}
