package model

import (
	"time"
)

// ApiKey belongs to User. UserID is a foreign key
type ApiKey struct {
	ID        uint   `gorm:"primary_key"`
	UserID    uint
	Key       string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}
