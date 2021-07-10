package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/errorreporter"
	"github.com/nestoroprysk/TelegramBots/internal/sql"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/util"
	"github.com/xwb1989/sqlparser"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	env, ok := env.New(env.Config{
		ResponseWriter: w,
		Request:        r,
		DBOpener:       sqlclient.NewOpener(),
		Poster:         http.DefaultClient,
		Telegram: telegramclient.Config{
			Token:   os.Getenv("EXPENSES_BOT_TOKEN"),
			AdminID: func() int { result, _ := strconv.Atoi(os.Getenv("ADMIN_ID")); return result }(),
		},
		SQL: sqlclient.Config{
			Name:                   "information_schema",
			User:                   "root",
			Password:               os.Getenv("BOT_SQL_ROOT_PASS"),
			InstanceConnectionName: os.Getenv("BOT_SQL_CONNECTION_NAME"),
		},
		ErrorReporter: errorreporter.Config{
			ProjectID:   os.Getenv("PROJECT_ID"),
			ServiceName: os.Getenv("SERVICE_NAME"),
		},
	})
	if !ok {
		return
	}
	defer env.Close()

	// TODO: Move to a small function
	id := "expenses" + strconv.Itoa(env.Update.Message.From.ID)

	// TODO: Drop
	if env.Update.Message.Text == "/error" {
		b := strings.TrimSpace(strings.TrimPrefix(env.Update.Message.Text, "/error"))
		env.Responder.Error(fmt.Errorf("%s", b))
		return
	} else if env.Update.Message.Text == "/start" {
		// TODO: Move to a function
		// TODO: Refactore inline SQL in some way
		if _, err := env.SQLClient.Exec(
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
		); err == nil {
			env.Responder.Succeed("Welcome! Type 'show tables' to begin...")
			return
		} else {
			env.Responder.Error(fmt.Errorf("failed to start the user (%s): %w", id, err))
			return
		}

		// TODO: Drop the user on stopping the bot in a while
	} else {
		// TODO: Move to a function
		s, err := env.UserSQLClient(id)
		if err != nil {
			env.Responder.Error(err)
			return
		}

		stmt, err := sqlparser.Parse(env.Update.Message.Text)
		if err == nil {
			switch stmt.(type) {
			case *sqlparser.Select, *sqlparser.Show, *sqlparser.OtherRead:
				result, err := s.Query(sqlclient.Query{Statement: env.Update.Message.Text})
				if err == nil {
					env.Responder.Succeed(sql.FormatTable(result))
					return
				} else {
					env.Responder.Succeed(fmt.Sprintf("invalid input SQL statement (%s): %s", env.Update.Message.Text, err.Error()))
					return
				}
			default:
				result, err := s.Exec(sqlclient.Query{Statement: env.Update.Message.Text})
				if err == nil {
					// TODO: Move to a small function
					env.Responder.Succeed(fmt.Sprintf("Query OK, %d %s affected", result.RowsAffected, util.Pluralize("row", int(result.RowsAffected))))
					return
				} else {
					env.Responder.Succeed(fmt.Sprintf("invalid input SQL statement (%s): %s", env.Update.Message.Text, err.Error()))
					return
				}
			}
		} else {
			env.Responder.Succeed(fmt.Sprintf("invalid input SQL statement (%s): %s", env.Update.Message.Text, err.Error()))
			return
		}
	}
}
