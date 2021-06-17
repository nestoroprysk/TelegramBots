package bot

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// TelegramClientConfig configures a Telegram client.
type TelegramClientConfig struct {
	// Token is a telegram bot token.
	Token string
}

// TelegramClient is an interface for sending text to chat.
type TelegramClient interface {
	Send(text string, chatID int) (response string, err error)
}

type telegramClient struct {
	token string
}

type mockTelegramClient struct {
	f func(text string, chatID int) (response string, err error)
}

// NewTelegramClient creates a Telegram client.
func NewTelegramClient(conf Telegram) TelegramClient {
	return &telegramClient{
		token: conf.Token,
	}
}

// NewTelegramClient creates a mock Telegram client.
func NewMockTelegramClient(f func(text string, chatID int) (response string, err error)) TelegramClient {
	return &mockTelegramClient{
		f: f,
	}
}

// Send sends text to chat.
func (tc telegramClient) Send(text string, chatID int) (string, error) {
	log.Printf("sending %s to chat_id: %d", text, chatID)

	response, err := http.PostForm(
		"https://api.telegram.org/bot"+tc.token+"/sendMessage",
		url.Values{
			"chat_id": {strconv.Itoa(chatID)},
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

// Send mocks sending text to chat.
func (tc mockTelegramClient) Send(text string, chatID int) (string, error) {
	return tc.f(text, chatID)
}
