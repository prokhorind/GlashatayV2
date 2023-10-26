package services

import (
	"database/sql"
	"errors"
	Usdomain "github.com/prokhorind/glashatayBotV2/internal/core/domain"
	"github.com/prokhorind/glashatayBotV2/internal/core/ports"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type gameService struct {
	gameRepository   ports.GameRepository
	phraseRepository ports.PhraseRepository
}

func NewGameService(gameRepository ports.GameRepository, phraseRepository ports.PhraseRepository) *gameService {
	return &gameService{
		gameRepository:   gameRepository,
		phraseRepository: phraseRepository,
	}
}

func (s gameService) Register(user Usdomain.User, chat Usdomain.Chat, year int) {
	s.gameRepository.CreateNewUser(user)
	s.gameRepository.CreateNewChat(chat)
	s.gameRepository.CreateNewGame(user, chat, year)
}

func (s gameService) RunGame(chat Usdomain.Chat) (*Usdomain.User, *Usdomain.Phrase, error) {

	resp, err := s.gameRepository.SelectChatById(chat.ID)
	if err != nil {
		return nil, nil, err
	}

	loc, _ := time.LoadLocation(os.Getenv("TZ"))
	now := time.Now().In(loc)
	lt := resp.LastTimeTriggered
	yr, mth, day := now.Date()
	if lt.Valid {
		tyr, tmth, tday := lt.Time.In(loc).Date()

		if tyr == yr && tmth == mth && tday == day {
			logrus.Info("job is already triggered")
			return nil, nil, errors.New("job is already triggered")
		}
	}

	user, err := s.gameRepository.SelectWinnerByChat(chat)
	if err != nil {
		return nil, nil, err
	}

	stat, err := s.gameRepository.SelectUserStatByYear(*user, chat, yr)
	s.gameRepository.SelectUserStatByYear(*user, chat, yr)

	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	if err != nil && err == sql.ErrNoRows {
		s.gameRepository.CreateNewGame(*user, chat, yr)
	}

	err = s.gameRepository.UpdateUserStat(*user, chat, yr, stat.Num+1)

	if err != nil {
		return nil, nil, err
	}
	chat.LastTimeTriggered = sql.NullTime{Time: now.UTC()}
	s.gameRepository.UpdateChatTriggeredTime(chat, user.ID)

	phrase, err := s.phraseRepository.SelectRandomPhrase(chat)

	if err != nil {
		return nil, nil, err
	}

	return user, phrase, err
}

func (s gameService) GetStatByYear(chat Usdomain.Chat, year int) ([]Usdomain.GameStat, error) {
	return s.gameRepository.SelectChatStatsByIdAndYear(chat.ID, year)
}

func (s gameService) GetStat(chat Usdomain.Chat) ([]Usdomain.GameStat, error) {
	return s.gameRepository.SelectChatStatsById(chat.ID)
}
