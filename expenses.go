package bot

import (
	"net/http"
)

func HandleExpenses(w http.ResponseWriter, r *http.Request) {
	// TODO: implement DB per customer
	HandleAdminSQL(w, r)
}
