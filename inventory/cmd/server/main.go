package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"com.MixieMelts.inventory/internal/database"
	"com.MixieMelts.inventory/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load; falling back to environment variables")
	}

	// Create database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	db, err := database.New(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize handlers with the db
	h := handlers.New(db)

	// Router setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS - allow frontend to call the API
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"frontend:"}, // adjust origin in env or here for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// Ingredient (material) routes - track bases, waxes, scent oils, etc.
	r.Get("/ingredients", h.GetIngredients)
	r.Get("/ingredients/{id}", h.GetIngredient)
	r.Post("/ingredients", h.CreateIngredient)
	r.Put("/ingredients/{id}", h.UpdateIngredient)
	r.Patch("/ingredients/{id}/adjust", h.AdjustIngredientStock)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Starting inventory service on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
