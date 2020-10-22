package repository

import (
	"palindromex/web/db"
)

type User struct {
	Connection *db.Connection
}

func NewUser(c *db.Connection) *User {
	return &User{c}
}
