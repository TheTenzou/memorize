package service

import (
	"crypto/rsa"
	"fmt"
	"log"
	"memorize/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// token claims
type idTokenCustomClaims struct {
	User *models.User `json:"user"`
	jwt.StandardClaims
}

// generateIDToken generates an IDToken which is a jwt with myCustomClaims
func generateToken(user *models.User, key *rsa.PrivateKey, expiration int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExpire := unixTime + expiration

	claims := idTokenCustomClaims{
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

// refreshTokenData holds the actual signed jwt string along with the ID
// We return the id so it can be used without re-parsing the JWT from signed string
type refreshTokenData struct {
	SignedToken string
	ID          uuid.UUID
	ExpiresIn   time.Duration
}

type refreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// generateRefreshToken creates a refresh token
// The refresh token stores only the user's ID, a string
func generateRefreshToken(uid uuid.UUID, key string, expiration int64) (*refreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(expiration) * time.Second)
	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := refreshTokenCustomClaims{
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

	return &refreshTokenData{
		SignedToken: signedToken,
		ID:          tokenID,
		ExpiresIn:   tokenExp.Sub(currentTime),
	}, nil
}

// returns the token's claims if the token is valid
func validateAccessToken(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*idTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}
