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
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Type      IngredientType `json:"type"`
	Unit      string         `json:"unit"`   // e.g. "g", "kg", "ml", "L"
	Amount    float64        `json:"amount"` // current on-hand quantity in Unit
	Notes     string         `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
