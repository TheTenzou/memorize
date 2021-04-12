package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"memorize/controllers"
	"memorize/repository"
	"memorize/services"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func inject(sources *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	userRepository := repository.NewUserRepository(sources.DB)
	tokenRepository := repository.NewTokenRepository(sources.RedisClient)

	userService := services.NewUserService(&services.UserServiceConfig{
		UserRepository: userRepository,
	})

	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	privateKeyString, err := ioutil.ReadFile(privateKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyString)

	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	publicKeyFile := os.Getenv("PUBLIC_KEY_FILE")
	publicKeyString, err := ioutil.ReadFile(publicKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyString)

	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	tokenExpiration := os.Getenv("TOKEN_EXP")
	refreshTokenExpiration := os.Getenv("REFRESH_TOKEN_EXP")

	tokenExpirationSec, err := strconv.ParseInt(tokenExpiration, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse TOKEN_EXP as int: %v", err)
	}

	refreshTokenExpirationSec, err := strconv.ParseInt(refreshTokenExpiration, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %v", err)
	}

	tokenService := services.NewTokenService(&services.TokenServiceConfig{
		TokenRepository:           tokenRepository,
		PrivateKey:                privateKey,
		PublicKey:                 publicKey,
		RefreshSecret:             refreshSecret,
		TokenExpirationSec:        tokenExpirationSec,
		RefreshTokenExpirationSec: refreshTokenExpirationSec,
	})

	router := gin.Default()

	baseUrl := os.Getenv("ACCOUNT_API_URL")
	controllers.NewController(&controllers.Config{
		Router:       router,
		UserService:  userService,
		TokenService: tokenService,
		BaseURL:      baseUrl,
	})

	return router, nil
}
