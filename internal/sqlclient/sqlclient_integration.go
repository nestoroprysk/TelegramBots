// +build integration

package sqlclient

import (
	"database/sql"
)

// NewIntegration creates an SQL client for integration tests.
func NewIntegration() (SQLClient, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/information_schema")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &sqlClient{db: mysqlDB{db: db}}, nil
}
