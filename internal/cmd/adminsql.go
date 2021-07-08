package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/errorreporter"
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
			Token:   os.Getenv("ADMIN_BOT_TOKEN"),
			AdminID: func() int { result, _ := strconv.Atoi(os.Getenv("ADMIN_ID")); return result }(),
		},
		DB: env.DB{
			Name:                   "information_schema",
			User:                   "root",
			Password:               os.Getenv("BOT_SQL_ROOT_PASS"),
			InstanceConnectionName: os.Getenv("BOT_SQL_CONNECTION_NAME"),
		},
		ErrorReporter: errorreporter.Config{
			ProjectID:   os.Getenv("PROJECT_ID"),
			ServiceName: os.Getenv("SERVICE_NAME"),
		},
	}

	v := validator.New()
	resp := responder.New(w)

	// Responder | Scenario | Action
	// HTTP      | Success  | Succeed()
	// Logger    | Success  | Succeed()
	// Telegram  | Success  | Succeed()
	// - Error reporter

	errorReporter, err := errorreporter.New(env.ErrorReporter)
	if err != nil {
		log.Print(fmt.Errorf("failed to initialize the error reporter: %w", err).Error())
	}
	defer errorReporter.Close()

	if err := v.Struct(env); err != nil {
		errorReporter.Error(err)
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

	if u.Message.From.ID != env.AdminID {
		resp.Fail(fmt.Errorf("user id (%d) is not authenticated to call the function", u.Message.From.ID))
		return
	}

	s, err := sqlclient.New(env.DB, sqlclient.NewOpener())
	if err != nil {
		errorReporter.Error(err)
		resp.Error(err)
		return
	}

	var text string

	stmt, err := sqlparser.Parse(u.Message.Text)
	if err == nil {
		switch stmt.(type) {
		case *sqlparser.Select, *sqlparser.Show, *sqlparser.OtherRead:
			result, err := s.Query(sqlclient.Query{Statement: u.Message.Text})
			if err == nil {
				text = util.Format(result)
			} else {
				err := fmt.Errorf("invalid input SQL statement (%s): %w", u.Message.Text, err)
				text = err.Error() // Hint user that the SQL statement is not ok.
			}
		default:
			result, err := s.Exec(sqlclient.Query{Statement: u.Message.Text})
			if err == nil {
				text = fmt.Sprintf("Query OK, %d %s affected", result.RowsAffected, util.Pluralize("row", int(result.RowsAffected)))
			} else {
				err := fmt.Errorf("invalid input SQL statement (%s): %w", u.Message.Text, err)
				text = err.Error() // Hint user that the SQL statement is not ok.
			}
		}
	} else {
		err := fmt.Errorf("invalid input SQL statement (%s): %w", u.Message.Text, err)
		text = err.Error() // Hint user that the SQL statement is not ok.
	}

	t := telegramclient.New(env.Telegram, u.Message.Chat.ID, http.DefaultClient)
	response, err := t.Send(text)
	if err != nil {
		errorReporter.Error(err)
		resp.Error(err)
		return
	}

	resp.Succeed(response)
}
