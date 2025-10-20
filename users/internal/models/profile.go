package models

// Profile represents a user's profile information.
type Profile struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}
