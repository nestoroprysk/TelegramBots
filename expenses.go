package api

import (
	"net/http"

	"github.com/nestoroprysk/TelegramBots/internal/cmd"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	cmd.Expenses(w, r)
}
