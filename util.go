package bot

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// ParseTelegramRequest handles incoming update from the Telegram web hook.
func ParseTelegramRequest(r *http.Request) (Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return Update{}, err
	}

	return update, nil
}
