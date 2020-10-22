package dto

// Credentials for a user
type Credentials struct {
	UserID          uint   `json:"user_id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
}
