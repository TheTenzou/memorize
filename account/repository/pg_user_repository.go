package repository

import (
	"context"
	"log"
	"memorize/models"
	"memorize/models/apperrors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgUserRepository struct {
	DB *sqlx.DB
}

// factory for initializing user repository
func NewUserRepository(db *sqlx.DB) models.UserRepository {
	return &pgUserRepository{
		DB: db,
	}
}

// create user record in database
func (repository *pgUserRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING *"

	if err := repository.DB.GetContext(ctx, user, query, user.Login, user.Password); err != nil {

		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf(
				"Could not create a user with login: %v. Reason %v\n",
				user.Login,
				err.Code.Name(),
			)

			return apperrors.NewConflict("login", user.Login)
		}

		log.Printf("Could not create a user with login: %v. Reason: %v", user.Login, err)
	}

	return nil
}

// fetch user by id from database
func (repository *pgUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user := &models.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	if err := repository.DB.GetContext(ctx, user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}

// fetch user by login from databse
func (repository *pgUserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	user := &models.User{}

	query := "SELECT * FROM users WHERE login=$1"

	if err := repository.DB.GetContext(ctx, user, query, login); err != nil {
		log.Printf("Unable to get user with login: %v. err %v\n", login, err)
		return nil, apperrors.NewNotFound("login", login)
	}

	return user, nil
}

// Update updates a user's properties
// TODO fix user mutation
func (repository *pgUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET name=:name, email=:email, website=:website
		WHERE uid=:uid
		RETURNING *;
	`

	preparedQuery, err := repository.DB.PrepareNamedContext(ctx, query)

	if err != nil {
		log.Printf("Unable to prepare user update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if err := preparedQuery.GetContext(ctx, user, user); err != nil {
		log.Printf("Failed to update details for user: %v\n", user)
		return apperrors.NewInternal()
	}

	return nil
}
