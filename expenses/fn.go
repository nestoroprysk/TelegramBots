package admin

import (
	"net/http"

	"github.com/nestoroprysk/TelegramBots/admin"
)

func HandleAdminSQL(w http.ResponseWriter, r *http.Request) {
	admin.HandleAdminSQL(w, r)
}
