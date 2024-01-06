package models

import "time"

type Question struct {
	ID           uint64    `db:"id"`
	QuestionText string    `db:"question_text"`
	Answer       string    `db:"answer"`
	Date         time.Time `db:"answered_at"`
	FromChatId   int64     `db:"from_chat_id"`
}
