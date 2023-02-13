package controller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/MaxKut3/BudgetBot/config"
	"github.com/MaxKut3/BudgetBot/internal/Message"
	"github.com/MaxKut3/BudgetBot/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client interface {
	Run()
	Stop()
}

type Currency interface {
	GetValue(cur string) int
}

/*
type Cache interface {
	Get(msg *Message.Message) (int, bool)
	Set(msg *Message.Message, c Currency)
}
*/

type client struct {
	bot      *tgbotapi.BotAPI
	cache    pkg.Cache
	currency Currency
}

func NewTgBot(cfg *config.TgBotConfig, cache pkg.Cache, currency Currency) *client {

	tgBot, err := tgbotapi.NewBotAPI(cfg.Key)
	if err != nil {
		log.Panic(fmt.Errorf("authorization failed: %v", err))
	}

	log.Println("Authorization was successful")

	return &client{
		bot:      tgBot,
		cache:    cache,
		currency: currency,
	}
}

func (c *client) Run() {
	c.bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := c.bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		wordList := strings.Split(update.Message.Text, " ")

		if !validateMessageText(wordList) {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невалидная строка. Строка должна быть следующего вида: Продукты 1000Rub")
			if _, sendError := c.bot.Send(msg); sendError != nil {
				log.Println(fmt.Errorf("send message failed: %v", sendError))
			}

			continue
		}

		msg := stringParser(wordList)

		if _, ok := c.cache.Get(msg); !ok {
			c.cache.Set(msg, c.currency)
		}

		val, _ := c.cache.Get(msg)

		msg = msg.SetSumRub(val)

		ans := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s, сумма вашей покупки составила: %d, в следующей валюте: %s. Сумма в рублях - %d . Категория покупки - %s ", update.Message.From.UserName, msg.Sum, msg.Cur, msg.SumRub/100, msg.Category))
		if _, sendError := c.bot.Send(ans); sendError != nil {
			log.Panic(fmt.Errorf("send message failed: %v", sendError))
		}

	}
}

func (c *client) Stop() {
	c.bot.StopReceivingUpdates()
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

var re, _ = regexp.Compile("[A-z]")

func stringParser(wordList []string) *Message.Message {

	listInd := re.FindStringIndex(wordList[1])
	i := listInd[0]

	sum, _ := strconv.Atoi(wordList[1][:i])

	return &Message.Message{
		Category: wordList[0],
		Sum:      sum,
		Cur:      wordList[1][i:],
	}
}
