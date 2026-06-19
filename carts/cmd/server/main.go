package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"com.MixieMeltsv2/carts/internal/db"
	"com.MixieMeltsv2/carts/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

/*
Main now wires the internal packages:
- uses internal/db.New(ctx, dbURL) to initialise the database (creates tables/indices and seeds)
- constructs internal/handlers with the DB and config
- mounts handler methods on the router
*/

func main() {
	// load optional .env
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Automatic migration-on-start removed.
	// Schema migrations are now applied explicitly in CI or via the top-level Makefile / service Makefile.

	d, err := db.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer d.Close()

	// threshold for marking running low by cart count
	threshold := 5
	if v := os.Getenv("RUNNING_LOW_CART_THRESHOLD"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			threshold = n
		}
	}
	inventoryBase := os.Getenv("INVENTORY_SERVICE_URL") // optional
	// Base URL for the users service (used to validate auth for merge operations).
	// Configure via env var USERS_SERVICE_URL (e.g. http://users:8080).
	usersBase := os.Getenv("USERS_SERVICE_URL")

	h := handlers.New(d, threshold, inventoryBase, usersBase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// health
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// cart endpoints wired to the internal handlers (authenticated user routes)
	r.Get("/carts/{userID}", h.GetCart)
	r.Post("/carts/{userID}/items", h.AddItem)
	r.Patch("/carts/{userID}/items/{itemID}", h.UpdateItem)
	r.Delete("/carts/{userID}/items/{itemID}", h.DeleteItem)

	// guest/session-backed routes (use cookie-based session_id)
	// GET /carts                     -> returns or creates a session cart
	// POST /carts/items               -> add item to session cart
	// PATCH /carts/items/{itemID}     -> update session cart item
	// DELETE /carts/items/{itemID}    -> delete session cart item
	r.Get("/carts", h.GetCartBySession)
	r.Post("/carts/items", h.AddItemBySession)
	r.Patch("/carts/items/{itemID}", h.UpdateItemBySession)
	r.Delete("/carts/items/{itemID}", h.DeleteItemBySession)

	// Merge endpoint: migrate session cart into authenticated user's cart after login.
	// Expects JSON body: { "user_id": <int> } and requires the session cookie to be present.
	r.Post("/carts/merge", h.MergeCart)

	// CORS wrapper
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"frontend:"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("starting carts service on :%s (running_low_threshold=%d inventory_base=%q)", port, threshold, inventoryBase)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
