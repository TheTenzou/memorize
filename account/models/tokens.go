package models

import "github.com/google/uuid"

// RefreshToken store token properties
type RefreshToken struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
	Token  string    `json:"refreshToken"`
}

// AccessToken store token properties
type AccessToken struct {
	Token string `json:"accessToken"`
}

type TokenPair struct {
	RefreshToken
	AccessToken
}
