package impl

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"pauuser/mood-bot/internal/models"
	"pauuser/mood-bot/internal/repository"
)

type questionRepositoryImpl struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewQuestionRepoSqliteImpl(db *sql.DB, logger *zap.Logger) repository.QuestionRepository {
	return &questionRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (q questionRepositoryImpl) Create(ctx context.Context, question *models.Question) (uint64, error) {
	createQuery := `INSERT INTO questions (question_text, answer, answered_at, from_chat_id)
					VALUES ($1, $2, $3, $4) RETURNING ID`

	row := q.db.QueryRow(createQuery, question.QuestionText, question.Answer, question.Date, question.FromChatId)
	var id uint64

	err := row.Scan(&id)
	if err != nil {
		q.logger.Error("insert into questions error", zap.Error(err))
		return 0, err
	}

	return id, nil

}

func (q questionRepositoryImpl) GetAll(ctx context.Context) ([]*models.Question, error) {
	getQuery := `SELECT * FROM questions`

	var questions = make([]*models.Question, 0)
	result, err := q.db.Query(getQuery)
	if err != nil {
		q.logger.Error("Could not")
	}
	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			q.logger.Error("Error closing rows")
		}
	}(result)

	for result.Next() {
		question := new(models.Question)
		if err := result.Scan(&question.ID,
			&question.QuestionText,
			&question.Answer,
			&question.Date,
			&question.FromChatId); err != nil {
			panic(err)
		}
		questions = append(questions, question)
	}

	return questions, nil

}
