package repository

import (
	"database/sql"
	Usdomain "github.com/prokhorind/glashatayBotV2/internal/core/domain"
)

type phraseRepository struct {
	db *sql.DB
}

func NewPhraseRepository(db *sql.DB) *phraseRepository {
	return &phraseRepository{db: db}
}

func (rep phraseRepository) SelectRandomPhrase(chat Usdomain.Chat) (*Usdomain.Phrase, error) {
	sqlGameStatement := `SELECT id, user_id ,phrase , type  FROM phrases WHERE user_id IN(SELECT user_id from game_stats where chat_id = $1 GROUP BY user_id)
ORDER BY RANDOM()
LIMIT 1`
	row := rep.db.QueryRow(sqlGameStatement, chat.ID)

	var phrase Usdomain.Phrase
	err := row.Scan(&phrase.Id, &phrase.UserID, &phrase.Phrase, &phrase.Type)

	return &phrase, err
}
