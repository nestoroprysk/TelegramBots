package bot

import (
	"log"
	"net/http"
	"os"

	"github.com/nestoroprysk/TelegramBots/internal/util"
)

// HangeAdminSQL sends an SQL query and response with the result to Telegram.
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
}
