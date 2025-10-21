package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"com.MixieMelts.users/internal/database"
	"com.MixieMelts.users/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var db *database.DB

// --- MAIN FUNCTION ---
func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	// Database connection
	db, err = database.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	r := chi.NewRouter()

	// Initialize handlers
	h := handlers.New(db, []byte(os.Getenv("JWT_SECRET_KEY")))

	// Middleware stack
	r.Use(middleware.Logger)    // Logs requests to the console
	r.Use(middleware.Recoverer) // Recovers from panics without crashing
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"frontend:"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}).Handler)

	// Public routes for local auth
	r.Post("/api/users/register", h.RegisterUser)
	r.Post("/api/users/login", h.LoginUser)

	// Public routes for OAuth
	r.Get("/api/users/oauth/google/login", h.HandleGoogleLogin)
	r.Get("/api/users/oauth/google/callback", h.HandleGoogleCallback)
	r.Get("/api/users/oauth/facebook/login", h.HandleFacebookLogin)
	r.Get("/api/users/oauth/facebook/callback", h.HandleFacebookCallback)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Get("/api/users/me", h.GetUserProfile)
		// Add other protected routes here (e.g., PUT /api/users/me)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting User Service on port " + port)
	log.Println("Ensure GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, FACEBOOK_CLIENT_ID, and FACEBOOK_CLIENT_SECRET environment variables are set.")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
