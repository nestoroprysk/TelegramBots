package telegram

// Update is a Telegram object that the handler receives every time a user interacts with the bot.
type Update struct {
	Message Message `json:"message" validate:"required"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text" validate:"required"`
	Chat Chat   `json:"chat" validate:"required"`
}

// Chat indicates the conversation to which the message belongs.
type Chat struct {
	ID int `json:"id" validate:"gt=0"`
}
