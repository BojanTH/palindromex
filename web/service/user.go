package service

import (
	"palindromex/web/dto"
	"palindromex/web/model"
	"palindromex/web/repository"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserRepository *repository.User
}

func NewUser(UserRepository *repository.User) *User {
	return &User{UserRepository: UserRepository}
}

func (service *User) GetUserByID(id int) model.User {
	return service.UserRepository.FindByID(id)
}

func (service *User) CreateNewUser(credentials *dto.Credentials) {
	user := model.User{}
	user.Name = credentials.Name
	user.Email = credentials.Email
	user.Password = service.getHashedPassword(credentials.Password)
	// Set the user to enabled, in real life this should probably happen when the user confirms the signup email
	user.Enabled = true

	service.UserRepository.CreateUser(user)
}

func (service *User) getHashedPassword(password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		panic(err)
	}

	return hashedPassword
}

func (service *User) GetUserByEmailAndPassword(email string, password string) (model.User, error) {
	user := service.UserRepository.FindByEmail(email)
	if user.ID == 0 {
		return user, errors.New("User with this email doesn't exist")
	}
	if nil == bcrypt.CompareHashAndPassword(user.Password, []byte(password)) {
		return user, nil
	}
	if !user.Enabled {
		return user, errors.New("This user is not enabled")
	}

	return user, errors.New("Invalid password")
}

func (service *User) IsAPIKeyValidForUser(userID int, key string) bool {
	return service.UserRepository.IsAPIKeyValidForUser(userID, key)
}