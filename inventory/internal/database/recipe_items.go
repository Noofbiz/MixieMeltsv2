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
		log.Println("Recipe items table already seeded, skipping")
		return
	}

	// seed entries that mirror the products service recipes (minimal example)
	recipeItems := []models.RecipeItem{
		// ProductID: 1 (Serene Sanctuary)
		{ProductID: 1, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 1, IngredientID: 2, Unit: "mL", Amount: 4.0},  // Lavender
		{ProductID: 1, IngredientID: 3, Unit: "mL", Amount: 3.0},  // Chamomile
		{ProductID: 1, IngredientID: 4, Unit: "mL", Amount: 2.0},  // Cedarwood
		{ProductID: 1, IngredientID: 5, Unit: "mL", Amount: 1.0},  // Ylang Ylang

		// ProductID: 2 (Citrus Sunshine)
		{ProductID: 2, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 2, IngredientID: 6, Unit: "mL", Amount: 3.5},  // Sweet Orange
		{ProductID: 2, IngredientID: 7, Unit: "mL", Amount: 3.0},  // Lemon
		{ProductID: 2, IngredientID: 8, Unit: "mL", Amount: 2.5},  // Bergamot
		{ProductID: 2, IngredientID: 9, Unit: "mL", Amount: 1.0},  // Spearmint

		// ProductID: 3 (Cozy Cashmere)
		{ProductID: 3, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 3, IngredientID: 10, Unit: "mL", Amount: 4.0}, // Vanilla Absolute
		{ProductID: 3, IngredientID: 11, Unit: "mL", Amount: 3.0}, // Sandalwood
		{ProductID: 3, IngredientID: 12, Unit: "mL", Amount: 2.0}, // Amyris
		{ProductID: 3, IngredientID: 13, Unit: "mL", Amount: 1.5}, // Peru Balsam

		// ProductID: 4 (Woodland Walk)
		{ProductID: 4, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 4, IngredientID: 14, Unit: "mL", Amount: 3.5}, // Eucalyptus
		{ProductID: 4, IngredientID: 16, Unit: "mL", Amount: 3.0}, // Cypress
		{ProductID: 4, IngredientID: 15, Unit: "mL", Amount: 2.0}, // Rosemary
		{ProductID: 4, IngredientID: 17, Unit: "mL", Amount: 1.5}, // Peppermint

		// ProductID: 5 (Rose Garden)
		{ProductID: 5, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 5, IngredientID: 18, Unit: "mL", Amount: 5.0}, // Rose Absolute
		{ProductID: 5, IngredientID: 19, Unit: "mL", Amount: 3.0}, // Palmarosa
		{ProductID: 5, IngredientID: 20, Unit: "mL", Amount: 2.0}, // Petitgrain

		// ProductID: 6 (April Showers)
		{ProductID: 6, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 6, IngredientID: 21, Unit: "mL", Amount: 3.5}, // Vetiver
		{ProductID: 6, IngredientID: 22, Unit: "mL", Amount: 3.0}, // Oakmoss
		{ProductID: 6, IngredientID: 20, Unit: "mL", Amount: 2.0}, // Petitgrain
		{ProductID: 6, IngredientID: 23, Unit: "mL", Amount: 1.5}, // Clary Sage

		// ProductID: 7 (Wildflower Meadow)
		{ProductID: 7, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 7, IngredientID: 24, Unit: "mL", Amount: 3.5}, // Geranium
		{ProductID: 7, IngredientID: 2, Unit: "mL", Amount: 3.0},  // Lavender
		{ProductID: 7, IngredientID: 3, Unit: "mL", Amount: 2.0},  // Chamomile
		{ProductID: 7, IngredientID: 25, Unit: "mL", Amount: 1.5}, // Lemongrass

		// ProductID: 8 (Lilac Bloom)
		{ProductID: 8, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 8, IngredientID: 26, Unit: "mL", Amount: 9.0}, // Lilac Natural Fragrance
		{ProductID: 8, IngredientID: 5, Unit: "mL", Amount: 1.5},  // Ylang Ylang

		// ProductID: 9 (Coastal Breeze)
		{ProductID: 9, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 9, IngredientID: 27, Unit: "mL", Amount: 4.0}, // Coconut Natural Fragrance
		{ProductID: 9, IngredientID: 28, Unit: "mL", Amount: 3.5}, // Lime
		{ProductID: 9, IngredientID: 9, Unit: "mL", Amount: 1.5},  // Spearmint
		{ProductID: 9, IngredientID: 12, Unit: "mL", Amount: 1.0}, // Amyris

		// ProductID: 10 (Sun-Kissed Peach)
		{ProductID: 10, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 10, IngredientID: 29, Unit: "mL", Amount: 6.0}, // Peach Natural Fragrance
		{ProductID: 10, IngredientID: 6, Unit: "mL", Amount: 2.5},  // Sweet Orange
		{ProductID: 10, IngredientID: 10, Unit: "mL", Amount: 2.0}, // Vanilla Absolute

		// ProductID: 11 (Tropical Getaway)
		{ProductID: 11, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 11, IngredientID: 30, Unit: "mL", Amount: 4.0}, // Pineapple Natural Fragrance
		{ProductID: 11, IngredientID: 27, Unit: "mL", Amount: 3.5}, // Coconut Natural Fragrance
		{ProductID: 11, IngredientID: 28, Unit: "mL", Amount: 1.5}, // Lime
		{ProductID: 11, IngredientID: 10, Unit: "mL", Amount: 1.0}, // Vanilla Absolute

		// ProductID: 12 (Autumn Harvest)
		{ProductID: 12, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 12, IngredientID: 31, Unit: "mL", Amount: 6.0}, // Apple Natural Fragrance
		{ProductID: 12, IngredientID: 32, Unit: "mL", Amount: 2.5}, // Cinnamon
		{ProductID: 12, IngredientID: 33, Unit: "mL", Amount: 1.5}, // Clove

		// ProductID: 13 (Bonfire Flannel)
		{ProductID: 13, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 13, IngredientID: 4, Unit: "mL", Amount: 4.0},  // Cedarwood
		{ProductID: 13, IngredientID: 34, Unit: "mL", Amount: 2.5}, // Frankincense
		{ProductID: 13, IngredientID: 21, Unit: "mL", Amount: 2.0}, // Vetiver
		{ProductID: 13, IngredientID: 35, Unit: "mL", Amount: 0.5}, // Birch Tar

		// ProductID: 14 (Pumpkin Spice)
		{ProductID: 14, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 14, IngredientID: 32, Unit: "mL", Amount: 3.0}, // Cinnamon
		{ProductID: 14, IngredientID: 33, Unit: "mL", Amount: 2.0}, // Clove
		{ProductID: 14, IngredientID: 36, Unit: "mL", Amount: 2.0}, // Ginger
		{ProductID: 14, IngredientID: 37, Unit: "mL", Amount: 1.5}, // Nutmeg
		{ProductID: 14, IngredientID: 38, Unit: "mL", Amount: 1.0}, // Cardamom

		// ProductID: 15 (Winter Woods)
		{ProductID: 15, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 15, IngredientID: 39, Unit: "mL", Amount: 4.0}, // Fir Balsam
		{ProductID: 15, IngredientID: 40, Unit: "mL", Amount: 3.0}, // Pine Needle
		{ProductID: 15, IngredientID: 16, Unit: "mL", Amount: 2.0}, // Cypress
		{ProductID: 15, IngredientID: 4, Unit: "mL", Amount: 1.0},  // Cedarwood

		// ProductID: 16 (Spiced Cranberry)
		{ProductID: 16, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 16, IngredientID: 41, Unit: "mL", Amount: 6.0}, // Cranberry Natural Fragrance
		{ProductID: 16, IngredientID: 6, Unit: "mL", Amount: 3.0},  // Sweet Orange
		{ProductID: 16, IngredientID: 32, Unit: "mL", Amount: 1.5}, // Cinnamon

		// ProductID: 17 (Peppermint Cocoa)
		{ProductID: 17, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 17, IngredientID: 17, Unit: "mL", Amount: 4.5}, // Peppermint
		{ProductID: 17, IngredientID: 42, Unit: "mL", Amount: 3.5}, // Cocoa Absolute
		{ProductID: 17, IngredientID: 10, Unit: "mL", Amount: 2.5}, // Vanilla Absolute

		// ProductID: 18 (Witches' Brew)
		{ProductID: 18, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 18, IngredientID: 43, Unit: "mL", Amount: 4.0}, // Patchouli
		{ProductID: 18, IngredientID: 34, Unit: "mL", Amount: 3.0}, // Frankincense
		{ProductID: 18, IngredientID: 32, Unit: "mL", Amount: 2.0}, // Cinnamon
		{ProductID: 18, IngredientID: 33, Unit: "mL", Amount: 1.0}, // Clove

		// ProductID: 19 (Haunted Hayride)
		{ProductID: 19, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 19, IngredientID: 44, Unit: "mL", Amount: 3.5}, // Hay Absolute
		{ProductID: 19, IngredientID: 21, Unit: "mL", Amount: 3.0}, // Vetiver
		{ProductID: 19, IngredientID: 12, Unit: "mL", Amount: 2.0}, // Amyris
		{ProductID: 19, IngredientID: 4, Unit: "mL", Amount: 1.5},  // Cedarwood

		// ProductID: 20 (Christmas Tree)
		{ProductID: 20, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 20, IngredientID: 39, Unit: "mL", Amount: 5.0}, // Fir Balsam
		{ProductID: 20, IngredientID: 40, Unit: "mL", Amount: 3.5}, // Pine Needle
		{ProductID: 20, IngredientID: 6, Unit: "mL", Amount: 1.5},  // Sweet Orange

		// ProductID: 21 (Gingerbread House)
		{ProductID: 21, IngredientID: 1, Unit: "g", Amount: 100.0}, // Soy Wax
		{ProductID: 21, IngredientID: 36, Unit: "mL", Amount: 3.0}, // Ginger
		{ProductID: 21, IngredientID: 32, Unit: "mL", Amount: 2.5}, // Cinnamon
		{ProductID: 21, IngredientID: 10, Unit: "mL", Amount: 2.0}, // Vanilla Absolute
		{ProductID: 21, IngredientID: 37, Unit: "mL", Amount: 1.5}, // Nutmeg
		{ProductID: 21, IngredientID: 33, Unit: "mL", Amount: 1.0}, // Clove
	}

	for _, ri := range recipeItems {
		query := `INSERT INTO recipe_items (product_id, ingredient_id, unit, amount, notes)
		VALUES ($1,$2,$3,$4,$5)`
		_, err := db.Exec(ctx, query, ri.ProductID, ri.IngredientID, ri.Unit, ri.Amount, ri.Notes)
		if err != nil {
			log.Printf("inventory: seed insert product_id %d: %v", ri.ProductID, err)
		}
	}
}

func (db *DB) createRecipeItemsTable(ctx context.Context) error {
	// Ensure the canonical table exists with the new schema.
	create := `CREATE TABLE IF NOT EXISTS recipe_items (
		id SERIAL PRIMARY KEY,
		product_id BIGINT NOT NULL,
		ingredient_id BIGINT REFERENCES ingredients(id) ON DELETE RESTRICT,
		unit VARCHAR(50),
		amount DOUBLE PRECISION,
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(ctx, create); err != nil {
		return err
	}

	return nil
}

// GetRecipe returns recipe items for a product.
func (db *DB) GetRecipe(ctx context.Context, productID int64) ([]models.RecipeItem, error) {
	rows, err := db.Query(ctx, `SELECT id, product_id, ingredient_id, unit, amount, notes, created_at, updated_at FROM recipe_items WHERE product_id = $1`, productID)
	if err != nil {
		return nil, fmt.Errorf("GetRecipe: %w", err)
	}
	defer rows.Close()

	var items []models.RecipeItem
	for rows.Next() {
		var r models.RecipeItem
		if err := rows.Scan(&r.ID, &r.ProductID, &r.IngredientID, &r.Unit, &r.Amount, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetRecipe scan: %w", err)
		}
		items = append(items, r)
	}
	return items, nil
}

// CreateRecipeItem adds an ingredient entry for a product.
func (db *DB) CreateRecipeItem(ctx context.Context, ri *models.RecipeItem) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO recipe_items (product_id, ingredient_id, unit, amount, notes) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		ri.ProductID, ri.IngredientID, ri.Unit, ri.Amount, ri.Notes).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateRecipeItem: %w", err)
	}
	return id, nil
}
