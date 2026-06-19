package db

import (
	"context"
	"fmt"

	"com.MixieMeltsv2/carts/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a thin wrapper around a pgx connection pool used by the carts service.
type DB struct {
	pool *pgxpool.Pool
}

// New creates a new DB wrapper. Schema creation and seeding are managed by SQL
// migrations (use golang-migrate). The previous in-code table/index/seed
// creation has been removed in favor of migrations.
func New(ctx context.Context, dbURL string) (*DB, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	d := &DB{pool: pool}
	return d, nil
}

// Close closes the underlying pool.
func (d *DB) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

// createTables removed: schema creation is handled by SQL migrations in ./migrations.
// This function intentionally left out to avoid duplicating schema management.
// See migrations/000001_create_carts_and_items.up.sql for the canonical schema.

// createIndices removed: index creation is managed by SQL migrations.
// See migrations/000001_create_carts_and_items.up.sql for the canonical indexes.

// seedIfEmpty removed: seeding is handled by migrations where appropriate.
// The migration `000001_create_carts_and_items.up.sql` contains a guarded
// seed that only runs when the carts table is empty (for local development).

// GetOrCreateCart returns a cart for the given user, creating one if it does not exist.
func (d *DB) GetOrCreateCart(ctx context.Context, userID int64) (*models.Cart, error) {
	var c models.Cart
	row := d.pool.QueryRow(ctx, `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=$1 LIMIT 1`, userID)
	if err := row.Scan(&c.ID, &c.UserID, &c.CreatedAt, &c.UpdatedAt); err == nil {
		// found existing cart; load items
		items, err := d.GetCartItems(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Items = items
		return &c, nil
	} else if err != nil && err != pgx.ErrNoRows {
		// unexpected DB error
		return nil, fmt.Errorf("GetOrCreateCart select: %w", err)
	}

	// not found -> create
	if err := d.pool.QueryRow(ctx, `INSERT INTO carts (user_id) VALUES ($1) RETURNING id, created_at, updated_at`, userID).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, fmt.Errorf("GetOrCreateCart insert: %w", err)
	}
	c.UserID = userID
	c.Items = []models.CartItem{}
	return &c, nil
}

// GetCartItems returns all items for a cart.
func (d *DB) GetCartItems(ctx context.Context, cartID int64) ([]models.CartItem, error) {
	rows, err := d.pool.Query(ctx, `SELECT id, cart_id, product_id, quantity, added_at FROM cart_items WHERE cart_id=$1`, cartID)
	if err != nil {
		return nil, fmt.Errorf("GetCartItems query: %w", err)
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var it models.CartItem
		if err := rows.Scan(&it.ID, &it.CartID, &it.ProductID, &it.Quantity, &it.AddedAt); err != nil {
			return nil, fmt.Errorf("GetCartItems scan: %w", err)
		}
		items = append(items, it)
	}
	return items, nil
}

// AddItem inserts or increments an item in a cart using an upsert.
func (d *DB) AddItem(ctx context.Context, cartID int64, productID int64, qty int) (*models.CartItem, error) {
	var it models.CartItem
	query := `INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1,$2,$3)
ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
RETURNING id, cart_id, product_id, quantity, added_at`
	if err := d.pool.QueryRow(ctx, query, cartID, productID, qty).Scan(&it.ID, &it.CartID, &it.ProductID, &it.Quantity, &it.AddedAt); err != nil {
		return nil, fmt.Errorf("AddItem upsert: %w", err)
	}
	return &it, nil
}

// RemoveItem deletes an item by id.
func (d *DB) RemoveItem(ctx context.Context, itemID int64) error {
	_, err := d.pool.Exec(ctx, `DELETE FROM cart_items WHERE id=$1`, itemID)
	return err
}

// UpdateItemQuantity sets the quantity for an item and returns the updated row.
func (d *DB) UpdateItemQuantity(ctx context.Context, itemID int64, qty int) (*models.CartItem, error) {
	_, err := d.pool.Exec(ctx, `UPDATE cart_items SET quantity=$1 WHERE id=$2`, qty, itemID)
	if err != nil {
		return nil, fmt.Errorf("UpdateItemQuantity exec: %w", err)
	}
	var it models.CartItem
	if err := d.pool.QueryRow(ctx, `SELECT id, cart_id, product_id, quantity, added_at FROM cart_items WHERE id=$1`, itemID).
		Scan(&it.ID, &it.CartID, &it.ProductID, &it.Quantity, &it.AddedAt); err != nil {
		return nil, fmt.Errorf("UpdateItemQuantity select: %w", err)
	}
	return &it, nil
}

// CountCartsWithProduct returns how many distinct carts contain the given product.
func (d *DB) CountCartsWithProduct(ctx context.Context, productID int64) (int, error) {
	var count int
	if err := d.pool.QueryRow(ctx, `SELECT COUNT(DISTINCT cart_id) FROM cart_items WHERE product_id=$1`, productID).Scan(&count); err != nil {
		return 0, fmt.Errorf("CountCartsWithProduct: %w", err)
	}
	return count, nil
}

// GetOrCreateCartBySession returns a cart associated with a session_id, creating one if missing.
// sessionID should be a stable opaque identifier stored in the user's cookie.
func (d *DB) GetOrCreateCartBySession(ctx context.Context, sessionID string) (*models.Cart, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session id required")
	}
	var c models.Cart
	row := d.pool.QueryRow(ctx, `SELECT id, user_id, session_id, created_at, updated_at FROM carts WHERE session_id=$1 LIMIT 1`, sessionID)
	if err := row.Scan(&c.ID, &c.UserID, &c.SessionID, &c.CreatedAt, &c.UpdatedAt); err == nil {
		items, err := d.GetCartItems(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Items = items
		return &c, nil
	}

	// create new cart with session_id
	if err := d.pool.QueryRow(ctx, `INSERT INTO carts (session_id) VALUES ($1) RETURNING id, created_at, updated_at`, sessionID).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, fmt.Errorf("GetOrCreateCartBySession insert: %w", err)
	}
	c.SessionID = sessionID
	c.Items = []models.CartItem{}
	return &c, nil
}

// MergeSessionCartIntoUserCart transfers items from a session cart into a user's cart and removes the session cart.
// It performs upserts so quantities are summed.
func (d *DB) MergeSessionCartIntoUserCart(ctx context.Context, sessionID string, userID int64) error {
	if sessionID == "" {
		return nil
	}
	// find session cart
	var sessionCartID int64
	err := d.pool.QueryRow(ctx, `SELECT id FROM carts WHERE session_id=$1 LIMIT 1`, sessionID).Scan(&sessionCartID)
	if err != nil {
		// no session cart to merge or DB error - treat missing as no-op
		return nil
	}

	// get or create user's cart
	userCart, err := d.GetOrCreateCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("merge: failed to get or create user cart: %w", err)
	}

	// get session cart items
	items, err := d.GetCartItems(ctx, sessionCartID)
	if err != nil {
		return fmt.Errorf("merge: failed to read session cart items: %w", err)
	}

	// upsert each item into the user's cart
	for _, it := range items {
		if _, err := d.AddItem(ctx, userCart.ID, it.ProductID, it.Quantity); err != nil {
			return fmt.Errorf("merge: failed to upsert item pid=%d: %w", it.ProductID, err)
		}
	}

	// remove the session cart (cascade deletes items)
	if _, err := d.pool.Exec(ctx, `DELETE FROM carts WHERE id=$1`, sessionCartID); err != nil {
		return fmt.Errorf("merge: failed to delete session cart: %w", err)
	}

	return nil
}
