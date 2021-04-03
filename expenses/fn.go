package expenses

import (
	"net/http"

	"github.com/nestoroprysk/TelegramBots/admin"
)

func HandleExpenses(w http.ResponseWriter, r *http.Request) {
	admin.HandleAdminSQL(w, r)
}
