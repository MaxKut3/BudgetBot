package main

type UpdateModel struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	From      From   `json:"from"`
	Text      string `json:"text"`
	Chat      Chat
}

type From struct {
	UserName string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}

type Request struct {
	Result []UpdateModel `json:"result"`
}
