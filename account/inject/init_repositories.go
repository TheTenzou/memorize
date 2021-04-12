package inject

import (
	"memorize/models"
	"memorize/repository"
)

// Stractue for repositories
type Repositories struct {
	UserRepository  models.UserRepository
	TokenRepository models.TokenRepository
}

// Inject data sources into repositories
func InitRepositories(sources *DataSources) *Repositories {
	userRepository := repository.NewUserRepository(sources.DB)
	tokenRepository := repository.NewTokenRepository(sources.RedisClient)

	return &Repositories{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}
}
