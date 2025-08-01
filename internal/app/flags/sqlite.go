package flags

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteFlags struct {
	Path string `mapstructure:"path"`
}

func (p *SqliteFlags) InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", p.Path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createQuestionsTable(db)
	if err != nil {
		return nil, err
	}
	err = createUsersTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createQuestionsTable(db *sql.DB) error {
	createQuestionsTable := `
	CREATE TABLE IF NOT EXISTS questions (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	question_text TEXT NOT NULL,
    	answer TEXT NOT NULL,
    	answered_at DATETIME NOT NULL,
    	from_chat_id INTEGER NOT NULL,
    	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(createQuestionsTable)
	if err != nil {
		return err
	}

	return nil
}

func createUsersTable(db *sql.DB) error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	chat_id INTEGER UNIQUE NOT NULL,
	    name TEXT,
	    username TEXT,
    	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	return nil
}
