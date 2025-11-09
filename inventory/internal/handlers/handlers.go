package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"com.MixieMelts.inventory/internal/database"
	"com.MixieMelts.inventory/internal/models"

	"github.com/go-chi/chi/v5"
)

// Handler provides HTTP handlers for the inventory service.
type Handler struct {
	db *database.DB
}

// New creates a new Handler.
func New(db *database.DB) *Handler {
	return &Handler{db: db}
}

// JSON helpers
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

// AdjustInventory adjusts the quantity for a finished-product inventory item.
type adjustPayload struct {
	Change int64  `json:"change"` // positive or negative integer
	Reason string `json:"reason,omitempty"`
}

// -------------------- Ingredient Handlers --------------------

// GetIngredients returns all ingredients (wax, bases, scent oils).
func (h *Handler) GetIngredients(w http.ResponseWriter, r *http.Request) {
	list, err := h.db.GetIngredients(r.Context())
	if err != nil {
		log.Printf("GetIngredients error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to list ingredients")
		return
	}
	respondWithJSON(w, http.StatusOK, list)
}

// GetIngredient returns a single ingredient by id.
func (h *Handler) GetIngredient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	it, err := h.db.GetIngredient(r.Context(), id)
	if err != nil {
		log.Printf("GetIngredient error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to get ingredient")
		return
	}
	if it == nil {
		respondWithError(w, http.StatusNotFound, "ingredient not found")
		return
	}
	respondWithJSON(w, http.StatusOK, it)
}

// CreateIngredient creates a new ingredient record.
func (h *Handler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	var it models.Ingredient
	if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.db.CreateIngredient(r.Context(), &it)
	if err != nil {
		log.Printf("CreateIngredient error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create ingredient")
		return
	}
	it.ID = id
	respondWithJSON(w, http.StatusCreated, it)
}

// UpdateIngredient updates ingredient metadata and stock (full replace semantics).
func (h *Handler) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var it models.Ingredient
	if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	it.ID = id

	if err := h.db.UpdateIngredient(r.Context(), &it); err != nil {
		log.Printf("UpdateIngredient error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to update ingredient")
		return
	}
	respondWithJSON(w, http.StatusOK, it)
}

// AdjustIngredientPayload is used for stock adjustments.
type AdjustIngredientPayload struct {
	Change    float64 `json:"change"`               // positive to add, negative to subtract
	Reason    string  `json:"reason,omitempty"`     // e.g., restock, sale, waste
	Reference string  `json:"reference,omitempty"`  // e.g., PO number, order id
	CreatedBy string  `json:"created_by,omitempty"` // user or system who adjusted
}

// AdjustIngredientStock adjusts stock and records the adjustment.
func (h *Handler) AdjustIngredientStock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	ingredientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p AdjustIngredientPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	newStock, err := h.db.AdjustIngredientStock(r.Context(), ingredientID, p.Change, p.Reason, p.Reference, p.CreatedBy)
	if err != nil {
		log.Printf("AdjustIngredientStock error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to adjust ingredient stock")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"ingredient_id": ingredientID,
		"new_stock":     newStock,
	})
}

// -------------------- Recipe Handlers --------------------

// GetRecipe returns recipe items associated with a product.
func (h *Handler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	prodIDStr := r.URL.Query().Get("product_id")
	if prodIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "product_id query param required")
		return
	}
	prodID, err := strconv.ParseInt(prodIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid product_id")
		return
	}
	items, err := h.db.GetRecipe(r.Context(), prodID)
	if err != nil {
		log.Printf("GetRecipe error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}
	respondWithJSON(w, http.StatusOK, items)
}

// CreateRecipeItem adds an ingredient usage entry for a product.
func (h *Handler) CreateRecipeItem(w http.ResponseWriter, r *http.Request) {
	var ri models.RecipeItem
	if err := json.NewDecoder(r.Body).Decode(&ri); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.db.CreateRecipeItem(r.Context(), &ri)
	if err != nil {
		log.Printf("CreateRecipeItem error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create recipe item")
		return
	}
	ri.ID = id
	respondWithJSON(w, http.StatusCreated, ri)
}

// -------------------- Utility / Diagnostics --------------------

// Health check (simple)
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Optional helper to compute production capacity for a product based on current ingredient stocks.
// Returns the maximum number of product units that can be produced (may be fractional).
func (h *Handler) ComputeProductCapacity(ctx context.Context, productID int64) (*models.ProductStockProjection, error) {
	recipe, err := h.db.GetRecipe(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}
	if len(recipe) == 0 {
		return &models.ProductStockProjection{
			ProductID:      productID,
			AvailableUnits: 0,
		}, nil
	}

	// For each recipe item, fetch ingredient to determine available units.
	var limitingItemID int64
	var limitingUnits float64 = 1e18 // very large
	for _, ri := range recipe {
		ing, err := h.db.GetIngredient(ctx, ri.IngredientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get ingredient %d: %w", ri.IngredientID, err)
		}
		if ing == nil {
			// missing ingredient -> cannot produce
			return &models.ProductStockProjection{
				ProductID:      productID,
				AvailableUnits: 0,
				LimitingItemID: ri.IngredientID,
				LimitingAvail:  0,
			}, nil
		}
		if ri.Amount <= 0 {
			continue
		}
		possible := ing.Stock / ri.Amount
		if possible < limitingUnits {
			limitingUnits = possible
			limitingItemID = ing.ID
		}
	}

	if limitingUnits > 1e17 {
		limitingUnits = 0
	}

	return &models.ProductStockProjection{
		ProductID:      productID,
		AvailableUnits: limitingUnits,
		LimitingItemID: limitingItemID,
		LimitingAvail:  limitingUnits,
	}, nil
}
