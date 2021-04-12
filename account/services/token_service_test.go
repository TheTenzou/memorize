package services

import (
	"context"
	"io/ioutil"
	"memorize/models"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPairFromUser(test *testing.T) {
	var tokenExpiration int64 = 15 * 60
	var refreshTokenExpiration int64 = 3 * 24 * 2600
	privateRSA, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateRSA)
	publicRSA, _ := ioutil.ReadFile("../rsa_public_test.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicRSA)
	secret := "verysecretsecret"

	tokenService := NewTokenService(&TokenServiceConfig{
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

	test.Run("Returns a token pair with proper values", func(test *testing.T) {
		ctx := context.TODO()
		tokenPair, err := tokenService.NewPairFromUser(ctx, user, "")
		assert.NoError(test, err)

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
		expectedExpiresAt := time.Now().Add(time.Duration(tokenExpiration) * time.Minute)
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
		expectedExpiresAt = time.Now().Add(time.Duration(refreshTokenExpiration) * time.Hour)
		assert.WithinDuration(test, expectedExpiresAt, expiresAt, 5*time.Second)
	})
}
