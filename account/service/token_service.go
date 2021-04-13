package service

import (
	"context"
	"crypto/rsa"
	"log"
	"memorize/models"
	"memorize/models/apperrors"

	"github.com/google/uuid"
)

type tokenService struct {
	TokenRepository           models.TokenRepository
	PrivateKey                *rsa.PrivateKey
	PublicKey                 *rsa.PublicKey
	RefreshSecret             string
	TokenExpirationSec        int64
	RefreshTokenExpirationSec int64
}

// parameter for creating token service
type TokenServiceConfig struct {
	TokenRepository           models.TokenRepository
	PrivateKey                *rsa.PrivateKey
	PublicKey                 *rsa.PublicKey
	RefreshSecret             string
	TokenExpirationSec        int64
	RefreshTokenExpirationSec int64
}

// function for initializing a UserService with its repository layer dependencie
func NewTokenService(config *TokenServiceConfig) models.TokenService {
	return &tokenService{
		TokenRepository:           config.TokenRepository,
		PrivateKey:                config.PrivateKey,
		PublicKey:                 config.PublicKey,
		RefreshSecret:             config.RefreshSecret,
		TokenExpirationSec:        config.TokenExpirationSec,
		RefreshTokenExpirationSec: config.RefreshTokenExpirationSec,
	}
}

// create new pair of tokens
// if previous token included, the previous token is removed
func (service *tokenService) NewPairFromUser(
	ctx context.Context,
	user *models.User,
	previousTokenID string,
) (*models.TokenPair, error) {

	if previousTokenID != "" {
		if err := service.TokenRepository.DeleteRefreshToken(ctx, user.UID.String(), previousTokenID); err != nil {
			log.Printf("Cold not delete previous refresh token for uid: %v, tokne %v\n", user.UID.String(), previousTokenID)

			return nil, err
		}
	}

	accessToken, err := generateToken(user, service.PrivateKey, service.TokenExpirationSec)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v, Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(user.UID, service.RefreshSecret, service.RefreshTokenExpirationSec)

	if err != nil {
		log.Printf("Error genaraating refreshToken for uid: %v. Error %v\n", user.UID, err.Error())
	}

	if err := service.TokenRepository.SetRefreshToken(
		ctx, user.UID.String(),
		refreshToken.ID.String(),
		refreshToken.ExpiresIn,
	); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	if previousTokenID != "" {
		if err := service.TokenRepository.DeleteRefreshToken(
			ctx,
			user.UID.String(),
			previousTokenID,
		); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", user.UID.String(), previousTokenID)
		}
	}

	return &models.TokenPair{
		AccessToken: models.AccessToken{
			Token: accessToken,
		},
		RefreshToken: models.RefreshToken{
			ID:     refreshToken.ID,
			UserID: user.UID,
			Token:  refreshToken.SignedToken,
		},
	}, nil
}

// validates the id token jwt string
// it returns the user extract from the IDTokenCustomClaims
func (service *tokenService) ValidateAccessToken(tokenString string) (*models.User, error) {
	claims, err := validateAccessToken(tokenString, service.PublicKey) // uses public RSA key

	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid
func (s *tokenService) ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &models.RefreshToken{
		Token:  tokenString,
		ID:     tokenUUID,
		UserID: claims.UserID,
	}, nil
}
