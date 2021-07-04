package telegramclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
)

// TODO: cover with unit tests

// TelegramClient is an interface for sending text to chat.
type TelegramClient interface {
	Send(text string) (telegram.Response, error)
}

// Poster posts an HTTP request.
type Poster interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

var _ TelegramClient = &telegramClient{}

type telegramClient struct {
	token  string
	chatID string
	client Poster
}

// NewTelegramClient creates a Telegram client.
func New(conf env.Telegram, chatID int, client Poster) TelegramClient {
	return &telegramClient{
		token:  conf.Token,
		chatID: strconv.Itoa(chatID),
		client: client,
	}
}

// Send sends text to chat.
func (tc telegramClient) Send(text string) (telegram.Response, error) {
	response, err := tc.client.PostForm(
		"https://api.telegram.org/bot"+tc.token+"/sendMessage",
		url.Values{
			"chat_id":    {tc.chatID},
			"text":       {"```" + text + "```"},
			"parse_mode": {"markdown"},
		},
	)
	if err != nil {
		err := fmt.Errorf("error when posting text to the chat %q: %w", tc.chatID, err)
		return telegram.Response{}, err
	}

	if response.StatusCode != http.StatusOK {
		return telegram.Response{}, fmt.Errorf("expecting status code %d for the Telegram response; got %d", http.StatusOK, response.StatusCode)
	}

	result, err := telegram.ParseResponse(response.Body)
	if err != nil {
		return telegram.Response{}, err
	}

	if result.Ok == false {
		return telegram.Response{}, fmt.Errorf("expecting ok; got %+v", result)
	}

	if result.ErrorCode != 0 {
		return telegram.Response{}, fmt.Errorf("expecting zero exit code; got %+v", result)
	}

	return result, nil
}
