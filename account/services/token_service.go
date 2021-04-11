package services

import (
	"context"
	"crypto/rsa"
	"log"
	"memorize/models"
	"memorize/models/apperrors"
)

type tokenService struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

type TokenServiceConfig struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

func NewTokenService(config *TokenServiceConfig) models.TokenService {
	return &tokenService{
		PrivateKey:    config.PrivateKey,
		PublicKey:     config.PublicKey,
		RefreshSecret: config.RefreshSecret,
	}
}

func (service *tokenService) NewPairFromUser(
	ctx context.Context,
	user *models.User,
	previousTokenID string,
) (*models.TokenPair, error) {
	idToken, err := generateToken(user, service.PrivateKey)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v, Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(user.UID, service.RefreshSecret)

	if err != nil {
		log.Printf("Error genaraating refreshToken for uid: %v. Error %v\n", user.UID, err.Error())
	}

	return &models.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SignedToken,
	}, nil
}
