package repository

import (
	"context"
	"pauuser/mood-bot/internal/models"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *models.Question) (uint64, error)
	GetAll(ctx context.Context) ([]*models.Question, error)
}
