package util

import (
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Format formats an SQL response in a CSV format.
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
				if b, ok := i.([]byte); ok {
					i = string(b)
				}

				items = append(items, i)
			}
		}

		t.AppendRow(items)
	}

	// TODO: Consider MD table
	return t.RenderCSV()
}
