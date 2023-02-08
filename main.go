package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
)

var cache = make(map[string]int)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}

func main() {

	val := os.Getenv("KEY")

	bot, err := tgbotapi.NewBotAPI(val)

	if err != nil {
		log.Panic(fmt.Errorf("authorization failed: %v", err))
	}

	log.Println("Authorization was successful")

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	go func() {

		for update := range updates {

			wordList := strings.Split(update.Message.Text, " ")

			if validateMessageText(wordList) != true {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невалидная строка. Строка должна быть следующего вида: Продукты 1000Rub")
				if _, sendError := bot.Send(msg); sendError != nil {
					log.Println(fmt.Errorf("send message failed: %v", sendError))
				}
				continue
			}

			category, sum, cur := stringParser(wordList)
			sumInt, _ := strconv.Atoi(sum)
			sumRub := 0

			if _, ok := cache[cur]; !ok {
				cache[cur] = getValue(cur)
			}

			sumRub = cache[cur] * sumInt / 100

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s, сумма вашей покупки составила: %s, в следующей валюте: %s. Сумма в рублях - %d . Категория покупки - %s ", update.Message.From.UserName, sum, cur, sumRub, category))
			if _, sendError := bot.Send(msg); sendError != nil {
				log.Panic(fmt.Errorf("send message failed: %v", sendError))
			}
		}

	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan

	log.Printf("%s signal caught", sig)

	bot.StopReceivingUpdates()
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

func validateMessageText(wordList []string) bool {
	if len(wordList) != 2 {
		return false
	}
	if matched, _ := regexp.MatchString("[0-9]", wordList[1]); matched != true {
		return false
	}
	if matched, _ := regexp.MatchString("[A-z]", wordList[1]); matched != true {
		return false
	}
	return true
}

type getCurValueFunc func(cur, key string) int

func getValue(cur string) int {

	providers := map[string]getCurValueFunc{
		"No key":                   coinGAteAPI,
		os.Getenv("FIXER"):         fixerAPI,
		os.Getenv("EXCHANGERATES"): exchangeratesAPI,
	}

	var wg sync.WaitGroup
	wg.Add(len(providers))

	ch := make(chan int)

	max := 0

	for key, provider := range providers {
		key := key
		go func() {
			defer wg.Done()
			ch <- provider(cur, key)
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for val := range ch {
		if val > max {
			max = val
		}
	}
	return max
}
