package repository

import (
	"palindromex/web/db"
	"palindromex/web/model"
)

type Message struct {
	Connection *db.Connection
}

func NewMessage(c *db.Connection) *Message {
	return &Message{c}
}

func (repo *Message) CreateNew(message model.Message) {
	repo.Connection.Open()
	defer repo.Connection.Close()

	repo.Connection.Conn.Create(&message)
}
