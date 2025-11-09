package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"com.MixieMelts.products/internal/models"
	"github.com/go-chi/chi/v5"
)

// Handler represents the HTTP handlers for the service.
type DBLayer interface {
	// GetProducts returns multiple products (optionally limited).
	GetProducts(ctx context.Context, limit int, prodID int64) ([]models.Product, error)

	// GetProduct returns a single product by id including its recipe.
	GetProduct(ctx context.Context, id int64) (*models.Product, error)

	// CreateProduct inserts a product without guaranteeing transactional recipe inserts.
	CreateProduct(ctx context.Context, product *models.Product) (int64, error)

	// CreateProductTx inserts a product and its recipe atomically in a transaction.
	CreateProductTx(ctx context.Context, product *models.Product) (int64, error)

	GetSubscriptionBoxes(ctx context.Context, limit int) ([]models.SubscriptionBox, error)
	CreateSubscriptionBox(ctx context.Context, box *models.SubscriptionBox) (int64, error)
}

type Handler struct {
	db DBLayer
}

// New creates a new handler.
func New(db DBLayer) *Handler {
	return &Handler{db: db}
}

// GetProducts handles GET requests to /products.
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}
	products, err := h.db.GetProducts(r.Context(), limit, 0)
	if err != nil {
		println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to get products")
		return
	}
	// Ensure recipe is a non-nil slice for the frontend (avoid JSON `null`).
	for i := range products {
		if products[i].Recipe == nil {
			products[i].Recipe = []models.Ingredient{}
		}
	}
	log.Printf("handlers: returning %d products (limit=%d)", len(products), limit)
	respondWithJSON(w, http.StatusOK, products)
}

// GetProduct handles GET requests to /products/{id}.
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	// If chi did not populate the URL param (tests often create a request
	// directly with a path like "/products/42"), attempt to extract the id
	// from the request path as a fallback.
	if idStr == "" {
		trimmed := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(trimmed, "/")
		if len(parts) > 0 {
			idStr = parts[len(parts)-1]
		}
	}

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "Missing product id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	product, err := h.db.GetProduct(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get product")
		return
	}
	if product == nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	// Ensure recipe is a non-nil slice to prevent JSON null in the frontend.
	if product.Recipe == nil {
		product.Recipe = []models.Ingredient{}
	}
	log.Printf("handlers: returning product id=%d with %d recipe items", product.ID, len(product.Recipe))
	respondWithJSON(w, http.StatusOK, product)
}

// CreateProduct handles POST requests to /products.
// This handler expects the database layer to provide a transactional create
// that inserts the product and its recipe atomically.
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Prefer transactional creation to ensure product + recipe are inserted atomically.
	productID, err := h.db.CreateProductTx(r.Context(), &product)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	// After creation, retrieve the product as stored in the DB so we return the
	// authoritative representation (including any recipe items inserted by the DB).
	created, err := h.db.GetProduct(r.Context(), productID)
	if err != nil {
		// Treat inability to retrieve the created product as an internal server error.
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve created product")
		return
	}
	if created == nil {
		respondWithError(w, http.StatusBadRequest, "Created product not found")
		return
	}

	// Ensure recipe is non-nil for JSON consumers.
	if created.Recipe == nil {
		created.Recipe = []models.Ingredient{}
	}

	respondWithJSON(w, http.StatusCreated, created)
}

// GetSubscriptionBoxes handles GET requests to /subscription-boxes.
func (h *Handler) GetSubscriptionBoxes(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	subscriptionBoxes, err := h.db.GetSubscriptionBoxes(r.Context(), limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get subscription boxes")
		return
	}
	respondWithJSON(w, http.StatusOK, subscriptionBoxes)
}

// CreateSubscriptionBox handles POST requests to /subscription-boxes.
func (h *Handler) CreateSubscriptionBox(w http.ResponseWriter, r *http.Request) {
	var subscriptionBox models.SubscriptionBox
	if err := json.NewDecoder(r.Body).Decode(&subscriptionBox); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	subscriptionBoxID, err := h.db.CreateSubscriptionBox(r.Context(), &subscriptionBox)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create subscription box")
		return
	}

	subscriptionBox.ID = subscriptionBoxID

	respondWithJSON(w, http.StatusCreated, subscriptionBox)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
