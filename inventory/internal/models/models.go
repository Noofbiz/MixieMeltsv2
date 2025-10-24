package models

import "time"

// IngredientType represents the category of an ingredient/material.
type IngredientType string

const (
	IngredientTypeWax   IngredientType = "wax"
	IngredientTypeBase  IngredientType = "base"
	IngredientTypeScent IngredientType = "scent"
	IngredientTypeOther IngredientType = "other"
)

// Ingredient represents a raw material used to make products (wax, scent base, additive, etc).
type Ingredient struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Type         IngredientType `json:"type"`
	Unit         string         `json:"unit"`                    // e.g. "g", "kg", "ml", "L"
	Stock        float64        `json:"stock"`                   // current on-hand quantity in Unit
	MinThreshold float64        `json:"min_threshold,omitempty"` // optional reorder threshold
	Notes        string         `json:"notes,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// InventoryAdjustment records a change to an Ingredient's stock level.
type InventoryAdjustment struct {
	ID           int64     `json:"id"`
	IngredientID int64     `json:"ingredient_id"`
	Change       float64   `json:"change"`               // positive for addition, negative for subtraction
	Reason       string    `json:"reason,omitempty"`     // e.g. "restock", "sale", "waste", "correction"
	Reference    string    `json:"reference,omitempty"`  // optional external reference (PO number, order id, etc)
	CreatedBy    string    `json:"created_by,omitempty"` // who made the adjustment
	CreatedAt    time.Time `json:"created_at"`
}

// RecipeItem represents the amount of a single ingredient required to produce one unit
// of a product (e.g., grams of wax, ml of scent per melt).
type RecipeItem struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`      // references products service product id
	IngredientID int64     `json:"ingredient_id"`   // references Ingredient.ID
	Quantity     float64   `json:"quantity"`        // amount of ingredient per unit product (in Ingredient.Unit)
	Notes        string    `json:"notes,omitempty"` // optional notes for the recipe step
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProductStockProjection is a helper model describing how many product units can be
// produced based on current ingredient stocks and a recipe. Useful for admin UI.
type ProductStockProjection struct {
	ProductID      int64     `json:"product_id"`
	AvailableUnits float64   `json:"available_units"`  // fractional units allowed (e.g., 12.5 melts)
	LimitingItemID int64     `json:"limiting_item_id"` // ingredient that limits production
	LimitingAvail  float64   `json:"limiting_avail"`   // how many units that limiting ingredient permits
	ComputedAt     time.Time `json:"computed_at"`
}

// Lightweight view types for API responses

// IngredientSummary returns a compact representation of ingredient and its stock.
type IngredientSummary struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Unit      string  `json:"unit"`
	Stock     float64 `json:"stock"`
	Threshold float64 `json:"min_threshold,omitempty"`
}

// Subscription: helper to set zero values for time in fixtures, if needed.
func ZeroTime() time.Time {
	return time.Time{}
}
