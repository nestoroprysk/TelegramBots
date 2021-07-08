package sql

import (
	"fmt"

	"github.com/nestoroprysk/TelegramBots/internal/util"
)

// Result is the result of a non-select SQL statement.
type Result struct {
	// RowsAffected counts affected rows.
	RowsAffected int64
	// LastInsertID indicates the last insertID.
	LastInsertID int64
}

// FormatResult formats an SQL exec response.
func FormatResult(r Result) string {
	return fmt.Sprintf("Query OK, %d %s affected", r.RowsAffected, util.Pluralize("row", int(r.RowsAffected)))
}
