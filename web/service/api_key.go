package service

import (
	"errors"
	"palindromex/web/model"
	"palindromex/web/repository"
)

type APIKey struct {
	APIKeyRepository *repository.APIKey
}

func NewAPIKey(APIKeyRepository *repository.APIKey) *APIKey {
	return &APIKey{APIKeyRepository: APIKeyRepository}
}

func (apiKeyService *APIKey) CreateNew(user model.User, apiKey string) error {
	newKey := model.ApiKey{
		UserID: user.ID,
		Key: apiKey,
		Enabled: true,
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

	apiKeyService.APIKeyRepository.CreateNew(newKey)

	return nil
}
