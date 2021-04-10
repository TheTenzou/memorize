package models

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Get(context context.Context, uid uuid.UUID) (*User, error)
}

type UserRepository interface {
	FindByID(context context.Context, uid uuid.UUID) (*User, error)
}
