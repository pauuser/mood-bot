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

	init := `
	CREATE TABLE IF NOT EXISTS questions (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	question_text TEXT,
    	answer TEXT,
    	answered_at TEXT,
    	from_chat_id INTEGER
)`
	_, err = db.Exec(init)
	if err != nil {
		return nil, err
	}

	return db, nil
}
