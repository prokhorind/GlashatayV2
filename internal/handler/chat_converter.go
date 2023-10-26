package handler

import (
	"github.com/prokhorind/glashatayBotV2/internal/core/domain"
	tb "gopkg.in/tucnak/telebot.v2"
)

type chatConverter struct{}

func (c chatConverter) convert(chat *tb.Chat) domain.Chat {
	return domain.Chat{
		Name: &chat.Title,
		ID:   chat.ID,
	}
}
