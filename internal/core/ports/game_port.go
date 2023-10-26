package ports

import Usdomain "github.com/prokhorind/glashatayBotV2/internal/core/domain"

type GameRepository interface {
	CreateNewUser(user Usdomain.User)
	CreateNewChat(chat Usdomain.Chat)
	CreateNewGame(user Usdomain.User, chat Usdomain.Chat, year int)
	SelectWinnerByChat(chat Usdomain.Chat) (*Usdomain.User, error)
	SelectUserStatByYear(user Usdomain.User, chat Usdomain.Chat, year int) (*Usdomain.GameStat, error)
	UpdateUserStat(user Usdomain.User, chat Usdomain.Chat, year int, val int) error
	UpdateChatTriggeredTime(chat Usdomain.Chat, userId int64) error
	SelectChatById(chatId int64) (*Usdomain.Chat, error)
	SelectChatStatsByIdAndYear(chatId int64, year int) ([]Usdomain.GameStat, error)
	SelectChatStatsById(chatId int64) ([]Usdomain.GameStat, error)
}

type PhraseRepository interface {
	SelectRandomPhrase(chat Usdomain.Chat) (*Usdomain.Phrase, error)
}

type GameService interface {
	Register(user Usdomain.User, chat Usdomain.Chat, year int)
	RunGame(chat Usdomain.Chat) (*Usdomain.User, *Usdomain.Phrase, error)
	GetStatByYear(chat Usdomain.Chat, year int) ([]Usdomain.GameStat, error)
	GetStat(chat Usdomain.Chat) ([]Usdomain.GameStat, error)
	HasJobAlreadyRun(chat Usdomain.Chat) (bool, *Usdomain.Chat, error)
}
