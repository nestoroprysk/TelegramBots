package env

import (
	"fmt"
	"net/http"

	"github.com/nestoroprysk/TelegramBots/internal/errorreporter"
	"github.com/nestoroprysk/TelegramBots/internal/httpresponder"
	"github.com/nestoroprysk/TelegramBots/internal/logger"
	"github.com/nestoroprysk/TelegramBots/internal/responder"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/validator"
)

// Config defines Env.
type Config struct {
	http.ResponseWriter
	*http.Request
	sqlclient.DBOpener
	telegramclient.Poster
	Telegram      telegramclient.Config
	SQL           sqlclient.Config
	ErrorReporter errorreporter.Config
	AssertAdminID bool
}

type Env struct {
	responder.Responder
	sqlclient.SQLClient
	UserSQLClient
	telegram.Update
}

type UserSQLClient func(id string) (sqlclient.SQLClient, error)

// New creates a new environment.
func New(conf Config) (Env, bool) {
	resp := httpresponder.Wrap(logger.New(), conf.ResponseWriter)
	resp, ok := errorreporter.Wrap(resp, conf.ErrorReporter)
	if !ok {
		return Env{}, false
	}

	v := validator.New()
	if err := v.Struct(conf); err != nil {
		resp.Error(fmt.Errorf("failed to initialize the environment: %w", err))
		resp.Close()
		return Env{}, false
	}

	u, err := telegram.ParseUpdate(conf.Request.Body)
	if err != nil {
		resp.Fail(fmt.Errorf("failed to parse the update: %w", err))
		resp.Close()
		return Env{}, false
	}

	if err := v.Struct(u); err != nil {
		resp.Fail(fmt.Errorf("failed to validate the update: %w", err))
		resp.Close()
		return Env{}, false
	}

	s, err := sqlclient.New(conf.SQL, conf.DBOpener)
	if err != nil {
		resp.Error(err)
		resp.Close()
		return Env{}, false
	}

	userSQLClient := func(id string) (sqlclient.SQLClient, error) {
		user := sqlclient.Config{
			Name:                   id,
			User:                   id,
			Password:               "",
			InstanceConnectionName: conf.SQL.InstanceConnectionName,
		}

		return sqlclient.New(user, conf.DBOpener)
	}

	resp = telegramclient.Wrap(resp, conf.Telegram, u.Message.Chat.ID, conf.Poster)

	if conf.AssertAdminID {
		if u.Message.From.ID != conf.Telegram.AdminID {
			resp.Fail(fmt.Errorf("user id (%d) is not authenticated to call the function", u.Message.From.ID))
			return Env{}, false
		}
	}

	return Env{
		Responder:     resp,
		SQLClient:     s,
		UserSQLClient: userSQLClient,
		Update:        u,
	}, true
}
