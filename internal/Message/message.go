package Message

type Message struct {
	Category string
	Sum      int
	Cur      string
	SumRub   int
}

func (msg *Message) SetSumRub(rate int) *Message {
	return &Message{
		Category: msg.Category,
		Sum:      msg.Sum,
		Cur:      msg.Cur,
		SumRub:   msg.Sum * rate,
	}
}
