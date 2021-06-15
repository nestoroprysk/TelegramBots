package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type SQLClient interface {
	Send(request string) (string, error)
}

type sqlClient struct {
	name                   string
	user                   string
	password               string
	instanceConnectionName string
}

// TODO: add interface assertion

// TODO: add mock SQL

// NewSQLClient creates an SQL client.
func NewSQLClient(conf DB) SQLClient {
	return &sqlClient{
		name:                   conf.Name,
		user:                   conf.User,
		password:               conf.Password,
		instanceConnectionName: conf.InstanceConnectionName,
	}
}

// Send sends an SQL query.
func (sc sqlClient) Send(query string) (string, error) {
	log.Printf("executing %q", query)

	const socketDir = "/cloudsql"

	dbURI := fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", sc.user, sc.password, socketDir, sc.instanceConnectionName, sc.name)

	// TODO: inject in-memory SQL for testing

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		// TODO: capture
		return "", fmt.Errorf("failed to connect to %q: %s", sc.name, err.Error())
	}

	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute %q: %v", query, err.Error())
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		// TODO: capture
		return "", fmt.Errorf("failed to get columns: %s", err.Error())
	}

	// TODO: move it out of here

	var all []string
	all = append(all, strings.Join(cols, " "))

	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return "", fmt.Errorf("rows.Scan: %v", err)
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		all = append(all, strings.Join(result, " "))
	}

	return strings.Join(all, "\n"), nil
}
