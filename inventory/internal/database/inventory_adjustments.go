package database

import (
	"context"
	"fmt"
	"log"

	"com.MixieMelts.inventory/internal/models"
)

func (db *DB) seedInventoryAdjustments(ctx context.Context) {
	// Only seed when the adjustments table is empty.
	var count int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM inventory_adjustments`).Scan(&count); err != nil {
		log.Printf("inventory: seed count query failed: %v", err)
		return
	}
	if count > 0 {
		log.Println("Adjustments table already seeded, skipping")
		return
	}
	adjustments := []models.InventoryAdjustment{
		{IngredientID: 1, Change: 100, Reason: "Initial stock", Reference: "Initial stock for Vanilla", CreatedBy: "system"},
		{IngredientID: 2, Change: 150, Reason: "Initial stock", Reference: "Initial stock for Chocolate", CreatedBy: "system"},
		{IngredientID: 3, Change: 200, Reason: "Initial stock", Reference: "Initial stock for Maple Brown Sugar", CreatedBy: "system"},
	}

	for _, adj := range adjustments {
		query := `INSERT INTO inventory_adjustments (ingredient_id, change, reason, reference, created_by)
		VALUES ($1,$2,$3,$4,$5)`
		_, err := db.Exec(ctx, query, adj.IngredientID, adj.Change, adj.Reason, adj.Reference, adj.CreatedBy)
		if err != nil {
			log.Printf("inventory: seed insert ingredient_id %d: %v", adj.IngredientID, err)
		}
	}
}

func (db *DB) createInventoryAdjustmentsTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS inventory_adjustments (
		id SERIAL PRIMARY KEY,
		ingredient_id BIGINT NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
		change DOUBLE PRECISION NOT NULL,
		reason TEXT,
		reference TEXT,
		created_by TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(ctx, query)
	return err
}

// AdjustIngredientStock performs a stock adjustment and records it in inventory_adjustments.
// change may be positive (restock) or negative (usage/waste). The operation is transactional.
func (db *DB) AdjustIngredientStock(ctx context.Context, ingredientID int64, change float64, reason, reference, createdBy string) (float64, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("AdjustIngredientStock begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // safe to call
	}()

	// Update stock
	var newStock float64
	err = tx.QueryRow(ctx, `UPDATE ingredients SET stock = stock + $1, updated_at = NOW() WHERE id = $2 RETURNING stock`, change, ingredientID).Scan(&newStock)
	if err != nil {
		return 0, fmt.Errorf("AdjustIngredientStock update: %w", err)
	}

	// Insert adjustment record
	_, err = tx.Exec(ctx, `INSERT INTO inventory_adjustments (ingredient_id, change, reason, reference, created_by, created_at) VALUES ($1,$2,$3,$4,$5,NOW())`,
		ingredientID, change, reason, reference, createdBy)
	if err != nil {
		return 0, fmt.Errorf("AdjustIngredientStock insert adjustment: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("AdjustIngredientStock commit: %w", err)
	}

	return newStock, nil
}
