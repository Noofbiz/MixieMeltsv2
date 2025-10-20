package models

// Credentials represents the authentication credentials for a user.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
