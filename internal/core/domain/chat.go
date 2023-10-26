package domain

import (
	"database/sql"
)

type Chat struct {
	ID                int64
	Name              *string
	AutoPick          bool
	AllPhrases        bool
	LastTimeTriggered sql.NullTime
	SelectedUserId    sql.NullInt64
}

func (m *Chat) TableName() string {
	return "chats"
}
