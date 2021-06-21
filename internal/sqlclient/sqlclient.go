package sqlclient

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nestoroprysk/TelegramBots/internal/env"
)

// SQLClient sends an SQL query and returns a response.
type SQLClient interface {
	Send(request string) (Response, error)
}

// Response is a generic SQL response.
type Response struct {
	Columns []string
	Rows    [][]string
}

type sqlClient struct {
	*sql.DB
}

// TODO: add interface assertion

// TODO: add mock SQL

// NewSQLClient creates an SQL client.
func New(conf env.DB) (SQLClient, error) {
	const socketDir = "/cloudsql"
	dbURI := fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", conf.User, conf.Password, socketDir, conf.InstanceConnectionName, conf.Name)

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %q: %w", conf.Name, err)
	}

	return &sqlClient{DB: db}, nil
}

// Send sends an SQL query.
func (sc sqlClient) Send(query string) (Response, error) {
	// TODO: use some library to parse query and validate if
	//       if all is fine, select either query or exec
	//       output different result for query or exec
	//       or even input the parsed query
	rows, err := sc.Query(query)
	if err != nil {
		return Response{}, fmt.Errorf("failed to execute %q: %w", query, err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		// TODO: capture
		return Response{}, fmt.Errorf("failed to get columns: %w", err)
	}

	rawResponse := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols)) // A temporary interface{} slice.
	for i, _ := range rawResponse {
		dest[i] = &rawResponse[i] // Put pointers to each string in the interface slice.
	}

	result := Response{
		Columns: cols,
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			// TODO: capture
			return Response{}, fmt.Errorf("failed to scan the result: %w", err)
		}

		result.Rows = append(result.Rows, []string{})
		for _, item := range rawResponse {
			result.Rows[len(result.Rows)-1] = append(result.Rows[len(result.Rows)-1], string(item))
		}
	}

	return result, nil
}
