package models

import "time"

// Product represents a product in the system.
type Product struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	Category     string       `json:"category"`
	Scent        string       `json:"scent"`
	Price        float64      `json:"price"`
	Subscription bool         `json:"subscription"`
	Image        string       `json:"image"`
	Recipe       []Ingredient `json:"recipe"`
	Description  string       `json:"description"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
