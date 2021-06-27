package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/responder"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/util"
	"github.com/nestoroprysk/TelegramBots/internal/validator"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	log.Printf("(%s) (%s)", string(b), err)
	return
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

	v := validator.New()
	resp := responder.New(w)

	if err := v.Struct(env); err != nil {
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

	s, err := sqlclient.New(env.DB)
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

	t := telegramclient.New(env.Telegram, u.Message.Chat.ID)
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
