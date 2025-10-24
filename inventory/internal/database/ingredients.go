package database

import (
	"context"
	"fmt"
	"log"

	"com.MixieMelts.inventory/internal/models"
	"github.com/jackc/pgx"
)

func (db *DB) seedIngredients(ctx context.Context) {
	// Only seed when the ingredients table is empty.
	var count int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM ingredients`).Scan(&count); err != nil {
		log.Printf("inventory: seed count query failed: %v", err)
		return
	}
	if count > 0 {
		log.Println("Ingredients table already seeded, skipping")
		return
	}

	ingredients := []models.Ingredient{
		{Name: "Vanilla", Type: "Scent", Unit: "ml", Stock: 500, MinThreshold: 100},
		{Name: "Chocolate", Type: "Scent", Unit: "ml", Stock: 500, MinThreshold: 100},
		{Name: "Maple Brown Sugar", Type: "Scent", Unit: "mg", Stock: 1000, MinThreshold: 100},
		{Name: "Milk", Type: "Dairy", Unit: "ml", Stock: 3000, MinThreshold: 500},
		{Name: "Butter", Type: "Dairy", Unit: "g", Stock: 1000, MinThreshold: 200},
	}

	for _, it := range ingredients {
		query := `INSERT INTO ingredients (name, type, unit, stock, min_threshold, notes)
		VALUES ($1,$2,$3,$4,$5,$6)`
		_, err := db.Exec(ctx, query, it.Name, it.Type, it.Unit, it.Stock, it.MinThreshold, it.Notes)
		if err != nil {
			log.Printf("inventory: seed insert %q: %v", it.Name, err)
		}
	}
}

func (db *DB) createIngredientsTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS ingredients (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		type TEXT NOT NULL,
		unit TEXT NOT NULL,
		stock DOUBLE PRECISION DEFAULT 0,
		min_threshold DOUBLE PRECISION DEFAULT 0,
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(ctx, query)
	return err
}

// GetIngredients returns all ingredients.
func (db *DB) GetIngredients(ctx context.Context) ([]models.Ingredient, error) {
	rows, err := db.Query(ctx, `SELECT id, name, type, unit, stock, min_threshold, notes, created_at, updated_at FROM ingredients`)
	if err != nil {
		return nil, fmt.Errorf("GetIngredients: %w", err)
	}
	defer rows.Close()

	var list []models.Ingredient
	for rows.Next() {
		var it models.Ingredient
		if err := rows.Scan(&it.ID, &it.Name, &it.Type, &it.Unit, &it.Stock, &it.MinThreshold, &it.Notes, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetIngredients scan: %w", err)
		}
		list = append(list, it)
	}
	return list, nil
}

// GetIngredient returns a single ingredient by id.
func (db *DB) GetIngredient(ctx context.Context, id int64) (*models.Ingredient, error) {
	row := db.QueryRow(ctx, `SELECT id, name, type, unit, stock, min_threshold, notes, created_at, updated_at FROM ingredients WHERE id = $1`, id)
	var it models.Ingredient
	if err := row.Scan(&it.ID, &it.Name, &it.Type, &it.Unit, &it.Stock, &it.MinThreshold, &it.Notes, &it.CreatedAt, &it.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetIngredient scan: %w", err)
	}
	return &it, nil
}

// CreateIngredient inserts a new ingredient and returns the new id.
func (db *DB) CreateIngredient(ctx context.Context, it *models.Ingredient) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO ingredients (name, type, unit, stock, min_threshold, notes) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		it.Name, it.Type, it.Unit, it.Stock, it.MinThreshold, it.Notes).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateIngredient: %w", err)
	}
	return id, nil
}

// UpdateIngredient updates editable fields of an ingredient.
func (db *DB) UpdateIngredient(ctx context.Context, it *models.Ingredient) error {
	_, err := db.Exec(ctx, `UPDATE ingredients SET name=$1, type=$2, unit=$3, stock=$4, min_threshold=$5, notes=$6, updated_at=NOW() WHERE id=$7`,
		it.Name, it.Type, it.Unit, it.Stock, it.MinThreshold, it.Notes, it.ID)
	if err != nil {
		return fmt.Errorf("UpdateIngredient: %w", err)
	}
	return nil
}
