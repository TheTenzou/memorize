package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type dataSources struct {
	DB *sqlx.DB
}

func initDataSources() (*dataSources, error) {
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
	db, err := sqlx.Open("postgres", pgConnString)

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return &dataSources{
		DB: db,
	}, nil
}

func (d *dataSources) close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	return nil
}
