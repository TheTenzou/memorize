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
func (r *redisTokenRepository) SetRefreshToken(
	ctx context.Context,
	userID string,
	tokenID string,
	expiresIn time.Duration,
) error {

	key := fmt.Sprintf("%s.%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}

// delete old refresh tokens
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {

	key := fmt.Sprintf("%s.%s", userID, tokenID)

	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/tokenID %s/%s doesnot exist\n", userID, tokenID)
		return apperrors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

// DeleteUserRefreshTokens delete all refresh tokens of specific user
func (r *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	interator := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for interator.Next(ctx) {
		if err := r.Redis.Del(ctx, interator.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", interator.Val())
			failCount++
		}
	}

	// check last value
	if err := interator.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", interator.Val())
	}

	if failCount > 0 {
		return apperrors.NewInternal()
	}

	return nil
}
