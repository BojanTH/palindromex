package service

import (
	"errors"
	"palindromex/web/dto"
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

func (service *Message) UpdateMessage(userID int, messageID int, content string) error {
	message, err := service.MessageRepository.FindMessage(userID, messageID)
	if err != nil {
		return err
	}

	message.Content = content
	message.Palindrome = IsPalindrome(content)

	return service.MessageRepository.UpdateMessage(message)
}

func (service *Message) DeleteMessage(userID, messageID int) error {
	return service.MessageRepository.DeleteMessage(userID, messageID)
}

func (service *Message) FindAllByUserID(userID int) ([]dto.Message, error) {
	messages := []dto.Message{}
	allMessages, err := service.MessageRepository.FindAllByUserID(userID)
	if err != nil {
		return messages, err
	}

	for _, message := range allMessages {
	    messageDTO := dto.Message {
			ID: message.ID,
			UserID: message.UserID,
	        Content: message.Content,
	        Palindrome: message.Palindrome,
	        CreatedAt: message.CreatedAt,
	        UpdatedAt: message.UpdatedAt,
	    }
		messages = append(messages, messageDTO)
	}

	return messages, nil
}

func (service *Message) FindMessage(userID, messageID int) (dto.Message, error) {
	message, err := service.MessageRepository.FindMessage(userID, messageID)
	var messageDTO dto.Message
	if err != nil {
		return messageDTO, err
	}

	// Return DTO instad of the model to hide associated user data
	messageDTO = dto.Message {
		ID: message.ID,
		UserID: message.UserID,
	    Content: message.Content,
	    Palindrome: message.Palindrome,
	    CreatedAt: message.CreatedAt,
	    UpdatedAt: message.UpdatedAt,
	}

	return messageDTO, nil
}