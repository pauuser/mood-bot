package repository_impl

import (
	"database/sql"
	"go.uber.org/zap"
	"pauuser/mood-bot/internal/models"
	"pauuser/mood-bot/internal/repository"
)

type userRepositoryImpl struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepoSqliteImpl(db *sql.DB, logger *zap.Logger) repository.UserRepository {
	return &userRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (u userRepositoryImpl) Create(user *models.User) error {
	createQuery := `INSERT INTO users (chat_id, name, username)
					VALUES ($1, $2, $3)`

	_, err := u.db.Exec(createQuery, user.ChatId, user.Name, user.Username)
	if err != nil {
		u.logger.Error("insert into questions error", zap.Error(err))
	}

	return err
}

func (u userRepositoryImpl) GetAll() ([]*models.User, error) {
	getQuery := `SELECT * FROM users`

	var users = make([]*models.User, 0)
	result, err := u.db.Query(getQuery)
	if err != nil {
		u.logger.Error("Could not query users")
	}
	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			u.logger.Error("Error closing rows")
		}
	}(result)

	for result.Next() {
		user := new(models.User)
		if err := result.Scan(&user.ID, &user.ChatId, &user.Name, &user.Username); err != nil {
			u.logger.Error("Could not query users")
		}
		users = append(users, user)
	}

	return users, nil
}

func (u userRepositoryImpl) GetUser(chatId int64) (*models.User, error) {
	getQuery := `SELECT * FROM users WHERE chat_id = $1`
	user := new(models.User)

	result := u.db.QueryRow(getQuery, chatId)
	if err := result.Scan(&user.ID, &user.ChatId, &user.Name, &user.Username); err != nil {
		u.logger.Error("Could not query user")
		return nil, err
	}
	if user.ChatId == 0 {
		return nil, nil
	}

	return user, nil
}
