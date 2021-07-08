package sql

import (
	"database/sql"

	"github.com/jedib0t/go-pretty/table"
)

// Table is the result of the select statement.
type Table struct {
	// Columns is a list of columns.
	Columns []string
	// Rows is a list of maps from a column name to the result.
	Rows []map[string]interface{}
}

// FormatTable formats an SQL query response in a CSV format.
func FormatTable(r Table) string {
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
