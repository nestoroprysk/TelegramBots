package bot

import (
	"log"
	"net/http"
)

// HangeAdminSQL sends an SQL query and response with the result to Telegram.
func HandleAdminSQL(w http.ResponseWriter, r *http.Request) {
	env, err := MakeEnv()
	if err != nil {
		log.Fatalf("failed to initialize the environment: %s", err.Error())
		return
	}

	telegram := NewTelegramClient(env.Telegram)
	admin := NewSQLClient(env.DB)

	update, err := ParseTelegramRequest(r)
	if err != nil {
		log.Printf("failed to parse the update: %s", err.Error())
		return
	}

	text, err := admin.Send(update.Message.Text)
	if err != nil {
		log.Printf("%s", err.Error())
		text = err.Error()
	}

	response, err := telegram.Send(text, update.Message.Chat.ID)
	if err == nil {
		log.Printf("successfully sent %q to %q: %s", text, update.Message.Chat.ID, response)
	} else {
		// TODO: capture
		log.Fatalf("failed to send %q to %q: %s", text, update.Message.Chat.ID, err.Error())
	}
}
