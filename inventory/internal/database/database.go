package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// DB is a thin wrapper around a pgx connection.
type DB struct {
	*pgx.Conn
}

// New creates a new database connection and ensures required tables exist.
func New(config string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db := &DB{Conn: conn}

	if err := db.createTables(context.Background()); err != nil {
		_ = conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Optional seed or migrations could run here
	db.seed(context.Background())

	log.Println("inventory: successfully connected to database")
	return db, nil
}

func (db *DB) createTables(ctx context.Context) error {
	if err := db.createIngredientsTable(ctx); err != nil {
		return err
	}
	if err := db.createInventoryAdjustmentsTable(ctx); err != nil {
		return err
	}
	if err := db.createRecipeItemsTable(ctx); err != nil {
		return err
	}
	return nil
}

// seed performs optional initial data seeding. Keep minimal so tests can run deterministically.
func (db *DB) seed(ctx context.Context) {
	db.seedIngredients(ctx)
	db.seedInventoryAdjustments(ctx)
	db.seedRecipeItems(ctx)
}
