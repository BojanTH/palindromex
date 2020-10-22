package model

import (
	"time"
)

// Message belongs to User. UserID is a foreign key
type Message struct {
	ID        uint   `gorm:"primary_key"`
	UserID    uint
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}
