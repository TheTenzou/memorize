package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"memorize/mocks"
	"memorize/models"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPairFromUser(test *testing.T) {
	var tokenExpiration int64 = 15 * 60
	var refreshTokenExpiration int64 = 3 * 24 * 60 * 60

	privateRSA, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateRSA)
	publicRSA, _ := ioutil.ReadFile("../rsa_public_test.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicRSA)
	secret := "verysecretsecret"

	mockTokenRepository := new(mocks.MockTokenRepository)

	tokenService := NewTokenService(&TokenServiceConfig{
		TokenRepository:           mockTokenRepository,
		PrivateKey:                privateKey,
		PublicKey:                 publicKey,
		RefreshSecret:             secret,
		TokenExpirationSec:        tokenExpiration,
		RefreshTokenExpirationSec: refreshTokenExpiration,
	})

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
		assert.IsType(test, s, tokenPair.IDToken)

		idTokenClaims := &TokenCustomClaims{}

		_, err = jwt.ParseWithClaims(
			tokenPair.IDToken,
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

		refreshTokenClaims := &RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(
			tokenPair.RefreshToken,
			refreshTokenClaims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		)

		assert.IsType(test, s, tokenPair.RefreshToken)

		assert.NoError(test, err)
		assert.Equal(test, user.UID, refreshTokenClaims.UID)

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
