package parser

import (
	"encoding/json"
	"io"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
)

type Parser interface {
	Update(body io.ReadCloser) (telegram.Update, error)
}

type parser struct{}

// TODO: Add interface implementation declaration

// TODO: Provide a mock

// TODO: Add logger as input

func New() Parser {
	return &parser{}
}

// Update parses the update and closes the body.
func (parser) Update(body io.ReadCloser) (telegram.Update, error) {
	defer body.Close()

	var update telegram.Update
	if err := json.NewDecoder(body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return telegram.Update{}, err
	}

	return update, nil
}
