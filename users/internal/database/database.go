package database

import (
	"context"
	"fmt"
	"log"

	"com.MixieMelts.users/internal/models"
	"github.com/jackc/pgx/v5"
)

// DB represents the database connection.
type DB struct {
	*pgx.Conn
}

// New creates a new database connection.
func New(config string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database.")

	dbWrapper := &DB{conn}
	if err := dbWrapper.createTables(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	dbWrapper.Seed(context.Background())

	return dbWrapper, nil
}

func (db *DB) createTables(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(ctx, query)
	return err
}

// CreateUser inserts a new user into the database.
func (db *DB) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id;
	`
	var userID int64
	err := db.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}

// GetUserByEmail retrieves a user from the database by their email.
func (db *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1"
	err := db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, fmt.Errorf("user not found in database with emal: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user from the database by their ID.
func (db *DB) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1"
	err := db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No user found is not an error
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}
