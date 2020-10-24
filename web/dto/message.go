package dto

import (
	"time"
)

// Message is used to expose only data that should be exposed in the response
type Message struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id`
	Content    string    `json:"content`
	Palindrome bool      `json:"palindrome"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
