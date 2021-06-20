package telegram

// Update is a Telegram object that the handler receives every time a user interacts with the bot.
type Update struct {
	Message Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Chat indicates the conversation to which the message belongs.
type Chat struct {
	ID int `json:"id"`
}
