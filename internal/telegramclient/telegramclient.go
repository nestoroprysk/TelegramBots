package telegramclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nestoroprysk/TelegramBots/internal/responder"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
	"github.com/nestoroprysk/TelegramBots/internal/util"
)

// TODO: cover with unit tests

// Poster posts an HTTP request.
type Poster interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

var _ responder.Responder = &telegramClient{}

type telegramClient struct {
	responder.Responder
	token  string
	chatID string
	client Poster
}

// Config defines the Telegram client.
type Config struct {
	// Token is a Telegram bot token.
	Token string `validate:"required"`
	// AdminID is an ID of the admin Telegram user.
	AdminID int `validate:"gt=0"`
}

// Wrap adds a Telegram client the responder.
func Wrap(r responder.Responder, conf Config, chatID int, client Poster) responder.Responder {
	return &telegramClient{
		Responder: r,
		token:     conf.Token,
		chatID:    strconv.Itoa(chatID),
		client:    client,
	}
}

func (tc telegramClient) Succeed(b interface{}) error {
	resp, err := tc.send(b)
	if err != nil {
		return tc.Responder.Error(err)
	}

	return tc.Responder.Succeed(resp)
}

func (tc telegramClient) Fail(err error) error {
	if _, errt := tc.send("Something is wrong with Telegram! Try again later..."); errt != nil {
		return tc.Responder.Error(util.CombineErrors(err, errt))
	}

	return tc.Responder.Fail(err)
}

func (tc telegramClient) Error(err error) error {
	if _, errt := tc.send("Something is wrong with me! Try again later..."); errt != nil {
		return tc.Responder.Error(util.CombineErrors(err, errt))
	}

	return tc.Responder.Error(err)
}

func (tc telegramClient) Close() error {
	return tc.Responder.Close()
}

// Send sends text to chat.
func (tc telegramClient) send(text interface{}) (telegram.Response, error) {
	response, err := tc.client.PostForm(
		"https://api.telegram.org/bot"+tc.token+"/sendMessage",
		url.Values{
			"chat_id":    {tc.chatID},
			"text":       {util.FormatCode(text)},
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
