package service

import (
	"crypto/rsa"
	"log"
	"memorize/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type TokenCustomClaims struct {
	User *models.User `json:"user"`
	jwt.StandardClaims
}

func generateToken(user *models.User, key *rsa.PrivateKey, expiration int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExpire := unixTime + expiration

	claims := TokenCustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExpire,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return signedToken, nil
}

type RefreshToken struct {
	SignedToken string
	ID          string
	ExpiresIn   time.Duration
}

type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

func generateRefreshToken(uid uuid.UUID, key string, expiration int64) (*RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(expiration) * time.Second)
	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := RefreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &RefreshToken{
		SignedToken: signedToken,
		ID:          tokenID.String(),
		ExpiresIn:   tokenExp.Sub(currentTime),
	}, nil
}
