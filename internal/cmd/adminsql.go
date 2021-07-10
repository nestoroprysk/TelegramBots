package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nestoroprysk/TelegramBots/internal/env"
	"github.com/nestoroprysk/TelegramBots/internal/errorreporter"
	"github.com/nestoroprysk/TelegramBots/internal/sql"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
	"github.com/nestoroprysk/TelegramBots/internal/util"

	"github.com/xwb1989/sqlparser"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	env, ok := env.New(env.Config{
		ResponseWriter: w,
		Request:        r,
		DBOpener:       sqlclient.NewOpener(),
		Poster:         http.DefaultClient,
		Telegram: telegramclient.Config{
			Token: os.Getenv("ADMIN_BOT_TOKEN"),
			// TODO: Move to a small function
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
		AssertAdminID: true,
	})
	if !ok {
		return
	}
	defer env.Close()

	stmt, err := sqlparser.Parse(env.Update.Message.Text)
	if err == nil {
		switch stmt.(type) {
		case *sqlparser.Select, *sqlparser.Show, *sqlparser.OtherRead:
			result, err := env.SQLClient.Query(sqlclient.Query{Statement: env.Update.Message.Text})
			if err == nil {
				env.Responder.Succeed(sql.FormatTable(result))
				return
			} else {
				// TODO: Move to a small function
				env.Responder.Succeed(fmt.Sprintf("invalid input SQL statement (%s): %s", env.Update.Message.Text, err.Error()))
				return
			}
		default:
			result, err := env.SQLClient.Exec(sqlclient.Query{Statement: env.Update.Message.Text})
			if err == nil {
				// TODO: Move to a small function
				env.Responder.Succeed(fmt.Sprintf("Query OK, %d %s affected", result.RowsAffected, util.Pluralize("row", int(result.RowsAffected))))
				return
			} else {
				// TODO: Move to a small function
				env.Responder.Succeed(fmt.Sprintf("invalid input SQL statement (%s): %s", env.Update.Message.Text, err.Error()))
				return
			}
		}
	} else {
		// TODO: Move to a small function
		env.Responder.Succeed(fmt.Sprintf("invalid input SQL statement (%s): %s", env.Update.Message.Text, err.Error()))
		return
	}
}
