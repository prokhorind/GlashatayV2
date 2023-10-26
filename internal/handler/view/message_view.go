package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

func GetUserName(bot *tgbotapi.BotAPI, update *tb.Message, userId int64) string {
	member := tgbotapi.GetChatMemberConfig{}
	member.ChatID = update.Chat.ID
	member.UserID = userId
	mm, err := bot.GetChatMember(member)
	if err != nil {
		logrus.Errorf("Can`t obtain user Info %d", userId)
		return ""
	}

	if len(strings.TrimSpace(mm.User.UserName)) != 0 {
		return "@" + mm.User.UserName
	}

	return mm.User.FirstName + " " + mm.User.LastName

}
