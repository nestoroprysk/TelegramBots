package mock

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
)

// Rows is a mock rows.
type Rows struct {
	cols []string
	left []row
}

var _ sqlclient.Rows = &Rows{}

// RowsOption defines rows.
type RowsOption func(Rows) Rows

// row is the result of the scan call/
type row struct {
	cols []interface{}
	err  error
}

type columnType struct {
	t reflect.Type
}

var _ sqlclient.ColumnType = &columnType{}

// NewRows creates mock rows defined by options.
func NewRows(opts ...RowsOption) Rows {
	result := Rows{}
	for _, o := range opts {
		result = o(result)
	}

	return result
}

// Cols sets columns.
func Cols(cols ...string) RowsOption {
	return func(r Rows) Rows {
		r.cols = cols
		return r
	}
}

//  Row adds row.
func Row(cols ...interface{}) RowsOption {
	return func(r Rows) Rows {
		r.left = append(r.left, row{cols: cols})
		return r
	}
}

//  RowErr adds row error.
func RowErr(err error) RowsOption {
	return func(r Rows) Rows {
		r.left = append(r.left, row{err: err})
		return r
	}
}

// Columns list columns of the resulting table.
func (r Rows) Columns() ([]string, error) {
	if len(r.cols) == 0 {
		return nil, fmt.Errorf("no columns")
	}

	return r.cols, nil
}

// Next iterates to the next row.
func (r Rows) Next() bool {
	return len(r.left) > 0
}

// Scan scans a row into dest.
func (r *Rows) Scan(dest ...interface{}) error {
	if len(r.left) == 0 {
		return sql.ErrNoRows
	}

	for i, d := range dest {
		row := r.left[0]

		if row.err != nil {
			return row.err
		}

		if len(row.cols) <= i {
			break
		}

		reflect.ValueOf(d).Elem().Set(reflect.ValueOf(row.cols[i]))
	}

	r.left = r.left[1:]

	return nil
}

/// ColumnTypes returns information on columns.
func (r Rows) ColumnTypes() ([]sqlclient.ColumnType, error) {
	if len(r.left) == 0 {
		return nil, nil
	}

	if err := r.left[0].err; err != nil {
		return nil, err
	}

	var result []sqlclient.ColumnType
	for _, c := range r.left[0].cols {
		result = append(result, &columnType{t: reflect.TypeOf(c)})
	}

	return result, nil
}

// ScanType returns a Go type suitable for scanning into using Rows.Scan.
func (ct columnType) ScanType() reflect.Type {
	return ct.t
}
