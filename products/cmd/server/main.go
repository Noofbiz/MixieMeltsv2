package main

import (
	"log"
	"net/http"
	"os"

	"com.MixieMelts.products/internal/database"
	"com.MixieMelts.products/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	// Database connection
	db, err := database.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Initialize handlers
	h := handlers.New(db)

	// Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// CORS middleware
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"frontend:"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}).Handler)

	r.Get("/products", h.GetProducts)
	r.Post("/products", h.CreateProduct)

	r.Get("/products/subscription-boxes", h.GetSubscriptionBoxes)
	r.Post("/products/subscription-boxes", h.CreateSubscriptionBox)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Println("Starting products service on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
