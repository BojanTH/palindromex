package repository

import (
	"palindromex/web/db"
	"palindromex/web/model"
)

type APIKey struct {
	Connection *db.Connection
}

func NewAPIKey(c *db.Connection) *APIKey {
	return &APIKey{c}
}

func (apiKeyRepo *APIKey) CreateNew(apiKey model.ApiKey) {
	apiKeyRepo.Connection.Open()
	defer apiKeyRepo.Connection.Close()

	apiKeyRepo.Connection.Conn.Create(&apiKey)
}
