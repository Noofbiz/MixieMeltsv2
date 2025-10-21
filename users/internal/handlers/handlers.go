package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"com.MixieMelts.users/internal/auth"
	"com.MixieMelts.users/internal/database"
	"com.MixieMelts.users/internal/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const oauthStateCookieName = "oauthstate"

type Handler struct {
	db           *database.DB
	jwtSecretKey []byte
}

func New(db *database.DB, jwtSecretKey []byte) *Handler {
	return &Handler{db: db, jwtSecretKey: jwtSecretKey}
}

// --- HANDLERS ---
// In a real app, these would be in an 'api' or 'handlers' package.

// RegisterUser handles new user creation.
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Check if user already exists
	_, err = h.db.GetUserByEmail(r.Context(), creds.Email)
	if err == nil {
		respondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	newUser := &models.User{
		Email:    creds.Email,
		Password: string(hashedPassword),
		Username: "New User", // In a real app, you'd get this from the request
	}

	userID, err := h.db.CreateUser(r.Context(), newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	newUser.ID = userID

	log.Printf("User registered: %s", creds.Email)
	respondWithJSON(w, http.StatusCreated, newUser)
}

// LoginUser handles user authentication and issues a JWT.
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), creds.Email)
	if err != nil || user == nil {
		spew.Dump(user)
		spew.Dump(creds)
		respondWithError(w, http.StatusUnauthorized, "Invalid credentialzorz!")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	h.issueJWT(w, *user)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetUserProfile is a protected handler that returns the user's profile.
func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// The userID is added to the request context by the AuthMiddleware.
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching profile for user ID: %d", userID)

	user, err := h.db.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// --- OAUTH HANDLERS ---
func (h *Handler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := h.generateStateOauthCookie(w)
	url := auth.GetGoogleAuthUrl(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie(oauthStateCookieName)
	if err != nil || r.FormValue("state") != state.Value {
		log.Printf("Invalid oauth state, expected '%s', got '%s'\n", state.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := auth.GetGoogleOAuthExchangeToken(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("Code exchange failed: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Failed getting user info: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed reading response body: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.Unmarshal(contents, &userInfo)

	// Check if user exists, if not, create one
	user, err := h.db.GetUserByEmail(r.Context(), userInfo.Email)
	if err != nil {
		// Handle error, but for now, assume it means user doesn't exist
		newUser := &models.User{
			Email:    userInfo.Email,
			Username: userInfo.Name,
			Password: "", // No password for OAuth users
		}
		userID, err := h.db.CreateUser(r.Context(), newUser)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		newUser.ID = userID
		user = newUser
		log.Printf("New user created via Google: %s", userInfo.Email)
	}

	h.issueJWT(w, *user)
}

func (h *Handler) HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	state := h.generateStateOauthCookie(w)
	url := auth.GetFacebookAuthUrl(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie(oauthStateCookieName)
	if err != nil || r.FormValue("state") != state.Value {
		log.Printf("Invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := auth.GetFacebookOAuthExchangeToken(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("Code exchange failed: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", token.AccessToken))
	if err != nil {
		log.Printf("Failed getting user info: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// Check if user exists, if not, create one
	user, err := h.db.GetUserByEmail(r.Context(), userInfo.Email)
	if err != nil {
		// Handle error, but for now, assume it means user doesn't exist
		newUser := &models.User{
			Email:    userInfo.Email,
			Username: userInfo.Name,
			Password: "", // No password for OAuth users
		}
		userID, err := h.db.CreateUser(r.Context(), newUser)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		newUser.ID = userID
		user = newUser
		log.Printf("New user created via Facebook: %s", userInfo.Email)
	}

	h.issueJWT(w, *user)
}

// --- UTILITY & MIDDLEWARE ---
func (h *Handler) generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return state
}

func (h *Handler) issueJWT(w http.ResponseWriter, user models.User) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprint(user.ID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecretKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	log.Printf("Issued JWT for user: %s", user.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return h.jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to the request context for protected handlers to use
		ctx := context.WithValue(r.Context(), "userID", claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
