package util

import (
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Format(r sqlclient.Response) string {
	t := table.NewWriter()

	var cols []interface{}
	for _, c := range r.Columns {
		cols = append(cols, c)
	}
	t.AppendHeader(cols)

	for _, row := range r.Rows {
		t.AppendRow(row)
	}

	return t.Render()
}
