package service

import (
	"errors"
	"palindromex/web/model"
	"palindromex/web/repository"
)

type Message struct {
	MessageRepository *repository.Message
}

func NewMessage(repo *repository.Message) *Message {
	return &Message{repo}
}

func (service *Message) CreateNewMessage(user model.User, content string) error {
	isPalindrome := IsPalindrome(content)
	message := model.Message {
		UserID: user.ID,
		Content: content,
		Palindrome: isPalindrome,
	}

	defer func() error {
		if r := recover(); r != nil {
			message := "Error while inserting new API key into the database. "

			switch x := r.(type) {
			case string:
				message += x
			case error:
				message += x.Error()
			default:
				message += "Unknown error"
			}

			return errors.New(message)
		}

		return nil
	}()

	service.MessageRepository.CreateNew(message)

	return nil
}