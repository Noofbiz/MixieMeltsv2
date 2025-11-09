package database

import (
	"context"
	"fmt"
	"log"

	"com.MixieMelts.products/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents the database connection.
type DB struct {
	*pgxpool.Pool
}

// New creates a new database connection.
func New(config string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	log.Println("Successfully created database connection pool.")

	dbWrapper := &DB{pool}
	if err := dbWrapper.createTables(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	if err := dbWrapper.createSubscriptionTables(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create subscription tables: %w", err)
	}

	dbWrapper.Seed(context.Background())

	return dbWrapper, nil
}

func (db *DB) createTables(ctx context.Context) error {
	// Create products + recipe_items canonical schema if missing.
	// Keep the CREATE statement minimal and idempotent.
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
	if _, err := db.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}

func (db *DB) createSubscriptionTables(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS subscription_boxes (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		price NUMERIC(10, 2) NOT NULL,
		image VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(ctx, query)
	return err
}

// GetProducts retrieves all products from the database, including their recipe items.
func (db *DB) GetProducts(ctx context.Context, limit int, prodID int64) ([]models.Product, error) {
	query := "SELECT id, name, category, scent, price, subscription, image, description, created_at, updated_at FROM products"
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}
	if prodID > 0 {
		query = fmt.Sprintf("%s WHERE id = %d", query, prodID)
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

		// Load recipe items for this product
		recipeQuery := `SELECT ingredient_id, unit, amount, notes, created_at, updated_at FROM recipe_items WHERE product_id = $1 ORDER BY id`
		rrows, err := db.Query(ctx, recipeQuery, product.ID)
		defer rrows.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to get recipe items for product %d: %w", product.ID, err)
		}
		var recipe []models.Ingredient
		for rrows.Next() {
			var ing models.Ingredient
			var typ string
			if err := rrows.Scan(&ing.ID, &ing.Unit, &ing.Amount, &ing.Notes, &ing.CreatedAt, &ing.UpdatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan recipe item for product %d: %w", product.ID, err)
			}
			ingredientQuery := `SELECT name, type FROM ingredients WHERE id = $1`
			ingRow := db.QueryRow(ctx, ingredientQuery, ing.ID)
			if err = ingRow.Scan(&ing.Name, &typ); err != nil {
				return nil, fmt.Errorf("failed to scan ingredient name for ingredient %d: %w", ing.ID, err)
			}
			ing.Type = models.IngredientType(typ)
			recipe = append(recipe, ing)
		}
		if len(recipe) == 0 {
			// Ensure we return an empty array instead of null when marshalled to JSON.
			recipe = []models.Ingredient{}
		}
		product.Recipe = recipe

		products = append(products, product)
	}

	return products, nil
}

// CreateProduct inserts a new product into the database and creates associated recipe items.
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

	// If recipe items were provided, insert them into recipe_items table.
	if len(product.Recipe) > 0 {
		insertQuery := `
		INSERT INTO recipe_items (product_id, name, type, unit, amount, notes)
		VALUES ($1, $2, $3, $4, $5, $6);
		`
		for _, ing := range product.Recipe {
			_, err := db.Exec(ctx, insertQuery, productID, ing.Name, string(ing.Type), ing.Unit, ing.Amount, ing.Notes)
			if err != nil {
				// Log and continue inserting remaining items - prefer partial success to failing entire request.
				log.Printf("failed to insert recipe item %q for product %d: %v", ing.Name, productID, err)
			}
		}
	}

	return productID, nil
}

// GetProduct returns a single product by id including its recipe items.
func (db *DB) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	prods, err := db.GetProducts(ctx, 0, id)
	return &prods[0], err
}

// CreateProductTx creates a product and its recipe items in a single transaction.
func (db *DB) CreateProductTx(ctx context.Context, product *models.Product) (int64, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("CreateProductTx begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var productID int64
	insertProduct := `
 	INSERT INTO products (name, category, scent, price, subscription, image, description)
 	VALUES ($1, $2, $3, $4, $5, $6, $7)
 	RETURNING id;
 	`
	if err := tx.QueryRow(ctx, insertProduct, product.Name, product.Category, product.Scent, product.Price, product.Subscription, product.Image, product.Description).Scan(&productID); err != nil {
		return 0, fmt.Errorf("CreateProductTx insert product: %w", err)
	}

	// Insert provided recipe items under the same transaction.
	if len(product.Recipe) > 0 {
		insertItem := `INSERT INTO recipe_items (product_id, name, type, unit, amount, notes) VALUES ($1, $2, $3, $4, $5, $6);`
		for _, ing := range product.Recipe {
			if _, err := tx.Exec(ctx, insertItem, productID, ing.Name, string(ing.Type), ing.Unit, ing.Amount, ing.Notes); err != nil {
				return 0, fmt.Errorf("CreateProductTx insert recipe item %q: %w", ing.Name, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("CreateProductTx commit: %w", err)
	}
	return productID, nil
}

func (db *DB) GetSubscriptionBoxes(ctx context.Context, limit int) ([]models.SubscriptionBox, error) {
	query := "SELECT id, name, description, price, image, created_at, updated_at FROM subscription_boxes"
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription boxes: %w", err)
	}
	defer rows.Close()

	var subscriptionBoxes []models.SubscriptionBox
	for rows.Next() {
		var subscriptionBox models.SubscriptionBox
		if err := rows.Scan(&subscriptionBox.ID, &subscriptionBox.Name, &subscriptionBox.Description, &subscriptionBox.Price, &subscriptionBox.Image, &subscriptionBox.CreatedAt, &subscriptionBox.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription box: %w", err)
		}
		subscriptionBoxes = append(subscriptionBoxes, subscriptionBox)
	}
	return subscriptionBoxes, nil
}

// CreateSubscriptionBox inserts a new subscription box into the database.
func (db *DB) CreateSubscriptionBox(ctx context.Context, subscriptionBox *models.SubscriptionBox) (int64, error) {
	query := `
	INSERT INTO subscription_boxes (name, description, price, image)
	VALUES ($1, $2, $3)
	RETURNING id;
	`
	var subscriptionBoxID int64
	err := db.QueryRow(ctx, query, subscriptionBox.Name, subscriptionBox.Price, subscriptionBox.Description, subscriptionBox.Image).Scan(&subscriptionBoxID)
	if err != nil {
		return 0, fmt.Errorf("failed to create subscription box: %w", err)
	}

	return subscriptionBoxID, nil
}
