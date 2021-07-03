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
	"github.com/xwb1989/sqlparser"
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

	id := "expenses" + strconv.Itoa(u.Message.From.ID)

	var text string

	if u.Message.Text == "/start" {
		s, err := sqlclient.New(e.DB, sqlclient.NewOpener())
		if err != nil {
			// TODO: Capture
			resp.Error(err)
			return
		}

		// TODO: Refactore inline SQL in some way
		_, err = s.Exec(
			sqlclient.Query{
				Statement: "CREATE DATABASE IF NOT EXISTS ?;",
				Args:      []interface{}{id},
			},
			sqlclient.Query{
				Statement: "CREATE TABLE IF NOT EXISTS ?.e (d date DEFAULT NULL, v int(10) unsigned DEFAULT NULL, c varchar(20) DEFAULT NULL);",
				Args:      []interface{}{id},
			},
			sqlclient.Query{
				Statement: "CREATE USER IF NOT EXISTS '?'@'%%';",
				Args:      []interface{}{id},
			},
			sqlclient.Query{
				Statement: "GRANT ALL PRIVILEGES ON ?.e TO '?'@'%%';",
				Args:      []interface{}{id, id},
			},
		)
		if err != nil {
			// TODO: Capture
			resp.Error(fmt.Errorf("failed to start the user (%s): %w", id, err))
			return
		}

		// TODO: Drop the user on stopping the bot
		text = "Welcome! Type 'show tables' to begin..."
	} else {
		user := env.DB{
			Name:                   id,
			User:                   id,
			InstanceConnectionName: os.Getenv("BOT_SQL_CONNECTION_NAME"),
		}

		s, err := sqlclient.New(user, sqlclient.NewOpener())
		if err != nil {
			// TODO: Capture
			resp.Error(err)
			return
		}

		stmt, err := sqlparser.Parse(u.Message.Text)
		if err != nil {
			resp.Fail(fmt.Errorf("invalid input SQL statement (%s): %w", u.Message.Text, err))
			return
		}

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
	}

	t := telegramclient.New(e.Telegram, u.Message.Chat.ID, http.DefaultClient)
	response, err := t.Send(text)
	if err != nil {
		// TODO: capture
		resp.Error(err)
		return
	}

	resp.Succeed(response)
}
