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
		{ID: 1, Name: "Soy Wax", Type: "Wax", Unit: "kg", Stock: 50.0, MinThreshold: 10.0},
		{ID: 2, Name: "Lavender", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 3, Name: "Chamomile", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 4, Name: "Cedarwood", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 5, Name: "Ylang Ylang", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 6, Name: "Sweet Orange", Type: "EssentialOil", Unit: "mL", Stock: 1000.0, MinThreshold: 200.0},
		{ID: 7, Name: "Lemon", Type: "EssentialOil", Unit: "mL", Stock: 1000.0, MinThreshold: 200.0},
		{ID: 8, Name: "Bergamot", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 9, Name: "Spearmint", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 10, Name: "Vanilla Absolute", Type: "Natural Fragrance", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 11, Name: "Sandalwood", Type: "EssentialOil", Unit: "mL", Stock: 100.0, MinThreshold: 25.0},
		{ID: 12, Name: "Amyris", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 13, Name: "Peru Balsam", Type: "Natural Fragrance", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 14, Name: "Eucalyptus", Type: "EssentialOil", Unit: "mL", Stock: 750.0, MinThreshold: 150.0},
		{ID: 15, Name: "Rosemary", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 16, Name: "Cypress", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 17, Name: "Peppermint", Type: "EssentialOil", Unit: "mL", Stock: 750.0, MinThreshold: 150.0},
		{ID: 18, Name: "Rose Absolute", Type: "Natural Fragrance", Unit: "mL", Stock: 100.0, MinThreshold: 25.0},
		{ID: 19, Name: "Palmarosa", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 20, Name: "Petitgrain", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 21, Name: "Vetiver", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 22, Name: "Oakmoss", Type: "Natural Fragrance", Unit: "mL", Stock: 100.0, MinThreshold: 25.0},
		{ID: 23, Name: "Clary Sage", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 24, Name: "Geranium", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 25, Name: "Lemongrass", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 26, Name: "Lilac Natural Fragrance", Type: "Natural Fragrance", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 27, Name: "Coconut Natural Fragrance", Type: "Natural Fragrance", Unit: "mL", Stock: 750.0, MinThreshold: 150.0},
		{ID: 28, Name: "Lime", Type: "EssentialOil", Unit: "mL", Stock: 750.0, MinThreshold: 150.0},
		{ID: 29, Name: "Peach Natural Fragrance", Type: "Natural Fragrance", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 30, Name: "Pineapple Natural Fragrance", Type: "Natural Fragrance", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 31, Name: "Apple Natural Fragrance", Type: "Natural Fragrance", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 32, Name: "Cinnamon", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 33, Name: "Clove", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 34, Name: "Frankincense", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 35, Name: "Birch Tar", Type: "EssentialOil", Unit: "mL", Stock: 100.0, MinThreshold: 25.0},
		{ID: 36, Name: "Ginger", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 37, Name: "Nutmeg", Type: "EssentialOil", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 38, Name: "Cardamom", Type: "EssentialOil", Unit: "mL", Stock: 100.0, MinThreshold: 25.0},
		{ID: 39, Name: "Fir Balsam", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 40, Name: "Pine Needle", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 41, Name: "Cranberry Natural Fragrance", Type: "Natural Fragrance", Unit: "mL", Stock: 250.0, MinThreshold: 50.0},
		{ID: 42, Name: "Cocoa Absolute", Type: "Natural Fragrance", Unit: "mL", Stock: 100.0, MinThreshold: 25.0},
		{ID: 43, Name: "Patchouli", Type: "EssentialOil", Unit: "mL", Stock: 500.0, MinThreshold: 100.0},
		{ID: 44, Name: "Hay Absolute", Type: "Natural Fragrance", Unit: "mL", Stock: 100.0, MinThreshold: 25.0},
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
