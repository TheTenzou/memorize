package inject

import (
	"fmt"
	"io/ioutil"
	"memorize/models"
	"memorize/service"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

// structure for services
type Services struct {
	UserService  models.UserService
	TokenService models.TokenService
}

// Inject repositories into services
func InitServices(repositories *Repositories) (*Services, error) {

	userService := service.NewUserService(&service.UserServiceConfig{
		UserRepository: repositories.UserRepository,
	})

	tokenService, err := initTokenService(repositories.TokenRepository)
	if err != nil {
		return nil, fmt.Errorf("faild to init tokenService: %w", err)
	}

	return &Services{
		UserService:  userService,
		TokenService: tokenService,
	}, nil
}

func initTokenService(tokenRepository models.TokenRepository) (models.TokenService, error) {

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

	return service.NewTokenService(&service.TokenServiceConfig{
		TokenRepository:           tokenRepository,
		PrivateKey:                privateKey,
		PublicKey:                 publicKey,
		RefreshSecret:             refreshSecret,
		TokenExpirationSec:        tokenExpirationSec,
		RefreshTokenExpirationSec: refreshTokenExpirationSec,
	}), nil
}
