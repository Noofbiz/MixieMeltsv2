package database

import (
	"context"
	"fmt"
	"log"

	"com.MixieMelts.products/internal/models"
	"github.com/jackc/pgx/v5"
)

// DB represents the database connection.
type DB struct {
	*pgx.Conn
}

// New creates a new database connection.
func New(config string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database.")

	dbWrapper := &DB{conn}
	if err := dbWrapper.createTables(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	dbWrapper.Seed(context.Background())

	return dbWrapper, nil
}

func (db *DB) createTables(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		category VARCHAR(255) NOT NULL,
		scent VARCHAR(255) NOT NULL,
		price NUMERIC(10, 2) NOT NULL,
		subscription BOOLEAN DEFAULT false,
		image VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(ctx, query)
	return err
}

// GetProducts retrieves all products from the database.
func (db *DB) GetProducts(ctx context.Context, limit int) ([]models.Product, error) {
	query := "SELECT id, name, category, scent, price, subscription, image, description, created_at, updated_at FROM products"
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Scent, &product.Price, &product.Subscription, &product.Image, &product.Description, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// CreateProduct inserts a new product into the database.
func (db *DB) CreateProduct(ctx context.Context, product *models.Product) (int64, error) {
	query := `
	INSERT INTO products (name, category, scent, price, subscription, image, description)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id;
	`
	var productID int64
	err := db.QueryRow(ctx, query, product.Name, product.Category, product.Scent, product.Price, product.Subscription, product.Image, product.Description).Scan(&productID)
	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}
	return productID, nil
}
