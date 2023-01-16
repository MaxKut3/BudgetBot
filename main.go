package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("5859735143:AAGFvrZyYDwF7RlOkrca4Ekfx8rPpJJQU9k")
	if err != nil {
		log.Panic(fmt.Errorf("authorization failed", err))
		panic(err)
	}

	log.Println("Authorization was successful")

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		res, parsError := stringParsing(update.Message.Text)
		if parsError != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невалидная строка")
			if _, sendError := bot.Send(msg); sendError != nil {
				fmt.Errorf("Send message failed", sendError)
			}
			panic(parsError)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s, сумма вашей покупки составила: %s, в следующей валюте: %s ", update.Message.From.UserName, res[0], res[1]))
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if _, sendError := bot.Send(msg); sendError != nil {
			log.Panic(fmt.Errorf("Send message failed", sendError))
			panic(sendError)
		}
	}
}

// Извлечение суммы, валюты, продукта из сообщения
func stringParsing(text string) ([]string, error) {

	msgList := strings.Split(text, " ")
	re, _ := regexp.Compile("[a-z]")

	//Проверка на валидность
	if matched, matchedError := regexp.MatchString("[a-z]", msgList[len(msgList)-1]); matched != true {
		return nil, fmt.Errorf("Невалидная строка", matchedError)
	}

	listInd := re.FindStringIndex(msgList[len(msgList)-1])
	fmt.Println(listInd)

	i := listInd[0]

	res := make([]string, 2, 2)
	res[0] = msgList[len(msgList)-1][:i] //сумма
	res[1] = msgList[len(msgList)-1][i:] //валюта

	return res, nil

}
