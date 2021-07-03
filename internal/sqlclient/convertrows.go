package sqlclient

import (
	_ "github.com/go-sql-driver/mysql"
)

// ConvertRows converts rows into a table.
//
// Source: https://kylewbanks.com/blog/query-result-to-map-in-golang.
func ConvertRows(rows Rows) (Table, error) {
	cols, err := rows.Columns()
	if err != nil {
		return Table{}, err
	}

	result := Table{Columns: cols}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return Table{}, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		result.Rows = append(result.Rows, m)
	}

	return result, nil
}
