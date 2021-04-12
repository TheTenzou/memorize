package inject

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type DataSources struct {
	DB          *sqlx.DB
	RedisClient *redis.Client
}

func InitDataSources() (*DataSources, error) {

	database, err := initPostgres()
	if err != nil {
		return nil, err
	}

	redisDB, err := initRedis()
	if err != nil {
		return nil, err
	}

	return &DataSources{
		DB:          database,
		RedisClient: redisDB,
	}, nil
}

func (d *DataSources) Close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing Redis clinet: %w", err)
	}

	return nil
}

func initPostgres() (*sqlx.DB, error) {

	log.Printf("Initilazing data sources\n")

	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")
	pgSSL := os.Getenv("POSTGRES_SSL")

	pgConnString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL,
	)

	log.Printf("Connecting to Postgresql\n")
	database, err := sqlx.Open("postgres", pgConnString)

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return database, nil
}

func initRedis() (*redis.Client, error) {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	log.Printf("Connecting to Redis\n")
	redisDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := redisDB.Ping(context.Background()).Result()

	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return redisDB, nil
}
