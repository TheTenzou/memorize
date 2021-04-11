package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"memorize/controllers"
	"memorize/repository"
	"memorize/services"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func inject(sources *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	userRepository := repository.NewUserRepository(sources.DB)

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

	tokenService := services.NewTokenService(&services.TokenServiceConfig{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		RefreshSecret: refreshSecret,
	})

	router := gin.Default()

	controllers.NewController(&controllers.Config{
		Router:       router,
		UserService:  userService,
		TokenService: tokenService,
	})

	return router, nil
}
