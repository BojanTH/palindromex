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

func (repo *Message) FindAllByUserID(userID int) ([]model.Message, error) {
	repo.Connection.Open()
	defer repo.Connection.Close()

	var users []model.Message
	result := repo.Connection.Conn.Find(&users, "user_id = ?", userID)

	return users, result.Error
}


func (repo *Message) FindByUserIDAndID(userID, messageID int) model.Message {
	repo.Connection.Open()
	defer repo.Connection.Close()

	message := model.Message{}
	repo.Connection.Conn.First(&message, "id = ? AND user_id = ?", messageID, userID)

	return message
}
