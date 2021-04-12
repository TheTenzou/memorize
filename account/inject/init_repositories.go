package inject

import (
	"memorize/models"
	"memorize/repository"
)

type Repositories struct {
	UserRepository  models.UserRepository
	TokenRepository models.TokenRepository
}

func InitRepositories(sources *DataSources) *Repositories {
	userRepository := repository.NewUserRepository(sources.DB)
	tokenRepository := repository.NewTokenRepository(sources.RedisClient)

	return &Repositories{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}
}
