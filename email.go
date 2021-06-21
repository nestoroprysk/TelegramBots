package api

import (
	"net/http"

	"github.com/nestoroprysk/TelegramBots/internal/cmd"
)

func Email(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	cmd.Email(w, r)
}
