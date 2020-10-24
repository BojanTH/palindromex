package repository

import (
	"palindromex/web/db"
	"palindromex/web/model"
)

type User struct {
	Connection *db.Connection
}

func NewUser(c *db.Connection) *User {
	return &User{c}
}

func (userRepo *User) ByID(id int) model.User {
	userRepo.Connection.Open()
	defer userRepo.Connection.Close()

	u := model.User{}
	userRepo.Connection.Conn.First(&u, id)

	return u
}
