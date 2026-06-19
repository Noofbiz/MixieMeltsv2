package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.MixieMeltsv2/carts/internal/handlers"
)

// mockDB implements handlers.DBLayer for tests.
type mockDB struct {
	GetOrCreateCartFunc            func(ctx context.Context, userID int64) (*handlers.Cart, error)
	GetOrCreateCartBySessionFunc   func(ctx context.Context, sessionID string) (*handlers.Cart, error)
	GetCartItemsFunc               func(ctx context.Context, cartID int64) ([]handlers.CartItem, error)
	AddItemFunc                    func(ctx context.Context, cartID int64, productID int64, qty int) (*handlers.CartItem, error)
	UpdateItemQuantityFunc         func(ctx context.Context, itemID int64, qty int) (*handlers.CartItem, error)
	RemoveItemFunc                 func(ctx context.Context, itemID int64) error
	CountCartsWithProductFunc      func(ctx context.Context, productID int64) (int, error)
	MergeSessionCartIntoUserCartFn func(ctx context.Context, sessionID string, userID int64) error
}

func (m *mockDB) GetOrCreateCart(ctx context.Context, userID int64) (*handlers.Cart, error) {
	if m.GetOrCreateCartFunc != nil {
		return m.GetOrCreateCartFunc(ctx, userID)
	}
	return &handlers.Cart{}, nil
}
func (m *mockDB) GetOrCreateCartBySession(ctx context.Context, sessionID string) (*handlers.Cart, error) {
	if m.GetOrCreateCartBySessionFunc != nil {
		return m.GetOrCreateCartBySessionFunc(ctx, sessionID)
	}
	return &handlers.Cart{}, nil
}
func (m *mockDB) GetCartItems(ctx context.Context, cartID int64) ([]handlers.CartItem, error) {
	if m.GetCartItemsFunc != nil {
		return m.GetCartItemsFunc(ctx, cartID)
	}
	return []handlers.CartItem{}, nil
}
func (m *mockDB) AddItem(ctx context.Context, cartID int64, productID int64, qty int) (*handlers.CartItem, error) {
	if m.AddItemFunc != nil {
		return m.AddItemFunc(ctx, cartID, productID, qty)
	}
	return &handlers.CartItem{}, nil
}
func (m *mockDB) UpdateItemQuantity(ctx context.Context, itemID int64, qty int) (*handlers.CartItem, error) {
	if m.UpdateItemQuantityFunc != nil {
		return m.UpdateItemQuantityFunc(ctx, itemID, qty)
	}
	return &handlers.CartItem{}, nil
}
func (m *mockDB) RemoveItem(ctx context.Context, itemID int64) error {
	if m.RemoveItemFunc != nil {
		return m.RemoveItemFunc(ctx, itemID)
	}
	return nil
}
func (m *mockDB) CountCartsWithProduct(ctx context.Context, productID int64) (int, error) {
	if m.CountCartsWithProductFunc != nil {
		return m.CountCartsWithProductFunc(ctx, productID)
	}
	return 0, nil
}
func (m *mockDB) MergeSessionCartIntoUserCart(ctx context.Context, sessionID string, userID int64) error {
	if m.MergeSessionCartIntoUserCartFn != nil {
		return m.MergeSessionCartIntoUserCartFn(ctx, sessionID, userID)
	}
	return nil
}

// TestMergeCart_Success verifies MergeCart requires a valid Authorization header,
// validates the token with the users service, and invokes the DB merge function.
func TestMergeCart_Success(t *testing.T) {
	// Mock users service: returns 200 + {"id":42} when Authorization == "Bearer good-token"
	usersSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users/me" {
			http.NotFound(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer good-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": 42})
	}))
	defer usersSrv.Close()

	called := false
	var gotSession string
	var gotUserID int64

	mock := &mockDB{
		MergeSessionCartIntoUserCartFn: func(ctx context.Context, sessionID string, userID int64) error {
			called = true
			gotSession = sessionID
			gotUserID = userID
			return nil
		},
	}

	// Create handler pointing at mocked users service
	h := handlers.New(mock, 5, "", usersSrv.URL)

	// Create request with Authorization and session cookie
	reqBody := bytes.NewReader([]byte(`{}`))
	req := httptest.NewRequest(http.MethodPost, "/carts/merge", reqBody)
	req.Header.Set("Authorization", "Bearer good-token")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "sess-abc-123", Path: "/"})

	rr := httptest.NewRecorder()
	h.MergeCart(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d; body=%s", rr.Code, rr.Body.String())
	}
	if !called {
		t.Fatalf("expected MergeSessionCartIntoUserCart to be called")
	}
	if gotSession != "sess-abc-123" {
		t.Fatalf("expected session 'sess-abc-123', got %q", gotSession)
	}
	if gotUserID != 42 {
		t.Fatalf("expected userID 42, got %d", gotUserID)
	}
	// Confirm cookie cleared in response (Set-Cookie present with session_id cleared)
	setCookies := rr.HeaderMap["Set-Cookie"]
	foundCleared := false
	for _, sc := range setCookies {
		if sc == "" {
			continue
		}
		// Should contain session_id and Max-Age=-1 or Value empty
		if (contains(sc, "session_id=") && (contains(sc, "Max-Age=-1") || contains(sc, "Expires="))) {
			foundCleared = true
			break
		}
	}
	if !foundCleared {
		t.Fatalf("expected session cookie to be cleared, set-cookie headers: %v", setCookies)
	}
}

// TestMergeCart_MissingAuth ensures the handler rejects requests without Authorization header.
func TestMergeCart_MissingAuth(t *testing.T) {
	mock := &mockDB{
		MergeSessionCartIntoUserCartFn: func(ctx context.Context, sessionID string, userID int64) error {
			t.Fatalf("merge should not be called when auth missing")
			return nil
		},
	}
	h := handlers.New(mock, 5, "", "http://unused")

	req := httptest.NewRequest(http.MethodPost, "/carts/merge", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "sess-1", Path: "/"})

	rr := httptest.NewRecorder()
	h.MergeCart(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when missing auth, got %d", rr.Code)
	}
}

// TestMergeCart_InvalidToken ensures invalid tokens are rejected and merge not performed.
func TestMergeCart_InvalidToken(t *testing.T) {
	// mock users service returns 401 for any token
	usersSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer usersSrv.Close()

	called := false
	mock := &mockDB{
		MergeSessionCartIntoUserCartFn: func(ctx context.Context, sessionID string, userID int64) error {
			called = true
			return nil
		},
	}
	h := handlers.New(mock, 5, "", usersSrv.URL)

	req := httptest.NewRequest(http.MethodPost, "/carts/merge", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "sess-2", Path: "/"})

	rr := httptest.NewRecorder()
	h.MergeCart(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", rr.Code)
	}
	if called {
		t.Fatalf("merge should not have been called on invalid token")
	}
}

// contains is a tiny helper to avoid importing strings in many places.
func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}
