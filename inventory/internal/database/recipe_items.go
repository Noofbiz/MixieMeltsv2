package database

import (
	"context"
	"fmt"
	"log"

	"com.MixieMelts.inventory/internal/models"
)

func (db *DB) seedRecipeItems(ctx context.Context) {
	// Only seed when the recipe table is empty.
	var count int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM recipe_items`).Scan(&count); err != nil {
		log.Printf("inventory: seed count query failed: %v", err)
		return
	}
	if count > 0 {
		log.Println("Ingredients table already seeded, skipping")
		return
	}

	recipeItems := []models.RecipeItem{
		{ProductID: 1, IngredientID: 1, Quantity: 10, Notes: "Vanilla for Mixie Melt"},
		{ProductID: 1, IngredientID: 4, Quantity: 200, Notes: "Milk for Mixie Melt"},
		{ProductID: 2, IngredientID: 2, Quantity: 15, Notes: "Chocolate for Choco Melt"},
		{ProductID: 2, IngredientID: 4, Quantity: 250, Notes: "Milk for Choco Melt"},
	}

	for _, ri := range recipeItems {
		query := `INSERT INTO recipe_items (product_id, ingredient_id, quantity, notes)
		VALUES ($1,$2,$3,$4)`
		_, err := db.Exec(ctx, query, ri.ProductID, ri.IngredientID, ri.Quantity, ri.Notes)
		if err != nil {
			log.Printf("inventory: seed insert product_id %d ingredient_id %d: %v", ri.ProductID, ri.IngredientID, err)
		}
	}
}

func (db *DB) createRecipeItemsTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS recipe_items (
		id SERIAL PRIMARY KEY,
		product_id BIGINT NOT NULL,
		ingredient_id BIGINT NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
		quantity DOUBLE PRECISION NOT NULL,
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(ctx, query)
	return err
}

// GetRecipe returns recipe items for a product.
func (db *DB) GetRecipe(ctx context.Context, productID int64) ([]models.RecipeItem, error) {
	rows, err := db.Query(ctx, `SELECT id, product_id, ingredient_id, quantity, notes, created_at, updated_at FROM recipe_items WHERE product_id = $1`, productID)
	if err != nil {
		return nil, fmt.Errorf("GetRecipe: %w", err)
	}
	defer rows.Close()

	var items []models.RecipeItem
	for rows.Next() {
		var r models.RecipeItem
		if err := rows.Scan(&r.ID, &r.ProductID, &r.IngredientID, &r.Quantity, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetRecipe scan: %w", err)
		}
		items = append(items, r)
	}
	return items, nil
}

// CreateRecipeItem adds an ingredient quantity requirement for a product.
func (db *DB) CreateRecipeItem(ctx context.Context, ri *models.RecipeItem) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO recipe_items (product_id, ingredient_id, quantity, notes) VALUES ($1,$2,$3,$4) RETURNING id`,
		ri.ProductID, ri.IngredientID, ri.Quantity, ri.Notes).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateRecipeItem: %w", err)
	}
	return id, nil
}
