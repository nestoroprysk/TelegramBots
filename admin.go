package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nestoroprysk/TelegramBots/internal/cmd"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s", os.Getenv("BOT_SQL_CONNECTION_NAME"))
	cmd.Admin(w, r)
}
