package sqlclient

import (
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nestoroprysk/TelegramBots/internal/sql"
)

// ConvertRows converts rows into a table.
//
// Source: https://kylewbanks.com/blog/query-result-to-map-in-golang.
func ConvertRows(rows Rows) (sql.Table, error) {
	cols, err := rows.Columns()
	if err != nil {
		return sql.Table{}, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return sql.Table{}, err
	}

	result := sql.Table{Columns: cols}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		for i := range columns {
			// Populating with pointers to concrete types.
			columns[i] = reflect.New(columnTypes[i].ScanType()).Interface()
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columns...); err != nil {
			return sql.Table{}, err
		}

		m := make(map[string]interface{})
		for i, c := range cols {
			// Getting values from pointers to types.
			m[c] = reflect.ValueOf(columns[i]).Elem().Interface()
		}

		result.Rows = append(result.Rows, m)
	}

	return result, nil
}
