package bot

import (
	"net/http"
)

func HandleExpenses(w http.ResponseWriter, r *http.Request) {
	HandleAdminSQL(w, r)
}
