package domain

type GameStat struct {
	ID     int
	UserID int64
	ChatID int64
	Num    int
	Year   int
}

func (m *GameStat) TableName() string {
	return "game_stats"
}
