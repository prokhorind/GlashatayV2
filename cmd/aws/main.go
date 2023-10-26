package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prokhorind/glashatayBotV2/internal/handler"
	tb "gopkg.in/tucnak/telebot.v2"
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

	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var u tb.Update
		if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
			tgBot.ProcessUpdate(u)
		}
		return events.APIGatewayProxyResponse{Body: "ok", StatusCode: 200}, nil
	})
}
