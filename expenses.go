package api

import (
	"net/http"

	"github.com/nestoroprysk/TelegramBots/internal/cmd"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	cmd.Expenses(w, r)
}
