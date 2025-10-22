package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"com.MixieMelts.products/internal/database"
	"com.MixieMelts.products/internal/models"
)

// Handler represents the HTTP handlers for the service.
type Handler struct {
	db *database.DB
}

// New creates a new handler.
func New(db *database.DB) *Handler {
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

	products, err := h.db.GetProducts(r.Context(), limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get products")
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

// CreateProduct handles POST requests to /products.
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	productID, err := h.db.CreateProduct(r.Context(), &product)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	product.ID = productID

	respondWithJSON(w, http.StatusCreated, product)
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
