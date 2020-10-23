package service

import (
	"palindromex/web/db"
	"palindromex/web/dto"
	"palindromex/web/model"
	"palindromex/web/repository"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Connection     *db.Connection
	UserRepository *repository.User
}

func NewUser(connection *db.Connection, UserRepository *repository.User) *User {
	return &User{Connection: connection, UserRepository: UserRepository}
}

// @TODO move queries to the repository
func (userService *User) CreateNewUser(credentials *dto.Credentials) {
	user := model.User{}
	user.Name = credentials.Name
	user.Email = credentials.Email
	user.Password = userService.getHashedPassword(credentials.Password)
	// Set the user to enabled, in real life this should probably happen when the user confirms the signup email
	user.Enabled = true

	userService.Connection.Open()
	defer userService.Connection.Close()

	existingUser := model.User{}
	userService.Connection.Conn.First(&existingUser, "email = ?", user.Email)
	if existingUser.ID != 0 {
		panic("User with this email already exists")
	}

	userService.Connection.Conn.Create(&user)
}

func (userService *User) getHashedPassword(password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		panic(err)
	}

	return hashedPassword
}

func (userService *User) GetUserByEmailAndPassword(email string, password string) (model.User, error) {
	userService.Connection.Open()
	defer userService.Connection.Close()

	u := model.User{}
	userService.Connection.Conn.First(&u, "email = ?", email)
	if u.ID == 0 {
		return u, errors.New("User with this email doesn't exist")
	}
	if !u.Enabled {
		return u, errors.New("This user is not enabled")
	}

	if nil == bcrypt.CompareHashAndPassword(u.Password, []byte(password)) {
		return u, nil
	}

	return u, errors.New("Invalid password")
}

func (userService *User) IsAPIKeyValidForUser(userID int, key string) bool {
	userService.Connection.Open()
	defer userService.Connection.Close()

	k := model.ApiKey{}
	userService.Connection.Conn.First(&k, "user_id = ? AND key = ?", userID, key)
	if k.ID == 0 {
		return false
	}
	if !k.Enabled {
		return false
	}

	return true
}