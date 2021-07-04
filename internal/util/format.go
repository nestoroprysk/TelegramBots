package util

import (
	"database/sql"
	"fmt"

	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Format formats an SQL query response in a CSV format.
func Format(r sqlclient.Table) string {
	t := table.NewWriter()

	var cols []interface{}
	for _, c := range r.Columns {
		cols = append(cols, c)
	}
	t.AppendHeader(cols)

	for _, row := range r.Rows {
		var items []interface{}
		for _, c := range r.Columns {
			if i, ok := row[c]; ok {
				if b, ok := i.(sql.RawBytes); ok {
					i = string(b)
				}

				items = append(items, i)
			}
		}

		t.AppendRow(items)
	}

	return t.Render()
}

// FormatResult formats an SQL exec response.
func FormatResult(r sqlclient.Result) string {
	return fmt.Sprintf("Query OK, %d %s affected", r.RowsAffected, Pluralize("row", int(r.RowsAffected)))
}

// Format code formats str as code in MD.
func FormatCode(s string) string {
	return "```\n" + s + "\n```"
}
