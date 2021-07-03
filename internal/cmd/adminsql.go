package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/responder"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/util"
	"github.com/nestoroprysk/TelegramBots/internal/validator"

	"github.com/xwb1989/sqlparser"
)

func Admin(w http.ResponseWriter, r *http.Request) {
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
		resp.Error(fmt.Errorf("failed to initialize the environment: %w", err))
		return
	}

	u, err := telegram.ParseUpdate(r.Body)
	if err != nil {
		resp.Fail(fmt.Errorf("failed to parse the update: %w", err))
		return
	}

	if err := v.Struct(u); err != nil {
		resp.Fail(fmt.Errorf("failed to validate the update: %w", err))
		return
	}

	const adminID = 381126698 // TODO: Move it to a secret
	if u.Message.From.ID != adminID {
		resp.Fail(fmt.Errorf("user id (%d) is not authenticated to call the function", u.Message.From.ID))
		return
	}

	s, err := sqlclient.New(env.DB, sqlclient.NewOpener())
	if err != nil {
		resp.Error(err)
		return
	}

	stmt, err := sqlparser.Parse(u.Message.Text)
	if err != nil {
		resp.Fail(fmt.Errorf("invalid input SQL statement (%s): %w", u.Message.Text, err))
		return
	}

	var text string

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		_ = stmt // Silence issues.

		result, err := s.Query(sqlclient.Query{Statement: u.Message.Text})
		if err == nil {
			text = util.Format(result)
		} else {
			text = err.Error() // Even if invalid SQL, send it.
		}
	default:
		result, err := s.Exec(sqlclient.Query{Statement: u.Message.Text})
		if err == nil {

			text = fmt.Sprintf("Query OK, affected %d rows", result.RowsAffected)
		} else {
			text = err.Error() // Even if execution failed, send the resulting error.
		}
	}

	t := telegramclient.New(env.Telegram, u.Message.Chat.ID, http.DefaultClient)
	response, err := t.Send(text)
	if err != nil {
		// TODO: capture
		resp.Error(err)
		return
	}

	resp.Succeed(response)
}
