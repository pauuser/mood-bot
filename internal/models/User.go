package models

type User struct {
	ID       uint64 `db:"id"`
	ChatId   int64  `db:"chat_id"`
	Name     string `db:"name"`
	Username string `db:"username"`
}
