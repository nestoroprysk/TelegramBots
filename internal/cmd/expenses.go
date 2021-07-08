package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/errorreporter"
	"github.com/nestoroprysk/TelegramBots/internal/responder"
	"github.com/nestoroprysk/TelegramBots/internal/sql"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/util"
	"github.com/nestoroprysk/TelegramBots/internal/validator"
	"github.com/xwb1989/sqlparser"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	e := env.Env{
		Telegram: telegramclient.Config{
			Token:   os.Getenv("EXPENSES_BOT_TOKEN"),
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

	errorReporter, err := errorreporter.New(e.ErrorReporter)
	if err != nil {
		log.Print(fmt.Errorf("failed to initialize the error reporter: %w", err).Error())
	}
	defer errorReporter.Close()

	if err := v.Struct(e); err != nil {
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

	id := "expenses" + strconv.Itoa(u.Message.From.ID)

	var text string

	if strings.HasPrefix(u.Message.Text, "/error") { // TODO: Drop.
		text = fmt.Sprintf("Reporting (%s)...", strings.TrimSpace(strings.TrimPrefix(u.Message.Text, "/error")))
		errorReporter.Error(fmt.Errorf("%s", text))
	} else if u.Message.Text == "/start" {
		s, err := sqlclient.New(e.DB, sqlclient.NewOpener())
		if err != nil {
			errorReporter.Error(err)
			resp.Error(err)
			return
		}

		// TODO: Refactore inline SQL in some way

		_, err = s.Exec(
			sqlclient.Query{
				Statement: fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", id),
			},
			sqlclient.Query{
				Statement: fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.e (d date DEFAULT NULL, v int(10) unsigned DEFAULT NULL, c varchar(20) DEFAULT NULL);", id),
			},
			sqlclient.Query{
				Statement: fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%';", id),
			},
			sqlclient.Query{
				Statement: fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.e TO '%s'@'%%';", id, id),
			},
		)
		if err != nil {
			errorReporter.Error(err)
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
			errorReporter.Error(err)
			resp.Error(err)
			return
		}

		stmt, err := sqlparser.Parse(u.Message.Text)
		if err == nil {
			switch stmt.(type) {
			case *sqlparser.Select, *sqlparser.Show, *sqlparser.OtherRead:
				result, err := s.Query(sqlclient.Query{Statement: u.Message.Text})
				if err == nil {
					text = sql.FormatTable(result)
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
	}

	t := telegramclient.New(e.Telegram, u.Message.Chat.ID, http.DefaultClient)
	response, err := t.Send(text)
	if err != nil {
		errorReporter.Error(err)
		resp.Error(err)
		return
	}

	resp.Succeed(response)
}
