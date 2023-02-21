package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	ChatID   int64
	Msg      string
	Category string
	Sum      int
	Cur      string
	SumRub   int
}

func NewMessage(update tgbotapi.Update) *Message {
	return &Message{
		ChatID: update.Message.Chat.ID,
		Msg:    update.Message.Text,
	}
}

func (msg *Message) ValidateMessage() bool {
	wordList := strings.Split(msg.Msg, " ")

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

func (msg *Message) ParseNewMessage() *Message {
	re, _ := regexp.Compile("[A-z]")

	wordList := strings.Split(msg.Msg, " ")
	listInd := re.FindStringIndex(wordList[1])
	i := listInd[0]

	sum, _ := strconv.Atoi(wordList[1][:i])

	return &Message{
		ChatID:   msg.ChatID,
		Category: wordList[0],
		Sum:      sum,
		Cur:      wordList[1][i:],
	}
}

func (msg *Message) SetSumRub(rate int) *Message {
	return &Message{
		ChatID:   msg.ChatID,
		Category: msg.Category,
		Sum:      msg.Sum,
		Cur:      msg.Cur,
		SumRub:   msg.Sum * rate,
	}
}

func (msg *Message) ResponseFormation() string {
	resp := fmt.Sprintf("Сумма вашей покупки составила: %d, в следующей валюте: %s. Сумма в рублях - %d . Категория покупки - %s ", msg.Sum, msg.Cur, msg.SumRub/100, msg.Category)
	return resp
}
