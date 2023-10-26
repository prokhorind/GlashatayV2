package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prokhorind/glashatayBotV2/internal/core/domain"
	"github.com/prokhorind/glashatayBotV2/internal/core/services"
	"github.com/prokhorind/glashatayBotV2/internal/handler/view"
	"github.com/prokhorind/glashatayBotV2/internal/locale"
	"github.com/prokhorind/glashatayBotV2/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"time"
)

var handlers map[string]func(bot *tgbotapi.BotAPI, update *tb.Message)

func init() {
	config := repository.NewDBConfig()
	gameRepository := repository.NewGameRepository(config)
	phraseRepository := repository.NewPhraseRepository(config)
	gameService := services.NewGameService(gameRepository, phraseRepository)
	translator, _ := locale.NewTranslator()

	handlers = map[string]func(bot *tgbotapi.BotAPI, update *tb.Message){
		"/gaytoday": func(bot *tgbotapi.BotAPI, update *tb.Message) {
			chat := chatConverter{}.convert(update.Chat)

			hasRun, chatRes, err := gameService.HasJobAlreadyRun(chat)
			if err != nil {
				logrus.Errorf("can't get chat info: %d %s", chat.ID, err.Error())
				return
			}

			if hasRun {
				message := tgbotapi.MessageConfig{}
				message.ChatID = update.Chat.ID
				viewUserName := view.GetUserName(bot, update, chatRes.SelectedUserId.Int64, true)

				text, _ := translator.Get("gayAlreadyChosen", language.Ukrainian)
				message.Text = fmt.Sprintf(text, viewUserName)
				bot.Send(message)
				return
			}

			user, phrase, error := gameService.RunGame(chat)
			if error != nil {
				logrus.Errorf("can't run game for chat %s with error: %s", chat.ID, error.Error())
				return
			}
			viewUserName := view.GetUserName(bot, update, user.ID, true)
			var text string
			if phrase.Type == "DYNAMIC" {
				text = strings.ReplaceAll(phrase.Phrase, "%gayname%", viewUserName)
			} else {
				text = phrase.Phrase
			}

			textLines := strings.Split(text, "&")

			if phrase.Type == "COMMON" {
				msgEnding, _ := translator.Get("commonMsgEnding", language.Ukrainian)
				textLines = append(textLines, fmt.Sprintf(msgEnding, viewUserName))
			}

			for _, line := range textLines {
				message := tgbotapi.MessageConfig{}
				message.ChatID = update.Chat.ID
				message.Text = line
				bot.Send(message)
				time.Sleep(1 * time.Second)
			}
		},
		"/gayreg": func(bot *tgbotapi.BotAPI, update *tb.Message) {

			user := userConverter{}.convert(update.Sender)
			chat := chatConverter{}.convert(update.Chat)

			gameService.Register(user, chat, time.Now().Year())
			message := tgbotapi.MessageConfig{}
			message.ChatID = update.Chat.ID
			text, _ := translator.Get("regMsg", language.Ukrainian)
			message.Text = text
			bot.Send(message)
		},

		"/stat": func(bot *tgbotapi.BotAPI, update *tb.Message) {

			chat := chatConverter{}.convert(update.Chat)
			params := strings.Split(update.Text, " ")

			var stats []domain.GameStat
			var statsError error
			sb := strings.Builder{}
			if len(params) == 1 {
				header, _ := translator.Get("yearStatHeader", language.Ukrainian)
				sb.WriteString(fmt.Sprintf(header, time.Now().Year()))
				stats, statsError = gameService.GetStatByYear(chat, time.Now().Year())
			} else {
				if "ALL" == params[1] {
					stats, statsError = gameService.GetStat(chat)
					header, _ := translator.Get("allStatHeader", language.Ukrainian)
					sb.WriteString(header)
				} else {
					year, err := strconv.Atoi(params[1])
					if err != nil {
						message := tgbotapi.MessageConfig{}
						message.ChatID = update.Chat.ID
						text, _ := translator.Get("wrongStatWithYearRequest", language.Ukrainian)
						message.Text = text
						bot.Send(message)
						return
					}
					if year < 2019 || year > time.Now().Year() {
						year = 2019

					}
					header, _ := translator.Get("yearStatHeader", language.Ukrainian)
					sb.WriteString(fmt.Sprintf(header, year))
					stats, statsError = gameService.GetStatByYear(chat, year)
				}

			}

			if len(stats) == 0 || statsError != nil {
				message := tgbotapi.MessageConfig{}
				message.ChatID = update.Chat.ID
				text, _ := translator.Get("noStatsMsg", language.Ukrainian)
				message.Text = text
				bot.Send(message)
				return
			}

			for id, stat := range stats {
				stat := stat
				viewUserName := view.GetUserName(bot, update, stat.UserID, false)
				sb.WriteString(fmt.Sprintf("\n <strong>%d.</strong> %s: %d", id+1, viewUserName, stat.Num))
			}

			message := tgbotapi.MessageConfig{}
			message.ChatID = update.Chat.ID
			message.Text = sb.String()
			message.ParseMode = "HTML"
			bot.Send(message)
		},
	}
}

func GetHandler(command string) func(bot *tgbotapi.BotAPI, update *tb.Message) {

	if v, ok := handlers[command]; ok {
		return v
	} else {
		return func(bot *tgbotapi.BotAPI, update *tb.Message) {}
	}
}
