package repository

import (
	"context"
	"fmt"
	"log"
	"memorize/models"
	"memorize/models/apperrors"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

// factory for initializing User Repositories
func NewTokenRepository(redisClient *redis.Client) models.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// stores a refresh token with an expiry time
func (repository *redisTokenRepository) SetRefreshToken(
	ctx context.Context,
	userID string,
	tokenID string,
	expiresIn time.Duration,
) error {

	key := fmt.Sprintf("%s.%s", userID, tokenID)
	if err := repository.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}

// delete old refresh tokens
func (repository *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {

	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := repository.Redis.Del(ctx, key).Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}
