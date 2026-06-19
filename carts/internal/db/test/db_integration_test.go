package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	cartdb "com.MixieMeltsv2/carts/internal/db"
	"com.MixieMeltsv2/carts/internal/models"
)

// TestMergeSessionCartIntoUserCartIntegration performs an integration test against a real Postgres
// instance. It requires DATABASE_URL environment variable to be set (matching the GH Actions service).
//
// The test:
//  - creates a session cart and adds items
//  - creates a user cart with one overlapping product
//  - calls MergeSessionCartIntoUserCart to move/merge items from session into user cart
//  - asserts that user cart quantities were summed and the session cart is effectively cleared
func TestMergeSessionCartIntoUserCartIntegration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set; skipping integration test")
	}

	// short context so tests fail fast on CI if DB isn't available
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	d, err := cartdb.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("db.New failed: %v", err)
	}
	defer d.Close()

	// Use unique ids unlikely to collide with seeded data
	sessionID := "itest-session-" + time.Now().Format("20060102150405")
	userID := int64(9001)

	// Ensure a session cart exists and add items
	sessCart, err := d.GetOrCreateCartBySession(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetOrCreateCartBySession failed: %v", err)
	}

	// Add some items to the session cart
	if _, err := d.AddItem(ctx, sessCart.ID, 10001, 2); err != nil {
		t.Fatalf("AddItem to session cart failed: %v", err)
	}
	if _, err := d.AddItem(ctx, sessCart.ID, 10002, 3); err != nil {
		t.Fatalf("AddItem to session cart failed: %v", err)
	}

	// Ensure a user cart exists and add an overlapping product (to test quantity summing)
	userCart, err := d.GetOrCreateCart(ctx, userID)
	if err != nil {
		t.Fatalf("GetOrCreateCart for user failed: %v", err)
	}
	if _, err := d.AddItem(ctx, userCart.ID, 10001, 1); err != nil {
		t.Fatalf("AddItem to user cart failed: %v", err)
	}

	// Perform the merge
	if err := d.MergeSessionCartIntoUserCart(ctx, sessionID, userID); err != nil {
		t.Fatalf("MergeSessionCartIntoUserCart failed: %v", err)
	}

	// Reload user cart and verify items
	uc, err := d.GetOrCreateCart(ctx, userID)
	if err != nil {
		t.Fatalf("GetOrCreateCart (after merge) failed: %v", err)
	}

	found10001 := false
	found10002 := false
	for _, it := range uc.Items {
		if it.ProductID == 10001 {
			found10001 = true
			// expected quantity = previous user(1) + session(2) = 3
			if it.Quantity != 3 {
				t.Fatalf("unexpected quantity for product 10001: want 3 got %d", it.Quantity)
			}
		}
		if it.ProductID == 10002 {
			found10002 = true
			if it.Quantity != 3 {
				t.Fatalf("unexpected quantity for product 10002: want 3 got %d", it.Quantity)
			}
		}
	}
	if !found10001 {
		t.Fatal("expected product 10001 in user cart after merge")
	}
	if !found10002 {
		t.Fatal("expected product 10002 in user cart after merge")
	}

	// Verify session cart was removed (GetOrCreateCartBySession should create an empty cart)
	afterSess, err := d.GetOrCreateCartBySession(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetOrCreateCartBySession after merge failed: %v", err)
	}
	if len(afterSess.Items) != 0 {
		t.Fatalf("expected session cart to be empty after merge, got %d items", len(afterSess.Items))
	}

	// Cleanup: try to delete test carts to keep DB tidy (best-effort)
	// Note: implementation uses cascade deletes for cart_items.
	_, _ = d.GetOrCreateCart(ctx, userID) // ensure exists
	// We won't assert on cleanup errors to avoid masking test results.
}
