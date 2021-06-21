package util

import (
	"strings"

	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
)

func Format(r sqlclient.Response) string {
	// TODO: use some pretty library

	var result []string
	result = append(result, strings.Join(r.Columns, " "))

	for _, row := range r.Rows {
		result = append(result, strings.Join(row, " "))
	}

	return strings.Join(result, "\n")
}
