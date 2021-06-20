package cmd

import (
	"net/http"
	"os"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/logger"
	"github.com/nestoroprysk/TelegramBots/internal/parser"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/validator"
)

func Admin(_ http.ResponseWriter, r *http.Request) {
	env := env.Env{
		Telegram: env.Telegram{
			Token: os.Getenv("ADMIN_BOT_TOKEN"),
		},
		DB: env.DB{
			Name:                   "information_schema",
			User:                   "root",
			Password:               os.Getenv("BOT_SQL_ROOT_PASS"),
			InstanceConnectionName: os.Getenv("BOT_SQL_CONNECTION_NAME"),
		},
	}

	l := logger.New()
	v := validator.New()

	if err := v.Struct(env); err != nil {
		l.Fatalf("failed to initialize the environment: %s", err.Error())
	}

	p := parser.New()
	u, err := p.Update(r.Body)
	if err != nil {
		l.Printf("failed to parse the update: %s", err.Error())
		return
	}

	if err := v.Struct(u); err != nil {
		l.Printf("failed to validate the update: %s", err.Error())
		return
	}

	s := sqlclient.New(env.DB)
	text, err := s.Send(u.Message.Text)
	if err != nil {
		l.Printf("%s", err.Error())
		text = err.Error()
	}

	t := telegramclient.New(env.Telegram, u.Message.Chat.ID)
	response, err := t.Send(text)
	if err == nil {
		l.Printf("successfully sent %q to %d: %s", text, u.Message.Chat.ID, response)
	} else {
		// TODO: capture
		l.Fatalf("failed to send %q to %d: %s", text, u.Message.Chat.ID, err.Error())
	}

	// TODO: respond in some way
}
