package telegramclient

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nestoroprysk/TelegramBots/internal/env"
)

// TelegramClient is an interface for sending text to chat.
type TelegramClient interface {
	Send(text string) (response string, err error)
}

type telegramClient struct {
	token  string
	chatID string
}

type mockTelegramClient struct {
	f func(text string, chatID int) (response string, err error)
}

// NewTelegramClient creates a Telegram client.
func New(conf env.Telegram, chatID int) TelegramClient {
	return &telegramClient{
		token:  conf.Token,
		chatID: strconv.Itoa(chatID),
	}
}

// Send sends text to chat.
func (tc telegramClient) Send(text string) (string, error) {
	log.Printf("sending %q to chat_id: %d", text, tc.chatID)

	response, err := http.PostForm(
		"https://api.telegram.org/bot"+tc.token+"/sendMessage",
		url.Values{
			"chat_id": {tc.chatID},
			"text":    {text},
		},
	)
	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error in parsing the Telegram response: %s", err.Error())
		return "", err
	}

	log.Printf("body of the Telegram response: %s", body)

	return string(body), nil
}
