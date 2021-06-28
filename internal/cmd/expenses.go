package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/responder"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/util"
	"github.com/nestoroprysk/TelegramBots/internal/validator"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	e := env.Env{
		Telegram: env.Telegram{
			Token: os.Getenv("EXPENSES_BOT_TOKEN"),
		},
		DB: env.DB{
			Name:                   "information_schema",
			User:                   "root",
			Password:               os.Getenv("BOT_SQL_ROOT_PASS"),
			InstanceConnectionName: os.Getenv("BOT_SQL_CONNECTION_NAME"),
		},
	}

	v := validator.New()
	resp := responder.New(w)

	if err := v.Struct(e); err != nil {
		// TODO: Capture
		resp.Respond(responder.Response{
			Status:  responder.Error,
			Data:    nil,
			Message: fmt.Errorf("failed to initialize the environment: %w", err).Error(),
		})
		return
	}

	u, err := telegram.Parse(r.Body)
	if err != nil {
		resp.Respond(responder.Response{
			Status:  responder.Fail,
			Data:    []byte(fmt.Errorf("failed to parse the update: %w", err).Error()),
			Message: "",
		})
		return
	}

	if err := v.Struct(u); err != nil {
		resp.Respond(responder.Response{
			Status:  responder.Fail,
			Data:    []byte(fmt.Errorf("failed to validate the update: %w", err).Error()),
			Message: "",
		})
		return
	}

	id := "expenses" + strconv.Itoa(u.Message.From.ID)

	if u.Message.Text == "/start" {
		_, err = sqlclient.New(e.DB)
		if err != nil {
			// TODO: Capture
			resp.Respond(responder.Response{
				Status:  responder.Error,
				Data:    nil,
				Message: err.Error(),
			})
			return
		}
		// TODO: execute the sql file from admin in the context to register a user if not registered

		s, err := sqlclient.New(e.DB)
		if err != nil {
			// TODO: Capture
			resp.Respond(responder.Response{
				Status:  responder.Error,
				Data:    nil,
				Message: err.Error(),
			})
			return
		}

		// TODO: Create a request type with input variables
		if err := s.Exec(
			fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", id),
			fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.e (d date DEFAULT NULL, v int(10) unsigned DEFAULT NULL, c varchar(20) DEFAULT NULL);", id),
			fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'%%';", id),
			fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.e TO '%s'@'%%';", id, id),
		); err != nil {
			// TODO: Capture
			resp.Respond(responder.Response{
				Status:  responder.Error,
				Data:    nil,
				Message: fmt.Errorf("failed to start the user (%s): %w", id, err).Error(),
			})
			return
		}

		// TODO: Drop the user on stopping the bot
	}

	user := env.DB{
		Name:                   id,
		User:                   id,
		InstanceConnectionName: os.Getenv("BOT_SQL_CONNECTION_NAME"),
	}

	s, err := sqlclient.New(user)
	if err != nil {
		// TODO: Capture
		resp.Respond(responder.Response{
			Status:  responder.Error,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	// TODO: parse SQL and error right away

	var text string
	result, err := s.Send(u.Message.Text)
	if err == nil {
		text = util.Format(result)
	} else {
		text = err.Error() // Even if invalid SQL, send it.
	}

	t := telegramclient.New(e.Telegram, u.Message.Chat.ID)
	response, err := t.Send(text)
	if err != nil {
		// TODO: capture
		resp.Respond(responder.Response{
			Status:  responder.Error,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	resp.Respond(responder.Response{
		Status:  responder.Success,
		Data:    response,
		Message: "",
	})
}
