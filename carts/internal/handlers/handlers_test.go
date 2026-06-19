package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.MixieMeltsv2/carts/internal/models"

	"github.com/go-chi/chi/v5"
)

// MockDB is a test double implementing the DBLayer interface.
type MockDB struct {
	GetOrCreateCartFunc      func(ctx context.Context, userID int64) (*models.Cart, error)
	GetCartItemsFunc         func(ctx context.Context, cartID int64) ([]models.CartItem, error)
	AddItemFunc              func(ctx context.Context, cartID int64, productID int64, qty int) (*models.CartItem, error)
	UpdateItemQuantityFunc   func(ctx context.Context, itemID int64, qty int) (*models.CartItem, error)
	RemoveItemFunc           func(ctx context.Context, itemID int64) error
	CountCartsWithProductFunc func(ctx context.Context, productID int64) (int, error)
}

func (m *MockDB) GetOrCreateCart(ctx context.Context, userID int64) (*models.Cart, error) {
	if m.GetOrCreateCartFunc != nil {
		return m.GetOrCreateCartFunc(ctx, userID)
	}
	return nil, errors.New("GetOrCreateCartFunc not implemented")
}

func (m *MockDB) GetCartItems(ctx context.Context, cartID int64) ([]models.CartItem, error) {
	if m.GetCartItemsFunc != nil {
		return m.GetCartItemsFunc(ctx, cartID)
	}
	return nil, errors.New("GetCartItemsFunc not implemented")
}

func (m *MockDB) AddItem(ctx context.Context, cartID int64, productID int64, qty int) (*models.CartItem, error) {
	if m.AddItemFunc != nil {
		return m.AddItemFunc(ctx, cartID, productID, qty)
	}
	return nil, errors.New("AddItemFunc not implemented")
}

func (m *MockDB) UpdateItemQuantity(ctx context.Context, itemID int64, qty int) (*models.CartItem, error) {
	if m.UpdateItemQuantityFunc != nil {
		return m.UpdateItemQuantityFunc(ctx, itemID, qty)
	}
	return nil, errors.New("UpdateItemQuantityFunc not implemented")
}

func (m *MockDB) RemoveItem(ctx context.Context, itemID int64) error {
	if m.RemoveItemFunc != nil {
		return m.RemoveItemFunc(ctx, itemID)
	}
	return errors.New("RemoveItemFunc not implemented")
}

func (m *MockDB) CountCartsWithProduct(ctx context.Context, productID int64) (int, error) {
	if m.CountCartsWithProductFunc != nil {
		return m.CountCartsWithProductFunc(ctx, productID)
	}
	return 0, errors.New("CountCartsWithProductFunc not implemented")
}

// helper to mount a handler on a chi router so URL params are populated.
func mountRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/carts/{userID}", h.GetCart)
	r.Post("/carts/{userID}/items", h.AddItem)
	r.Patch("/carts/{userID}/items/{itemID}", h.UpdateItem)
	r.Delete("/carts/{userID}/items/{itemID}", h.DeleteItem)
	return r
}

func TestGetCart_Success(t *testing.T) {
	mock := &MockDB{
		GetOrCreateCartFunc: func(ctx context.Context, userID int64) (*models.Cart, error) {
			return &models.Cart{
				ID:     10,
				UserID: userID,
				Items:  []models.CartItem{},
			}, nil
		},
	}

	h := New(mock, 5, "")
	r := mountRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/carts/42", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var cart models.Cart
	if err := json.NewDecoder(rr.Body).Decode(&cart); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if cart.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", cart.UserID)
	}
}

func TestAddItem_Success(t *testing.T) {
	mock := &MockDB{
		GetOrCreateCartFunc: func(ctx context.Context, userID int64) (*models.Cart, error) {
			return &models.Cart{ID: 7, UserID: userID}, nil
		},
		AddItemFunc: func(ctx context.Context, cartID int64, productID int64, qty int) (*models.CartItem, error) {
			return &models.CartItem{ID: 123, CartID: cartID, ProductID: productID, Quantity: qty}, nil
		},
		CountCartsWithProductFunc: func(ctx context.Context, productID int64) (int, error) {
			return 0, nil
		},
	}

	h := New(mock, 5, "")
	r := mountRouter(h)

	body := map[string]interface{}{
		"product_id": 5,
		"quantity":   2,
	}
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/carts/100/items", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d. body=%s", rr.Code, rr.Body.String())
	}

	var item models.CartItem
	if err := json.NewDecoder(rr.Body).Decode(&item); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if item.ProductID != 5 || item.Quantity != 2 {
		t.Fatalf("unexpected item in response: %+v", item)
	}
}

func TestAddItem_InvalidBody(t *testing.T) {
	mock := &MockDB{}
	h := New(mock, 5, "")
	r := mountRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/carts/1/items", bytes.NewReader([]byte(`{"product_id":`))) // malformed JSON
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 bad request for malformed body, got %d", rr.Code)
	}
}

func TestUpdateItem_Success(t *testing.T) {
	mock := &MockDB{
		UpdateItemQuantityFunc: func(ctx context.Context, itemID int64, qty int) (*models.CartItem, error) {
			return &models.CartItem{ID: itemID, CartID: 2, ProductID: 9, Quantity: qty}, nil
		},
		CountCartsWithProductFunc: func(ctx context.Context, productID int64) (int, error) {
			return 0, nil
		},
	}

	h := New(mock, 5, "")
	r := mountRouter(h)

	body := map[string]interface{}{"quantity": 4}
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/carts/200/items/55", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
	var item models.CartItem
	if err := json.NewDecoder(rr.Body).Decode(&item); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if item.ID != 55 || item.Quantity != 4 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestDeleteItem_Success(t *testing.T) {
	removed := false
	mock := &MockDB{
		RemoveItemFunc: func(ctx context.Context, itemID int64) error {
			removed = true
			return nil
		},
	}

	h := New(mock, 5, "")
	r := mountRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/carts/1/items/500", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", rr.Code)
	}
	if !removed {
		t.Fatalf("expected RemoveItem to be called")
	}
}

func TestGetCart_DBError(t *testing.T) {
	mock := &MockDB{
		GetOrCreateCartFunc: func(ctx context.Context, userID int64) (*models.Cart, error) {
			return nil, errors.New("db down")
		},
	}
	h := New(mock, 5, "")
	r := mountRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/carts/99", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when DB returns error, got %d", rr.Code)
	}
}

func TestAddItem_DBError(t *testing.T) {
	mock := &MockDB{
		GetOrCreateCartFunc: func(ctx context.Context, userID int64) (*models.Cart, error) {
			return &models.Cart{ID: 1, UserID: userID}, nil
		},
		AddItemFunc: func(ctx context.Context, cartID int64, productID int64, qty int) (*models.CartItem, error) {
			return nil, errors.New("cannot insert")
		},
	}
	h := New(mock, 5, "")
	r := mountRouter(h)

	body := map[string]interface{}{"product_id": 3, "quantity": 1}
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/carts/1/items", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when AddItem fails, got %d", rr.Code)
	}
}
