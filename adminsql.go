package bot

import (
	"log"
	"net/http"
)

func HandleAdminSQL(_ http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	var update, err = ParseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	text, err := SendSQLRequest(update.Message.Text)
	if err != nil {
		log.Printf("error interacting with SQL, %s", err.Error())
		text = err.Error()
	}

	// Send the SQL response back to Telegram
	response, err := SendTextToTelegramChat(update.Message.Chat.Id, text)
	if err != nil {
		log.Printf("got error %s from telegram, reponse body is %s", err.Error(), response)
	} else {
		log.Printf("response %s successfuly distributed to chat id %d", update.Message.Text, update.Message.Chat.Id)
	}
}
