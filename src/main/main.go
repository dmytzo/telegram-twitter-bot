package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"os"
)




var port = os.Getenv("PORT")
var cKey = os.Getenv("cKey")
var cSecret = os.Getenv("cSecret")
var t = os.Getenv("t")
var tSecret = os.Getenv("tSecret")
var botApi = os.Getenv("botApi")
var WebHookURL = os.Getenv("webHookUrl")

func setUpTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(cKey, cSecret)
	token := oauth1.NewToken(t, tSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client
}

func setUpTelegramBot() *tgbotapi.BotAPI{
	bot, err := tgbotapi.NewBotAPI(botApi)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebHookURL))
	if err != nil {
		log.Fatal(err)
	}

	return bot
}

func main() {

	searchButtons := []tgbotapi.KeyboardButton{
		tgbotapi.KeyboardButton{Text: "Search in Twitter"},
	}

	searchOptions := map[string]bool{
		"Search in Twitter": false,
	}

	twitterClient := setUpTwitterClient()
	bot := setUpTelegramBot()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+port, nil)

	for update := range updates {
		var message tgbotapi.MessageConfig
		msg := update.Message.Text

		switch msg {

		case "Search in Twitter":
			message = tgbotapi.NewMessage(update.Message.Chat.ID, "<Type keyword to search>")
			message.ReplyMarkup = tgbotapi.NewReplyKeyboard(searchButtons)
			bot.Send(message)

			searchOptions["Search in Twitter"] = true

		default:

			if searchOptions["Search in Twitter"] {
				query := update.Message.Text
				search, _, _ := twitterClient.Search.Tweets(&twitter.SearchTweetParams{
					Query: query, Count: 10, ResultType: "recent",
				})

				if search != nil {
					message = tgbotapi.NewMessage(update.Message.Chat.ID, "********** \n RESULTS: \n**********")
					message.ReplyMarkup = tgbotapi.NewReplyKeyboard(searchButtons)
					bot.Send(message)

					for _, tweet := range search.Statuses {
						message = tgbotapi.NewMessage(update.Message.Chat.ID, "Text: " + tweet.Text + "\n" + " Link: https://twitter.com/statuses/" + tweet.IDStr)
						message.ReplyMarkup = tgbotapi.NewReplyKeyboard(searchButtons)
						bot.Send(message)
					}
				}
				searchOptions["Search in Twitter"] = false

			} else {
				message = tgbotapi.NewMessage(update.Message.Chat.ID, "Press button to search")
				message.ReplyMarkup = tgbotapi.NewReplyKeyboard(searchButtons)
				bot.Send(message)
			}
		}
	}
}
