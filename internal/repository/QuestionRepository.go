package repository

import (
	"pauuser/mood-bot/internal/models"
)

type QuestionRepository interface {
	Create(question *models.Question) error
	GetAll() ([]*models.Question, error)
}
