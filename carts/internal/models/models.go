package models

import "time"

// Cart represents a user's shopping cart.
type Cart struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	SessionID string     `json:"session_id,omitempty"`
	Items     []CartItem `json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
}

// CartItem represents a single item in a cart.
type CartItem struct {
	ID         int64     `json:"id"`
	CartID     int64     `json:"cart_id"`
	ProductID  int64     `json:"product_id"`
	Quantity   int       `json:"quantity"`
	RunningLow bool      `json:"running_low,omitempty"`
	AddedAt    time.Time `json:"added_at,omitempty"`
}

// AddItemRequest is the JSON payload expected when adding an item to a cart.
type AddItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// UpdateItemRequest is the JSON payload for updating an existing cart item.
type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

// InventoryStockResponse is used when probing an inventory service for
// stock/threshold information for a product. This shape matches what the
// cart service's inventory probe expects by default.
type InventoryStockResponse struct {
	Stock     float64 `json:"stock"`
	Threshold float64 `json:"threshold"`
}

// NewCart creates a Cart with sensible zero values for easy construction.
func NewCart(userID int64) *Cart {
	return &Cart{
		UserID:    userID,
		Items:     []CartItem{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

// NewCartItem creates a CartItem with a positive default quantity of 1.
func NewCartItem(cartID, productID int64, quantity int) CartItem {
	if quantity <= 0 {
		quantity = 1
	}
	return CartItem{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
		AddedAt:   time.Time{},
	}
}

// Seed fixtures for development and tests. These are small, intentionally
// minimal examples to be used by DB seeding logic (if desired).
var SeedCarts = []Cart{
	{
		ID:     1,
		UserID: 1001,
		Items: []CartItem{
			{ID: 1, CartID: 1, ProductID: 1, Quantity: 2},
			{ID: 2, CartID: 1, ProductID: 2, Quantity: 1},
		},
	},
	{
		ID:     2,
		UserID: 1002,
		Items: []CartItem{
			{ID: 3, CartID: 2, ProductID: 2, Quantity: 3},
		},
	},
}
