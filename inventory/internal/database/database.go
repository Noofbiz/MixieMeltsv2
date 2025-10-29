package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a thin wrapper around a pgx connection pool.
type DB struct {
	*pgxpool.Pool
}

// New creates a new database connection pool and ensures required tables exist.
func New(config string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	db := &DB{pool}

	if err := db.createTables(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Optional seed or migrations could run here
	db.seed(context.Background())

	log.Println("inventory: successfully created database connection pool")
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
	// Seed recipe items (mirrors products' recipes). Enabled so inventory has recipe data for product lookups.
	db.seedRecipeItems(ctx)
}
