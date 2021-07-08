package sqlclient_test

import (
	"fmt"

	"github.com/nestoroprysk/TelegramBots/internal/mock"
	"github.com/nestoroprysk/TelegramBots/internal/sql"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = It("Fails to open as expected", func() {
	err := fmt.Errorf("cannot open")
	opener := mock.NewErrOpener(err)
	_, err = sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).To(MatchError("failed to connect to \"test\": cannot open"))
})

var _ = It("Fails to ping as expected", func() {
	err := fmt.Errorf("fail to ping")
	opener := mock.NewOpener(mock.PingErr(err))
	_, err = sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).To(MatchError("fail to ping"))
})

var _ = It("Queries just fine", func() {
	opener := mock.NewOpener(mock.Return(
		mock.NewRows(mock.Cols("1"), mock.Row(1)),
	).For("select 1;"))
	db, err := sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).NotTo(HaveOccurred())
	result, err := db.Query(sqlclient.Query{Statement: "select 1;"})
	Expect(err).NotTo(HaveOccurred())
	Expect(result).To(Equal(
		sql.Table{
			Columns: []string{"1"},
			Rows: []map[string]interface{}{
				{"1": 1},
			},
		},
	))
})

var _ = It("Fails to query as expected", func() {
	err := fmt.Errorf("fail to query")
	opener := mock.NewOpener(mock.ReturnErr(err).For("select 1;"))
	db, err := sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Query(sqlclient.Query{Statement: "select 1;"})
	Expect(err).To(MatchError("fail to query"))
})

var _ = It("Executes just fine", func() {
	opener := mock.NewOpener(mock.Begins(
		mock.NewTx(
			mock.ReturnTx(
				mock.Result{AffectedRows: 3, LastInsertID: 13},
			).ForTx("insert into val(v) values (?, ?, ?);", 1, 2, 3),
			mock.ReturnTx(
				mock.Result{AffectedRows: 2, LastInsertID: 42},
			).ForTx("insert into val(v) values (?, ?);", 4, 5),
		),
	))
	db, err := sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).NotTo(HaveOccurred())
	result, err := db.Exec(
		sqlclient.Query{
			Statement: "insert into val(v) values (?, ?, ?);",
			Args:      []interface{}{1, 2, 3},
		},
		sqlclient.Query{
			Statement: "insert into val(v) values (?, ?);",
			Args:      []interface{}{4, 5},
		},
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).To(Equal(sql.Result{
		RowsAffected: 5,
		LastInsertID: 42,
	}))
})

var _ = It("Tx exec error is returned nicely", func() {
	err := fmt.Errorf("tx exec errors")
	opener := mock.NewOpener(mock.Begins(
		mock.NewTx(
			mock.ReturnTxErr(err).ForTx("insert into val(v) values (1);"),
		),
	))
	db, err := sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Exec(
		sqlclient.Query{
			Statement: "insert into val(v) values (1);",
		},
	)
	Expect(err).To(MatchError("tx exec errors"))
})

var _ = It("Tx affected rows err is returned nicely", func() {
	err := fmt.Errorf("tx affected rows errors")
	opener := mock.NewOpener(mock.Begins(
		mock.NewTx(
			mock.ReturnTx(
				mock.Result{AffectedRowsErr: err},
			).ForTx("insert into val(v) values (1);"),
		),
	))
	db, err := sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Exec(
		sqlclient.Query{
			Statement: "insert into val(v) values (1);",
		},
	)
	Expect(err).To(MatchError("tx affected rows errors"))
})

var _ = It("Tx last affected ID err is returned nicely", func() {
	err := fmt.Errorf("tx last insert ID errors")
	opener := mock.NewOpener(mock.Begins(
		mock.NewTx(
			mock.ReturnTx(
				mock.Result{LastInsertIDErr: err},
			).ForTx("insert into val(v) values (1);"),
		),
	))
	db, err := sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Exec(sqlclient.Query{
		Statement: "insert into val(v) values (1);",
	})
	Expect(err).To(MatchError("tx last insert ID errors"))
})

var _ = It("Begin errors nicely", func() {
	opener := mock.NewOpener()
	db, err := sqlclient.New(sqlclient.Config{Name: "test"}, opener)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Exec(sqlclient.Query{
		Statement: "insert into val(v) values (1);",
	})
	Expect(err).To(MatchError("transaction to execute is not defined"))
})
