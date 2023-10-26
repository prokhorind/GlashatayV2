package repository

import (
	"database/sql"
	Usdomain "github.com/prokhorind/glashatayBotV2/internal/core/domain"
	"github.com/sirupsen/logrus"
)

type gameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *gameRepository {
	return &gameRepository{db: db}
}

func (rep gameRepository) CreateNewUser(user Usdomain.User) {
	sqlUserStatement := `INSERT INTO users (id, active) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`
	_, err := rep.db.Exec(sqlUserStatement, user.ID, true)
	if err != nil {
		logrus.Error("Can't create user")
	}
}

func (rep gameRepository) CreateNewChat(chat Usdomain.Chat) {
	sqlChatStatement := `INSERT INTO chats (id,name, auto_pick, all_phrases) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`
	_, err := rep.db.Exec(sqlChatStatement, chat.ID, chat.Name, true, false)
	if err != nil {
		logrus.Error("Can't create chat")
	}
}

func (rep gameRepository) CreateNewGame(user Usdomain.User, chat Usdomain.Chat, year int) {
	sqlGameStatement := `INSERT INTO game_stats (user_id, chat_id, num , year) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id,chat_id,year) DO NOTHING`
	_, err := rep.db.Exec(sqlGameStatement, user.ID, chat.ID, 0, year)
	if err != nil {
		logrus.Error("Can't create new game")
	}
}

func (rep gameRepository) SelectWinnerByChat(chat Usdomain.Chat) (*Usdomain.User, error) {
	sqlGameStatement := `SELECT id, active FROM users 
JOIN ( SELECT user_id FROM game_stats WHERE chat_id = $1 GROUP BY user_id , chat_id ) stats 
ON stats.user_id = users.id
WHERE active = true 
ORDER BY RANDOM()
LIMIT 1`
	row := rep.db.QueryRow(sqlGameStatement, chat.ID)

	var user Usdomain.User
	err := row.Scan(&user.ID, &user.Active)

	return &user, err
}

func (rep gameRepository) SelectUserStatByYear(user Usdomain.User, chat Usdomain.Chat, year int) (*Usdomain.GameStat, error) {
	sqlGameStatement := `SELECT id, user_id , chat_id , num , year  FROM game_stats WHERE chat_id = $1 AND user_id = $2 AND year = $3`

	row := rep.db.QueryRow(sqlGameStatement, chat.ID, user.ID, year)

	var stat Usdomain.GameStat
	err := row.Scan(&stat.ID, &stat.UserID, &stat.ChatID, &stat.Num, &stat.Year)

	return &stat, err
}

func (rep gameRepository) UpdateUserStat(user Usdomain.User, chat Usdomain.Chat, year int, val int) error {
	sqlGameStatement := `UPDATE game_stats SET num = $1 WHERE chat_id = $2 AND user_id = $3 AND year = $4`

	_, err := rep.db.Exec(sqlGameStatement, val, chat.ID, user.ID, year)

	return err
}

func (rep gameRepository) UpdateChatTriggeredTime(chat Usdomain.Chat, userId int64) error {
	sqlGameStatement := `UPDATE chats SET last_time_triggered = $1 , selected_user_id = $2 WHERE id = $3`

	_, err := rep.db.Exec(sqlGameStatement, chat.LastTimeTriggered.Time, userId, chat.ID)

	return err
}

func (rep gameRepository) SelectChatById(chatId int64) (*Usdomain.Chat, error) {
	sqlGameStatement := `SELECT id, name, auto_pick , all_phrases , last_time_triggered, selected_user_id  FROM chats WHERE id = $1`

	row := rep.db.QueryRow(sqlGameStatement, chatId)

	var chat Usdomain.Chat
	err := row.Scan(&chat.ID, &chat.Name, &chat.AutoPick, &chat.AllPhrases, &chat.LastTimeTriggered, &chat.SelectedUserId)

	return &chat, err
}

func (rep gameRepository) SelectChatStatsByIdAndYear(chatId int64, year int) ([]Usdomain.GameStat, error) {
	sqlGameStatement := `SELECT user_id ,SUM(num)  FROM game_stats as g 
WHERE chat_id = $1 AND year = $2
GROUP BY user_id
ORDER BY sum(num) DESC
LIMIT 10`

	rows, err := rep.db.Query(sqlGameStatement, chatId, year)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	stats := make([]Usdomain.GameStat, 0)
	for rows.Next() {
		var stat Usdomain.GameStat
		rows.Scan(&stat.UserID, &stat.Num)
		stats = append(stats, stat)
	}
	return stats, err
}

func (rep gameRepository) SelectChatStatsById(chatId int64) ([]Usdomain.GameStat, error) {
	sqlGameStatement := `SELECT user_id ,SUM(num)  FROM game_stats as g 
WHERE chat_id = $1
GROUP BY user_id
ORDER BY sum(num) DESC
LIMIT 10`

	rows, err := rep.db.Query(sqlGameStatement, chatId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	stats := make([]Usdomain.GameStat, 0)
	for rows.Next() {
		var stat Usdomain.GameStat
		rows.Scan(&stat.UserID, &stat.Num)
		stats = append(stats, stat)
	}
	return stats, err
}
