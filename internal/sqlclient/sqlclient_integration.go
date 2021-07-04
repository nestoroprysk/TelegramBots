package sqlclient

import (
	"database/sql"
	"fmt"
	"os"
)

// NewIntegration creates an SQL client for integration tests.
func NewIntegration() (SQLClient, error) {
	host := "127.0.0.1"
	if os.Getenv("CI") != "" {
		host = "mysql"
	}

	db, err := sql.Open("mysql", fmt.Sprintf("root:root@tcp(%s:3306)/information_schema", host))
	if err != nil {
		return nil, err
	}

	result := mysqlDB{db: db}

	if err := result.Ping(); err != nil {
		return nil, err
	}

	return &sqlClient{db: result}, nil
}
