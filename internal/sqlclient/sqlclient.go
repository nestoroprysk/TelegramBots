package sqlclient

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/nestoroprysk/TelegramBots/internal/env"
)

// SQLClient selects and executes.
type SQLClient interface {
	// Query selects into the table.
	Query(Query) (Table, error)
	// Exec executes queries inside a transaction.
	Exec(...Query) (Result, error)
}

var _ SQLClient = &sqlClient{}

// Query is a SQL query with arguments.
type Query struct {
	// Statement is a SQL statement.
	Statement string
	// Args are arguments to the statement.
	Args []interface{}
}

// Table is the result of the select statement.
type Table struct {
	// Columns is a list of columns.
	Columns []string
	// Rows is a list of maps from a column name to the result.
	Rows []map[string]interface{}
}

// Result is the result of a non-select SQL statement.
type Result struct {
	// RowsAffected counts affected rows.
	RowsAffected int64
	// LastInsertID indicates the last insertID.
	LastInsertID int64
}

// Rows is the result of a select.
type Rows interface {
	// Columns list columns of the resulting table.
	Columns() ([]string, error)
	// Next iterates to the next row.
	Next() bool
	// Scan scans a row into dest.
	Scan(dest ...interface{}) error
	/// ColumnTypes returns information on columns.
	ColumnTypes() ([]ColumnType, error)
}

// ColumnType provides information about columns.
type ColumnType interface {
	// ScanType returns a Go type suitable for scanning into using Rows.Scan.
	ScanType() reflect.Type
}

// DB queries and begins transactions.
type DB interface {
	// Query executes a query and returns rows.
	Query(query string, args ...interface{}) (Rows, error)
	// Begin opens a transaction.
	Begin() (Tx, error)
	// Ping verifies that the connection is alive.
	Ping() error
}

// Tx is a transaction.
type Tx interface {
	// Exec executes inside a transaction.
	Exec(query string, args ...interface{}) (sql.Result, error)
	// Rollback undoes all the queries that are a part of the transaction.
	Rollback() error
	// Commit commits the transaction.
	Commit() error
}

var _ Tx = &sql.Tx{}

// DBOpener opens a database connection (e.g., sql.Open).
type DBOpener func(driverName, dataSourceName string) (DB, error)

type sqlClient struct {
	db DB
}

// New creates an SQL client.
func New(conf env.DB, open DBOpener) (SQLClient, error) {
	const socketDir = "/cloudsql"
	dbURI := fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", conf.User, conf.Password, socketDir, conf.InstanceConnectionName, conf.Name)

	db, err := open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %q: %w", conf.Name, err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &sqlClient{db: db}, nil
}

// Query selects into the table.
func (sc sqlClient) Query(q Query) (Table, error) {
	result, err := sc.db.Query(q.Statement, q.Args...)
	if err != nil {
		return Table{}, err
	}

	return ConvertRows(result)
}

// Exec executes queries inside a transaction.
//
// The result is the sum of the affected rows.
func (sc sqlClient) Exec(qs ...Query) (Result, error) {
	tx, err := sc.db.Begin()
	if err != nil {
		return Result{}, err
	}

	var (
		result       int64
		lastInsertID int64
	)

	for _, q := range qs {
		r, err := tx.Exec(q.Statement, q.Args...)
		if err != nil {
			_ = tx.Rollback()
			return Result{}, err
		}

		rowsAffected, err := r.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return Result{}, err
		}

		result += rowsAffected

		lastInsertID, err = r.LastInsertId()
		if err != nil {
			_ = tx.Rollback()
			return Result{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Result{}, err
	}

	return Result{
		RowsAffected: result,
		LastInsertID: lastInsertID,
	}, nil
}
