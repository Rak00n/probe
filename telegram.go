package main

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

var bot *tgbotapi.BotAPI

func runTelegramAPI() {
	bot, _ = tgbotapi.NewBotAPI(telegramBotKey)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Command() {
			case "help":
				msg.Text = "This bot can help you monitor your infrastructure and get notifications in your telegram. First of all you must requst new probe with \"/newprobe\". Get your probes data with \"/status\". Get your Telegram ID with \"/getmyid\""
			case "start":
				msg.Text = "This bot can help you monitor your infrastructure and get notifications in your telegram. First of all you must requst new probe with \"/newprobe\". Get your probes data with \"/status\". Get your Telegram ID with \"/getmyid\""
			case "newprobe":
				uuidWithHyphen := uuid.New()
				uuidString := uuidWithHyphen.String()
				msg.Text = "New probe ID: " + uuidString
			case "status":
				msg.Text = "Your probes data"
			case "getmyid":
				msg.Text = "Your ID: " + strconv.Itoa(int(update.Message.Chat.ID))
			default:
				msg.Text = "Unknown command. Try \"/help\""
			}
			bot.Send(msg)
		} else if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func sendDataToTelegram(dataToSend string) {
	fmt.Println("Sending Data")
	dataToSendString, _ := GetAESDecrypted(dataToSend)
	msg := tgbotapi.NewMessage(421964774, string(dataToSendString))
	bot.Send(msg)
}
