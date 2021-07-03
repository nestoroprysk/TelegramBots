package sqlclient

import "database/sql"

type mysqlDB struct {
	*sql.DB
}

// Query executes a query and returns rows.
func (db mysqlDB) Query(query string, args ...interface{}) (Rows, error) {
	return db.Query(query, args...)
}

// Begin opens a transaction.
func (db mysqlDB) Begin() (Tx, error) {
	return db.Begin()
}

// Ping verifies that the connection is alive.
func (db mysqlDB) Ping() error {
	return db.DB.Ping()
}

var _ DB = &mysqlDB{}

// NewOpener creates a new mysql DB opener.
func NewOpener() DBOpener {
	return func(driverName, dataSourceName string) (DB, error) {
		db, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}

		return &mysqlDB{DB: db}, nil
	}
}
