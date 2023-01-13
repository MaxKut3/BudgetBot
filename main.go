package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("5859735143:AAGFvrZyYDwF7RlOkrca4Ekfx8rPpJJQU9k")
	if err != nil {
		log.Panic(fmt.Errorf("Authorization failed", err))
		panic(err)
	}

	log.Println("Authorization was successful")

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if _, sendError := bot.Send(msg); sendError != nil {
			log.Panic(fmt.Errorf("Send message failed", sendError))
			panic(sendError)
		}
	}
}


