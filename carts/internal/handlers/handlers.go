package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"com.MixieMeltsv2/carts/internal/db"
	"com.MixieMeltsv2/carts/internal/models"

	"github.com/go-chi/chi/v5"
)

// DBLayer abstracts the database operations the handlers need.
// This allows tests to inject a mock implementation.
type DBLayer interface {
	GetOrCreateCart(ctx context.Context, userID int64) (*models.Cart, error)
	// GetOrCreateCartBySession returns or creates a cart associated with an anonymous session id (guest carts).
	GetOrCreateCartBySession(ctx context.Context, sessionID string) (*models.Cart, error)

	GetCartItems(ctx context.Context, cartID int64) ([]models.CartItem, error)
	AddItem(ctx context.Context, cartID int64, productID int64, qty int) (*models.CartItem, error)
	UpdateItemQuantity(ctx context.Context, itemID int64, qty int) (*models.CartItem, error)
	RemoveItem(ctx context.Context, itemID int64) error

	CountCartsWithProduct(ctx context.Context, productID int64) (int, error)

	// MergeSessionCartIntoUserCart migrates items from a session cart into a user's cart (summing quantities).
	MergeSessionCartIntoUserCart(ctx context.Context, sessionID string, userID int64) error
}

// compile-time check: ensure the concrete DB type implements the DBLayer interface.
// This will fail to compile if the db.DB type (from internal/db) doesn't satisfy DBLayer.
var _ DBLayer = (*db.DB)(nil)

 // Handler exposes HTTP handlers for the carts service.
 type Handler struct {
 	db                  DBLayer
 	runningLowThreshold int
 	inventoryBase       string
 	// usersBase is an optional base URL for the users service (e.g. "http://users:8080").
 	// If empty, the handler will attempt to read USERS_SERVICE_URL from the environment.
 	usersBase           string
 	httpClient          *http.Client
 }

// New constructs a Handler.
func New(d DBLayer, runningLowThreshold int, inventoryBase string, usersBase string) *Handler {
	return &Handler{
		db:                  d,
		runningLowThreshold: runningLowThreshold,
		inventoryBase:       inventoryBase,
		usersBase:           usersBase,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// GetCart responds with the user's cart, creating one if necessary.
// Route: GET /carts/{userID}
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseIDParam(r, "userID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	cart, err := h.db.GetOrCreateCart(ctx, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load cart")
		return
	}

	cart.Items = h.computeRunningLow(ctx, cart.Items)
	respondJSON(w, http.StatusOK, cart)
}

// GetCartBySession responds with the session cart, creating one if necessary.
// Expects a "session_id" cookie to identify the guest session.
// Route: GET /carts
func (h *Handler) GetCartBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		respondError(w, http.StatusBadRequest, "missing session_id cookie")
		return
	}
	sessionID := cookie.Value

	cart, err := h.db.GetOrCreateCartBySession(ctx, sessionID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load cart")
		return
	}

	cart.Items = h.computeRunningLow(ctx, cart.Items)
	respondJSON(w, http.StatusOK, cart)
}

// AddItem adds an item to the user's cart and returns the created/updated item.
// Route: POST /carts/{userID}/items
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseIDParam(r, "userID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req models.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ProductID == 0 || req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "product_id and positive quantity required")
		return
	}

	cart, err := h.db.GetOrCreateCart(ctx, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load cart")
		return
	}

	item, err := h.db.AddItem(ctx, cart.ID, req.ProductID, req.Quantity)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to add item")
		return
	}

	// Compute running_low for the returned item only.
	items := []models.CartItem{*item}
	items = h.computeRunningLow(ctx, items)
	item.RunningLow = items[0].RunningLow

	respondJSON(w, http.StatusCreated, item)
}

// AddItemBySession adds an item to the session cart and returns the created/updated item.
// Expects a "session_id" cookie to identify the guest session.
// Route: POST /carts/items
func (h *Handler) AddItemBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		respondError(w, http.StatusBadRequest, "missing session_id cookie")
		return
	}
	sessionID := cookie.Value

	var req models.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ProductID == 0 || req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "product_id and positive quantity required")
		return
	}

	cart, err := h.db.GetOrCreateCartBySession(ctx, sessionID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load cart")
		return
	}

	item, err := h.db.AddItem(ctx, cart.ID, req.ProductID, req.Quantity)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to add item")
		return
	}

	// Compute running_low for the returned item only.
	items := []models.CartItem{*item}
	items = h.computeRunningLow(ctx, items)
	item.RunningLow = items[0].RunningLow

	respondJSON(w, http.StatusCreated, item)
}

// UpdateItem updates the quantity for an item in a cart and returns the updated item.
// Route: PATCH /carts/{userID}/items/{itemID}
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := parseIDParam(r, "userID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	itemID, err := parseIDParam(r, "itemID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req models.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "quantity must be > 0")
		return
	}

	item, err := h.db.UpdateItemQuantity(ctx, itemID, req.Quantity)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update item")
		return
	}

	items := []models.CartItem{*item}
	items = h.computeRunningLow(ctx, items)
	item.RunningLow = items[0].RunningLow

	respondJSON(w, http.StatusOK, item)
}

// UpdateItemBySession updates the quantity for an item in a session cart and returns the updated item.
// Expects a "session_id" cookie to identify the guest session.
// Route: PATCH /carts/items/{itemID}
func (h *Handler) UpdateItemBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		respondError(w, http.StatusBadRequest, "missing session_id cookie")
		return
	}

	itemID, err := parseIDParam(r, "itemID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req models.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "quantity must be > 0")
		return
	}

	item, err := h.db.UpdateItemQuantity(ctx, itemID, req.Quantity)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update item")
		return
	}

	items := []models.CartItem{*item}
	items = h.computeRunningLow(ctx, items)
	item.RunningLow = items[0].RunningLow

	respondJSON(w, http.StatusOK, item)
}

// DeleteItem removes an item from the cart.
// Route: DELETE /carts/{userID}/items/{itemID}
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := parseIDParam(r, "userID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	itemID, err := parseIDParam(r, "itemID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.db.RemoveItem(ctx, itemID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to remove item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteItemBySession removes an item from the session cart.
// Expects a "session_id" cookie to identify the guest session.
// Route: DELETE /carts/items/{itemID}
func (h *Handler) DeleteItemBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		respondError(w, http.StatusBadRequest, "missing session_id cookie")
		return
	}

	itemID, err := parseIDParam(r, "itemID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.db.RemoveItem(ctx, itemID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to remove item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// MergeCart migrates the session cart into the authenticated user's cart.
// Expects a JSON body: { "user_id": <int> } and requires the session_id cookie.
// Route: POST /carts/merge
func (h *Handler) MergeCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		respondError(w, http.StatusBadRequest, "missing session_id cookie")
		return
	}
	sessionID := cookie.Value

	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == 0 {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.db.MergeSessionCartIntoUserCart(ctx, sessionID, req.UserID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to merge carts")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Helpers ---

func parseIDParam(r *http.Request, name string) (int64, error) {
	idStr := chi.URLParam(r, name)
	if idStr == "" {
		return 0, fmt.Errorf("missing %s", name)
	}
	return strconv.ParseInt(idStr, 10, 64)
}

func respondJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, code int, message string) {
	http.Error(w, message, code)
}

// computeRunningLow enriches the provided items with the RunningLow flag using two signals:
// - popularity (count of distinct carts containing the product)
// - optional inventory probe (inventoryBase)
//
// This is intentionally synchronous and best-effort (inventory probe timeouts/errors are ignored).
func (h *Handler) computeRunningLow(ctx context.Context, items []models.CartItem) []models.CartItem {
	for i := range items {
		it := &items[i]
		count, err := h.db.CountCartsWithProduct(ctx, it.ProductID)
		if err == nil && count >= h.runningLowThreshold {
			it.RunningLow = true
			continue
		}

		// best-effort inventory probe
		if h.inventoryBase != "" {
			if s, err := h.fetchInventoryStock(ctx, it.ProductID); err == nil {
				if s.Stock <= s.Threshold {
					it.RunningLow = true
					continue
				}
			}
		}

		it.RunningLow = false
	}
	return items
}

// fetchInventoryStock queries the configured inventory service for a product's stock.
// Expected response shape matches models.InventoryStockResponse.
func (h *Handler) fetchInventoryStock(ctx context.Context, productID int64) (*models.InventoryStockResponse, error) {
	if h.inventoryBase == "" {
		return nil, fmt.Errorf("no inventory base configured")
	}
	url := fmt.Sprintf("%s/products/%d/stock", h.inventoryBase, productID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inventory responded with status %d", resp.StatusCode)
	}
	var out models.InventoryStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
