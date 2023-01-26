package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
)

func main() {

	// Блок подключения к Боту

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	val := os.Getenv("KEY")
	fixer := os.Getenv("FIXER")
	exchangerates := os.Getenv("EXCHANGERATES")

	bot, err := tgbotapi.NewBotAPI(val)

	if err != nil {
		log.Panic(fmt.Errorf("authorization failed: %v", err))
		panic(err)
	}

	log.Println("Authorization was successful")

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// Блок обработки сообщений

	for update := range updates {

		// Проверка строки на валидность
		wordList := strings.Split(update.Message.Text, " ")
		fmt.Println(wordList)

		if len(wordList) != 2 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невалидная строка. Строка должна быть следующего вида: Продукты 1000Rub")
			if _, sendError := bot.Send(msg); sendError != nil {
				log.Println(fmt.Errorf("send message failed: %v", sendError))
			}
			continue
		}

		if matched, _ := regexp.MatchString("[0-9]", wordList[1]); matched != true {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невалидная строка. Строка должна быть следующего вида: Продукты 1000Rub")
			if _, sendError := bot.Send(msg); sendError != nil {
				log.Println(fmt.Errorf("send message failed: %v", sendError))
			}
			continue
		}

		if matched, _ := regexp.MatchString("[A-z]", wordList[1]); matched != true {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невалидная строка. Строка должна быть следующего вида: Продукты 1000Rub")
			if _, sendError := bot.Send(msg); sendError != nil {
				log.Println(fmt.Errorf("send message failed: %v", sendError))
			}
			continue
		}

		category, sum, cur := stringParser(wordList)

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			fmt.Println("fixer", fixerAPI(cur, fixer))
		}()

		go func() {
			defer wg.Done()
			fmt.Println("coinGAteAPI", coinGAteAPI(cur))
		}()

		go func() {
			defer wg.Done()
			fmt.Println("exchangeratesAPI", exchangeratesAPI(cur, exchangerates))
		}()

		wg.Wait()

		// Отправка ответного сообщения
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s, сумма вашей покупки составила: %s, в следующей валюте: %s. Категория покупки - %s ", update.Message.From.UserName, sum, cur, category))
		if _, sendError := bot.Send(msg); sendError != nil {
			log.Panic(fmt.Errorf("send message failed: %v", sendError))
		}
	}
}

var re, _ = regexp.Compile("[A-z]")

func stringParser(str []string) (category, sum, cur string) {

	listInd := re.FindStringIndex(str[1])
	i := listInd[0]

	category = str[0]
	sum = str[1][:i]
	cur = str[1][i:]

	return category, sum, cur
}
