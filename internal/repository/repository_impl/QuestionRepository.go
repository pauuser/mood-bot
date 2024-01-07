package repository_impl

import (
	"database/sql"
	"go.uber.org/zap"
	"pauuser/mood-bot/internal/models"
	"pauuser/mood-bot/internal/repository"
	"time"
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

func (q questionRepositoryImpl) Create(question *models.Question) error {
	createQuery := `INSERT INTO questions (question_text, answer, answered_at, from_chat_id)
					VALUES ($1, $2, $3, $4)`

	_, err := q.db.Exec(createQuery, question.QuestionText, question.Answer, question.Date, question.FromChatId)
	if err != nil {
		q.logger.Error("insert into questions error", zap.Error(err))
	}

	return err
}

func (q questionRepositoryImpl) GetAll() ([]*models.Question, error) {
	getQuery := `SELECT * FROM questions`

	var questions = make([]*models.Question, 0)
	result, err := q.db.Query(getQuery)
	if err != nil {
		q.logger.Error("Could not query questions")
	}
	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			q.logger.Error("Error closing rows")
		}
	}(result)

	for result.Next() {
		question := new(models.Question)
		var date string
		if err := result.Scan(&question.ID,
			&question.QuestionText,
			&question.Answer,
			&date,
			&question.FromChatId); err != nil {
			q.logger.Error("Could not parse answers!")
		}
		question.Date, err = time.Parse("2006-01-02 15:04:05.000000000-07:00", date)
		questions = append(questions, question)
	}

	return questions, nil
}
