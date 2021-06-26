package api

import (
	"net/http"

	"github.com/nestoroprysk/TelegramBots/internal/cmd"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	cmd.Admin(w, r)
}
