package sqlclient

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDB struct {
	db *sql.DB
}

var _ DB = &mysqlDB{}

type rows struct {
	r *sql.Rows
}

var _ Rows = &rows{}

// Query executes a query and returns rows.
func (db mysqlDB) Query(query string, args ...interface{}) (Rows, error) {
	result, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows{r: result}, nil
}

// Begin opens a transaction.
func (db mysqlDB) Begin() (Tx, error) {
	return db.db.Begin()
}

// Ping verifies that the connection is alive.
func (db mysqlDB) Ping() error {
	return db.db.Ping()
}

// NewOpener creates a new mysql DB opener.
func NewOpener() DBOpener {
	return func(driverName, dataSourceName string) (DB, error) {
		db, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}

		return &mysqlDB{db: db}, nil
	}
}

// Columns list columns of the resulting table.
func (r rows) Columns() ([]string, error) {
	return r.r.Columns()
}

// Next iterates to the next row.
func (r rows) Next() bool {
	return r.r.Next()
}

// Scan scans a row into dest.
func (r rows) Scan(dest ...interface{}) error {
	return r.r.Scan(dest...)
}

/// ColumnTypes returns information on columns.
func (r rows) ColumnTypes() ([]ColumnType, error) {
	cts, err := r.r.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var result []ColumnType
	for _, ct := range cts {
		result = append(result, ct)
	}

	return result, nil
}
