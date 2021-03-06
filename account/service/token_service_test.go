package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"memorize/mocks"
	"memorize/models"
	"memorize/models/apperrors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPairFromUser(test *testing.T) {
	// init params for token service
	var tokenExpiration int64 = 15 * 60
	var refreshTokenExpiration int64 = 3 * 24 * 60 * 60

	privateRSA, _ := ioutil.ReadFile("../rsa/rsa_private_test.pem")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateRSA)
	publicRSA, _ := ioutil.ReadFile("../rsa/rsa_public_test.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicRSA)
	secret := "verysecretsecret"

	mockTokenRepository := new(mocks.MockTokenRepository)

	// instance of tested service
	tokenService := NewTokenService(&TokenServiceConfig{
		TokenRepository:           mockTokenRepository,
		PrivateKey:                privateKey,
		PublicKey:                 publicKey,
		RefreshSecret:             secret,
		TokenExpirationSec:        tokenExpiration,
		RefreshTokenExpirationSec: refreshTokenExpiration,
	})

	// init users
	uid, _ := uuid.NewRandom()
	user := &models.User{
		UID:      uid,
		Login:    "Alice",
		Password: "alicepassword",
	}

	uidErrorCase, _ := uuid.NewRandom()
	userErrorCase := &models.User{
		UID:      uidErrorCase,
		Email:    "failure",
		Password: "failurePassword",
	}
	previousID := "a_previous_tokenID"

	// Setup mock call responses in setup before test.Run statements
	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		user.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		uidErrorCase.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPreviousIDArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		user.UID.String(),
		previousID,
	}

	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("Error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPreviousIDArguments...).Return(nil)

	test.Run("Returns a token pair with proper values", func(test *testing.T) {
		ctx := context.Background()
		tokenPair, err := tokenService.NewPairFromUser(ctx, user, previousID)
		assert.NoError(test, err)

		mockTokenRepository.AssertCalled(test, "SetRefreshToken", setSuccessArguments...)
		mockTokenRepository.AssertCalled(test, "DeleteRefreshToken", deleteWithPreviousIDArguments...)

		var s string
		assert.IsType(test, s, tokenPair.AccessToken.Token)

		idTokenClaims := &idTokenCustomClaims{}

		// parse token claims
		_, err = jwt.ParseWithClaims(
			tokenPair.AccessToken.Token,
			idTokenClaims,
			func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			},
		)

		assert.NoError(test, err)

		expectedClaims := []interface{}{
			user.UID,
			user.Email,
			user.Name,
			user.ImageURL,
			user.Website,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Name,
			idTokenClaims.User.ImageURL,
			idTokenClaims.User.Website,
		}

		assert.ElementsMatch(test, expectedClaims, actualIDClaims)
		assert.Empty(test, idTokenClaims.User.Password)

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(tokenExpiration) * time.Second)
		assert.WithinDuration(test, expectedExpiresAt, expiresAt, 5*time.Second)

		// parse refresh token
		refreshTokenClaims := &refreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(
			tokenPair.RefreshToken.Token,
			refreshTokenClaims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		)

		assert.IsType(test, s, tokenPair.RefreshToken.Token)

		assert.NoError(test, err)
		assert.Equal(test, user.UID, refreshTokenClaims.UserID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshTokenExpiration) * time.Second)
		assert.WithinDuration(test, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	test.Run("Error setting refresh token", func(test *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, userErrorCase, "")
		assert.Error(test, err)

		mockTokenRepository.AssertCalled(test, "SetRefreshToken", setErrorArguments...)
		mockTokenRepository.AssertNotCalled(test, "DeleteRefreshToken")
	})

	test.Run("Empty string provided for prevID", func(test *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, user, "")
		assert.NoError(test, err)

		mockTokenRepository.AssertCalled(test, "SetRefreshToken", setSuccessArguments...)
		mockTokenRepository.AssertNotCalled(test, "DeleteRefreshToken")
	})
}

func TestValidateToken(test *testing.T) {
	var idExp int64 = 15 * 60

	private, _ := ioutil.ReadFile("../rsa/rsa_private_test.pem")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(private)
	public, _ := ioutil.ReadFile("../rsa/rsa_public_test.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(public)

	tokenService := NewTokenService(&TokenServiceConfig{
		PrivateKey:         privateKey,
		PublicKey:          publicKey,
		TokenExpirationSec: idExp,
	})

	uid, _ := uuid.NewRandom()
	user := &models.User{
		UID:      uid,
		Login:    "alice",
		Password: "passwordpass",
	}

	test.Run("Valid token", func(test *testing.T) {
		token, _ := generateToken(user, privateKey, idExp)

		uFromToken, err := tokenService.ValidateAccessToken(token)
		assert.NoError(test, err)

		assert.ElementsMatch(
			test,
			[]interface{}{user.Email, user.Name, user.UID, user.Website, user.ImageURL},
			[]interface{}{uFromToken.Email, uFromToken.Name, uFromToken.UID, uFromToken.Website, uFromToken.ImageURL},
		)
	})

	test.Run("Expired token", func(test *testing.T) {
		token, _ := generateToken(user, privateKey, -1)

		expectedErr := apperrors.NewAuthorization("Unable to verify user from idToken")

		_, err := tokenService.ValidateAccessToken(token)
		assert.EqualError(test, err, expectedErr.Message)
	})

	test.Run("Invalid signature", func(test *testing.T) {
		token, _ := generateToken(user, privateKey, -1)

		expectedErr := apperrors.NewAuthorization("Unable to verify user from idToken")

		_, err := tokenService.ValidateAccessToken(token)
		assert.EqualError(test, err, expectedErr.Message)
	})
}

func TestValidateRefreshToken(test *testing.T) {
	var refreshExp int64 = 3 * 24 * 2600
	secret := "anotsorandomtestsecret"

	tokenService := NewTokenService(&TokenServiceConfig{
		RefreshSecret:             secret,
		RefreshTokenExpirationSec: refreshExp,
	})

	uid, _ := uuid.NewRandom()
	user := &models.User{
		UID:      uid,
		Login:    "alice",
		Password: "passwordsssss",
	}

	test.Run("Valid token", func(t *testing.T) {
		testRefreshToken, _ := generateRefreshToken(user.UID, secret, refreshExp)

		validatedRefreshToken, err := tokenService.ValidateRefreshToken(testRefreshToken.SignedToken)
		assert.NoError(t, err)

		assert.Equal(t, user.UID, validatedRefreshToken.UserID)
		assert.Equal(t, testRefreshToken.SignedToken, validatedRefreshToken.Token)
		assert.Equal(t, user.UID, validatedRefreshToken.UserID)
	})

	test.Run("Expired token", func(test *testing.T) {
		testRefreshToken, _ := generateRefreshToken(user.UID, secret, -1)

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken.SignedToken)
		assert.EqualError(test, err, expectedErr.Message)
	})
}

func TestSignout(test *testing.T) {
	mockTokenRepository := new(mocks.MockTokenRepository)
	tokenService := NewTokenService(&TokenServiceConfig{
		TokenRepository: mockTokenRepository,
	})

	test.Run("No error", func(test *testing.T) {
		uidSuccess, _ := uuid.NewRandom()
		mockTokenRepository.
			On("DeleteUserRefreshTokens", mock.AnythingOfType("*context.emptyCtx"), uidSuccess.String()).
			Return(nil)

		ctx := context.Background()
		err := tokenService.Signout(ctx, uidSuccess)
		assert.NoError(test, err)
	})

	test.Run("Error", func(test *testing.T) {
		uidError, _ := uuid.NewRandom()
		mockTokenRepository.
			On("DeleteUserRefreshTokens", mock.AnythingOfType("*context.emptyCtx"), uidError.String()).
			Return(apperrors.NewInternal())

		ctx := context.Background()
		err := tokenService.Signout(ctx, uidError)

		assert.Error(test, err)

		apperr, ok := err.(*apperrors.Error)
		assert.True(test, ok)
		assert.Equal(test, apperr.Type, apperrors.Internal)
	})
}
