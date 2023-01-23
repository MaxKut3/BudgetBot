package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

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

		if matched, _ := regexp.MatchString("[A-z]", wordList[1]); matched != true {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невалидная строка. Строка должна быть следующего вида: Продукты 1000Rub")
			if _, sendError := bot.Send(msg); sendError != nil {
				log.Println(fmt.Errorf("send message failed: %v", sendError))
			}
			continue
		}

		category, sum, cur := stringParser(wordList)

		// Отправка ответного сообщения
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s, сумма вашей покупки составила: %s, в следующей валюте: %s. Категория покупки - %s ", update.Message.From.UserName, sum, cur, category))
		if _, sendError := bot.Send(msg); sendError != nil {
			log.Panic(fmt.Errorf("send message failed: %v", sendError))
		}
	}
}

func stringParser(str []string) (category, sum, cur string) {

	re, err := regexp.Compile("[A-z]")
	if err != nil {
		log.Println(err)
	}

	listInd := re.FindStringIndex(str[1])
	i := listInd[0]

	category = str[0]
	sum = str[1][:i]
	cur = str[1][i:]

	return category, sum, cur
}

func fixerAPI(cur, sum string) int {

	URL := fmt.Sprintf("https://api.apilayer.com/fixer/convert?to=RUB&from=%s&amount=%s", cur, sum)

	client := &http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("apikey", "rd8XN7uKMiFl6k9fXUMFt3TqEsF1EJhb")

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)

	fmt.Println(string(body))
	fmt.Println(body)

	var fixerJSON FixerJSON

	unmarshalErr := json.Unmarshal(body, &fixerJSON)
	if unmarshalErr != nil {
		log.Println(unmarshalErr)
	}
	fmt.Println(fixerJSON)
	fmt.Printf("Результат: %d", int(fixerJSON.Result))
	fmt.Println()

	return int(fixerJSON.Result)
}
