package cmd

import (
	"net/http"
	"os"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/logger"
	"github.com/nestoroprysk/TelegramBots/internal/parser"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
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
	if err := env.IsValid(); err != nil {
		// TODO: Capture
		l.Fatalf("failed to initialize the environment: %s", err.Error())
	}

	p := parser.New()
	u, err := p.Update(r.Body)
	if err != nil {
		l.Printf("failed to parse the update: %s", err.Error())
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
		l.Printf("successfully sent %q to %q: %s", text, u.Message.Chat.ID, response)
	} else {
		// TODO: capture
		l.Fatalf("failed to send %q to %q: %s", text, u.Message.Chat.ID, err.Error())
	}

	// TODO: respond in some way if applicable
}

// HangeAdminSQL sends an SQL query and response with the result to Telegram.
/*
func HandleAdminSQL(_ http.ResponseWriter, r *http.Request) {
	util.Bar()
	logger := log.New(os.Stderr, "HandleAdminSQL", log.LstdFlags)

	env, err := MakeEnv()
	if err != nil {
		logger.Fatalf("failed to initialize the environment: %s", err.Error())
	}

	update, err := ParseTelegramRequest(r)
	if err != nil {
		logger.Printf("failed to parse the update: %s", err.Error())
		return
	}

	telegram := NewTelegramClient(env.Telegram)
	adminSQL := NewSQLClient(env.DB)

	HandleAdminSQLUpdate(logger, update, telegram, adminSQL)
}

// TODO: unit test

func HandleAdminSQLUpdate(logger Logger, update Update, telegram TelegramClient, adminSQL SQLClient) {
	text, err := adminSQL.Send(update.Message.Text)
	if err != nil {
		logger.Printf("%s", err.Error())
		text = err.Error()
	}

	response, err := telegram.Send(text, update.Message.Chat.ID)
	if err == nil {
		logger.Printf("successfully sent %q to %q: %s", text, update.Message.Chat.ID, response)
	} else {
		// TODO: capture
		logger.Fatalf("failed to send %q to %q: %s", text, update.Message.Chat.ID, err.Error())
	}

	// TODO: respond in some way if applicable
}*/
