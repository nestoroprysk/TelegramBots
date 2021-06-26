package telegram

import (
	"encoding/json"
	"fmt"
	"io"
)

// Update is a Telegram object that the handler receives every time a user interacts with the bot.
type Update struct {
	Message Message `json:"message" validate:"required"`
	From    From    `json:"from", validate:"required"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text" validate:"required"`
	Chat Chat   `json:"chat" validate:"required"`
}

// From is a Telegram object that indicates a user that sends a message.
type From struct {
	ID int `json:"id" validate:"gt=0"`
}

// Chat indicates the conversation to which the message belongs.
type Chat struct {
	ID int `json:"id" validate:"gt=0"`
}

// Parse parses the update and closes the body.
func Parse(body io.ReadCloser) (Update, error) {
	defer body.Close()

	var update Update
	if err := json.NewDecoder(body).Decode(&update); err != nil {
		err := fmt.Errorf("could not decode an incoming update: %w", err)
		return Update{}, err
	}

	return update, nil
}
