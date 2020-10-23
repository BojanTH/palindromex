package model

import (
	"time"
)

// User has many messages
type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"size:255"`
	Email     string `gorm:"unique;not null;size:255"`
	Password  []byte `gorm:"size:255"`
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []Message
	ApiKeys   []ApiKey
}

// TableName sets the table name to "_user" because "user" is a reserved word in Postgres
func (User) TableName() string {
	return "_user"
}
