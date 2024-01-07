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
    	question_text TEXT,
    	answer TEXT,
    	answered_at TEXT,
    	from_chat_id INTEGER
	)`
	_, err := db.Exec(createQuestionsTable)

	return err
}

func createUsersTable(db *sql.DB) error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	chat_id INTEGER,
	    name TEXT,
	    username TEXT
	)`
	_, err := db.Exec(createUsersTable)

	return err
}
