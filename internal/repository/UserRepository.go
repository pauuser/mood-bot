package repository

import (
	"pauuser/mood-bot/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetAll() ([]*models.User, error)
	GetUser(chatId int64) (*models.User, error)
}
