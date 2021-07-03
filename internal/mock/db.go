package mock

import (
	"database/sql"
	"fmt"

	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
)

// DB is a mock DB.
type DB struct {
	pingErr      error
	rowsPerQuery map[string]QueryClause
	tx           []Tx
}

var _ sqlclient.DB = &DB{}

// DBOption defines a mock DB.
type DBOption func(DB) DB

// QueryClause is an auxilary struct for the option Return(row).For("select 1").
type QueryClause struct {
	rows Rows
	err  error
}

// Tx is a mock transaction.
type Tx struct {
	resultPerQuery map[string]QueryTxClause
}

var _ sqlclient.Tx = &Tx{}

// TxOption defines a mock Tx.
type TxOption func(Tx) Tx

// Result is a mock result.
type Result struct {
	AffectedRows    int64
	AffectedRowsErr error
	LastInsertID    int64
	LastInsertIDErr error
}

var _ sql.Result = &Result{}

// QueryTxClause is an auxilary struct for the option Return(row).For("select 1").
type QueryTxClause struct {
	result Result
	err    error
}

// NewOpener creates a mock DB opener.
func NewOpener(opts ...DBOption) sqlclient.DBOpener {
	result := DB{rowsPerQuery: map[string]QueryClause{}}
	for _, o := range opts {
		result = o(result)
	}

	return func(driverName, dataSourceName string) (sqlclient.DB, error) {
		return result, nil
	}
}

// NewErrOpener creates an opener that errors with err.
func NewErrOpener(err error) sqlclient.DBOpener {
	return func(_, _ string) (sqlclient.DB, error) {
		return nil, err
	}
}

// PingErr return err on ping.
func PingErr(err error) DBOption {
	return func(db DB) DB {
		db.pingErr = err
		return db
	}
}

// Return defines rows to return for a query call.
func Return(r Rows) QueryClause {
	return QueryClause{rows: r}
}

// ReturnErr defines err to return for a query call.
func ReturnErr(err error) QueryClause {
	return QueryClause{err: err}
}

// For defines input for the returned rows.
func (qc QueryClause) For(query string, args ...interface{}) DBOption {
	return func(db DB) DB {
		db.rowsPerQuery[key(query, args...)] = qc
		return db
	}
}

// Query executes a query and returns rows.
func (db DB) Query(query string, args ...interface{}) (sqlclient.Rows, error) {
	result, ok := db.rowsPerQuery[key(query, args...)]
	if !ok {
		return nil, fmt.Errorf("output not defined for the query (%q) with input of len (%d)", query, len(args))
	}

	return &result.rows, result.err
}

// Begin opens a transaction.
func (db DB) Begin() (sqlclient.Tx, error) {
	if len(db.tx) == 0 {
		return nil, fmt.Errorf("transaction to execute is not defined")
	}

	result := db.tx[0]
	db.tx = db.tx[1:]

	return result, nil
}

// Ping verifies that the connection is alive.
func (db DB) Ping() error {
	return db.pingErr
}

// Begins defines input for the returned rows.
func Begins(tx Tx) DBOption {
	return func(db DB) DB {
		db.tx = append(db.tx, tx)
		return db
	}
}

// NewTx creates a mock transaction.
func NewTx(opts ...TxOption) Tx {
	result := Tx{resultPerQuery: map[string]QueryTxClause{}}
	for _, o := range opts {
		result = o(result)
	}

	return result
}

// ReturnTx defines result to return for a query call.
func ReturnTx(r Result) QueryTxClause {
	return QueryTxClause{result: r}
}

// ReturnTxErr defines err to return for a query call.
func ReturnTxErr(err error) QueryTxClause {
	return QueryTxClause{err: err}
}

// ForTx defines input for the returned rows.
func (qc QueryTxClause) ForTx(query string, args ...interface{}) TxOption {
	return func(tx Tx) Tx {
		tx.resultPerQuery[key(query, args...)] = qc
		return tx
	}
}

// Exec executes inside a transaction.
func (tx Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, ok := tx.resultPerQuery[key(query, args...)]
	if !ok {
		return nil, fmt.Errorf("output not defined for the query (%q) with input of len (%d)", query, len(args))
	}

	return &result.result, result.err
}

// Rollback does nothing.
func (Tx) Rollback() error {
	return nil
}

// Commit does nothing.
func (Tx) Commit() error {
	return nil
}

// LastInsertId returns last insert ID.
func (r Result) LastInsertId() (int64, error) {
	return r.LastInsertID, r.LastInsertIDErr
}

// RowsAffected returns rows affected.
func (r Result) RowsAffected() (int64, error) {
	return r.AffectedRows, r.AffectedRowsErr
}

// key converts query input to string.
func key(query string, args ...interface{}) string {
	result := query
	for _, a := range args {
		result += fmt.Sprintf("%v", a)
	}

	return result
}
