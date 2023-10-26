package handler

import (
	"github.com/prokhorind/glashatayBotV2/internal/core/domain"
	tb "gopkg.in/tucnak/telebot.v2"
)

type userConverter struct{}

func (c userConverter) convert(user *tb.User) domain.User {
	return domain.User{
		ID: user.ID,
	}

}
