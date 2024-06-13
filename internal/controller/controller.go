package controller

import (
	"fmt"
	"log"

	"github.com/MaxKut3/BudgetBot/config"
	"github.com/MaxKut3/BudgetBot/internal/models"
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

type Repository interface {
	Insert(msg *models.Message)
	GetTotalAmount(update tgbotapi.Update) int
}

/*
	type Cache interface {
		Get(msg *models.Message) (int, bool)
		Set(msg *models.Message, c Currency)
	}
*/
type BotClient struct {
	Bot      *tgbotapi.BotAPI
	cache    pkg.Cache
	currency Currency
	conn     Repository
}

func NewTgBot(cfg *config.TgBotConfig, cache pkg.Cache, currency Currency, rep Repository) *BotClient {

	tgBot, err := tgbotapi.NewBotAPI(cfg.Key)
	if err != nil {
		log.Panic(fmt.Errorf("authorization failed: %v", err))
	}

	log.Println("Authorization was successful")

	return &BotClient{
		Bot:      tgBot,
		cache:    cache,
		currency: currency,
		conn:     rep,
	}
}

func (c *BotClient) Send(update tgbotapi.Update, resp string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
	if _, sendError := c.Bot.Send(msg); sendError != nil {
		log.Println(fmt.Errorf("send message failed: %v", sendError))
	}
}

func (c *BotClient) Run() {
	c.Bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := c.Bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message.Text == "Итог" {
			fmt.Println(c.conn.GetTotalAmount(update))
			c.Send(update, fmt.Sprintf("Сумма ваших трат за месяц: %d", c.conn.GetTotalAmount(update)))
			continue
		}

		models := models.NewMessage(update)

		if !models.ValidateMessage() {
			c.Send(update, "Невалидная строка. Строка должна быть следующего вида: Продукты 1000Rub")
			continue
		}

		models = models.ParseNewMessage()

		if _, ok := c.cache.Get(models); !ok {
			c.cache.Set(models, c.currency)
		}

		val, _ := c.cache.Get(models)

		models = models.SetSumRub(val)

		fmt.Println(models)

		c.conn.Insert(models)

		c.Send(update, models.ResponseFormation())

	}
}

func (c *BotClient) Stop() {
	c.Bot.StopReceivingUpdates()
}
