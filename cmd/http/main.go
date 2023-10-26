package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prokhorind/glashatayBotV2/internal/handler"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"os"
	"strings"
)

func main() {

	settings := tb.Settings{
		Token:       os.Getenv("BOT_TOKEN"),
		Synchronous: true,
		Verbose:     true,
	}
	api, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	tgBot, err := tb.NewBot(settings)
	if err != nil {
		fmt.Println(err)
		panic("can't create bot")
	}
	tgBot.Handle(tb.OnText, func(m *tb.Message) {
		hn := handler.GetHandler(strings.Split(m.Text, " ")[0])
		hn(api, m)
	})

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		var u tb.Update
		err := json.NewDecoder(req.Body).Decode(&u)
		if err == nil {
			tgBot.ProcessUpdate(u)
		}
	})

	http.ListenAndServe("0.0.0.0:8081", nil)
}
