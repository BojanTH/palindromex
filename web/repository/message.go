package repository

import (
	"palindromex/web/db"
)

type Message struct {
	Connection *db.Connection
}

func NewMessage(c *db.Connection) *Message {
	return &Message{c}
}
