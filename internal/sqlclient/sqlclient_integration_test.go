package sqlclient_test

import (
	gosql "database/sql"

	"github.com/nestoroprysk/TelegramBots/internal/sql"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// Input defines input to query and exec integration tests.
type Input struct {
	// Setup executes queries at the setup stage.
	Setup []sqlclient.Query
	// Exec executes just after setup and optionally asserts results.
	Exec map[*[]sqlclient.Query]*Result
	// Queries maps queries to expected results just after exec.
	Queries map[*sqlclient.Query]sql.Table
	// Cleanup makes sure that nothing is left after the execution.
	Cleanup []sqlclient.Query
}

// Result is either exec result or error.
type Result struct {
	// Result asserts exec result.
	Result sql.Result
	// Err asserts exec error.
	Err bool
}

var _ = DescribeTable("Queries and executes like a pro", func(e Input) {
	db, err := sqlclient.NewIntegration()
	Expect(err).NotTo(HaveOccurred(), "run `make up` to spin a mysql instance")
	defer func() {
		_, err = db.Exec(e.Cleanup...)
		Expect(err).NotTo(HaveOccurred())
	}()

	_, err = db.Exec(e.Setup...)
	Expect(err).NotTo(HaveOccurred())

	for q, r := range e.Exec {
		actualResult, err := db.Exec(*q...)
		if r != nil {
			if r.Err {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(actualResult).To(Equal(r.Result))
			}
		}
	}

	for q, t := range e.Queries {
		actualTable, err := db.Query(*q)
		Expect(err).NotTo(HaveOccurred())
		Expect(actualTable).To(Equal(t))
	}
},
	Entry("Selects a",
		Input{
			Queries: map[*sqlclient.Query]sql.Table{
				&sqlclient.Query{Statement: `select "a"`}: sql.Table{
					Columns: []string{"a"},
					Rows: []map[string]interface{}{
						{"a": gosql.RawBytes("a")},
					},
				},
			},
		},
	),

	Entry("Selects 1",
		Input{
			Queries: map[*sqlclient.Query]sql.Table{
				&sqlclient.Query{Statement: `select 1`}: sql.Table{
					Columns: []string{"1"},
					Rows: []map[string]interface{}{
						{"1": int64(1)},
					},
				},
			},
		},
	),

	Entry("Selects 1 as foo",
		Input{
			Queries: map[*sqlclient.Query]sql.Table{
				&sqlclient.Query{Statement: `select 1 as "foo"`}: sql.Table{
					Columns: []string{"foo"},
					Rows: []map[string]interface{}{
						{"foo": int64(1)},
					},
				},
			},
		},
	),

	Entry("Selects 1 as argument",
		Input{
			Queries: map[*sqlclient.Query]sql.Table{
				&sqlclient.Query{
					Statement: `select ? as "foo"`,
					Args:      []interface{}{1},
				}: sql.Table{
					Columns: []string{"foo"},
					Rows: []map[string]interface{}{
						{"foo": int64(1)},
					},
				},
			},
		},
	),

	Entry("Creates a database, a table, and select from it",
		Input{
			Setup: []sqlclient.Query{
				{
					Statement: "create database test1;",
				},
				{
					Statement: "use test1;",
				},
				{
					Statement: "create table t(i int);",
				},
				{
					Statement: "insert into t(i) values (?), (?), (?);",
					Args:      []interface{}{1, 2, 3},
				},
			},
			Queries: map[*sqlclient.Query]sql.Table{
				&sqlclient.Query{
					Statement: `select * from t;`,
				}: sql.Table{
					Columns: []string{"i"},
					Rows: []map[string]interface{}{
						{"i": gosql.NullInt64{Int64: 1, Valid: true}},
						{"i": gosql.NullInt64{Int64: 2, Valid: true}},
						{"i": gosql.NullInt64{Int64: 3, Valid: true}},
					},
				},
			},
			Cleanup: []sqlclient.Query{
				{
					Statement: "drop database test1;",
				},
			},
		},
	),

	Entry("Rollbacks just fine",
		Input{
			Setup: []sqlclient.Query{
				{
					Statement: "create database test2;",
				},
				{
					Statement: "use test2;",
				},
				{
					Statement: "create table t(i int);",
				},
			},
			Exec: map[*[]sqlclient.Query]*Result{
				&[]sqlclient.Query{
					{
						Statement: "insert into t(i) values (?);",
						Args:      []interface{}{-1},
					},
				}: &Result{
					Result: sql.Result{
						RowsAffected: 1,
					},
				},
				&[]sqlclient.Query{
					{
						Statement: "insert into t(i) values (?), (?), (?);",
						Args:      []interface{}{1, 2, 3}, // Should rollback!
					},
					{
						Statement: "insert into t(i) values (?);",
						Args:      []interface{}{"string arg"}, // Errors...
					},
				}: &Result{Err: true},
			},
			Queries: map[*sqlclient.Query]sql.Table{
				&sqlclient.Query{
					Statement: `select * from t;`,
				}: sql.Table{
					Columns: []string{"i"},
					Rows: []map[string]interface{}{
						{"i": gosql.NullInt64{Int64: -1, Valid: true}},
					},
				},
			},
			Cleanup: []sqlclient.Query{
				{
					Statement: "drop database test2",
				},
			},
		},
	),
)
