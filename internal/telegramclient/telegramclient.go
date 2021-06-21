package telegramclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nestoroprysk/TelegramBots/internal/env"
)

// TelegramClient is an interface for sending text to chat.
type TelegramClient interface {
	Send(text string) (response []byte, err error)
}

type telegramClient struct {
	token  string
	chatID string
}

// NewTelegramClient creates a Telegram client.
func New(conf env.Telegram, chatID int) TelegramClient {
	return &telegramClient{
		token:  conf.Token,
		chatID: strconv.Itoa(chatID),
	}
}

// Send sends text to chat.
func (tc telegramClient) Send(text string) ([]byte, error) {
	// TODO: inject HTTP client for testing
	//       and repond just like the real telegram
	//       for both success and error cases

	response, err := http.PostForm(
		"https://api.telegram.org/bot"+tc.token+"/sendMessage",
		url.Values{
			"chat_id": {tc.chatID},
			"text":    {text},
		},
	)
	if err != nil {
		err := fmt.Errorf("error when posting text to the chat %q: %w", tc.chatID, err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err := fmt.Errorf("error in parsing the Telegram response: %w", err)
		return nil, err
	}

	// TODO: consider parsing the response and checking if ok
	// TODO: consider checking the exit code

	return body, nil
}
