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


func (repo *User) CreateUser(user model.User) {
	existingUser := repo.FindByEmail(user.Email)
	if existingUser.ID != 0 {
		panic("User with this email already exists")
	}

	repo.Connection.Open()
	defer repo.Connection.Close()
	repo.Connection.Conn.Create(&user)
}

func (repo *User) FindByID(id int) model.User {
	repo.Connection.Open()
	defer repo.Connection.Close()

	user := model.User{}
	repo.Connection.Conn.First(&user, id)

	return user
}

func (repo *User) FindByEmail(email string) model.User {
	repo.Connection.Open()
	defer repo.Connection.Close()

	user := model.User{}
	repo.Connection.Conn.First(&user, "email = ?", email)

	return user
}

func (repo *User) IsAPIKeyValidForUser(userID int, key string) bool {
	repo.Connection.Open()
	defer repo.Connection.Close()

	k := model.ApiKey{}
	repo.Connection.Conn.First(&k, "user_id = ? AND key = ?", userID, key)
	if k.ID == 0 {
		return false
	}
	if !k.Enabled {
		return false
	}

	return true
}