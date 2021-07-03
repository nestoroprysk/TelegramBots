package mock_test

import (
	"fmt"

	"github.com/nestoroprysk/TelegramBots/internal/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	driver = ""
	dsn    = ""
)

var _ = It("Err DB opener errors", func() {
	err := fmt.Errorf("cannot open")
	_, err = mock.NewErrOpener(err)(driver, dsn)
	Expect(err).To(MatchError("cannot open"))
})

var _ = It("Ping errors as expected", func() {
	err := fmt.Errorf("ping failed")
	db, err := mock.NewOpener(mock.PingErr(err))(driver, dsn)
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Ping()).To(MatchError("ping failed"))
})

var _ = It("Returns expected rows", func() {
	db, err := mock.NewOpener(mock.Return(
		mock.NewRows(
			mock.Cols("1"),
			mock.Row(1),
		),
	).For("select ?;", "1"))(driver, dsn)
	Expect(err).NotTo(HaveOccurred())
	result, err := db.Query("select ?;", "1")
	Expect(err).NotTo(HaveOccurred())
	cols, err := result.Columns()
	Expect(err).NotTo(HaveOccurred())
	Expect(cols).To(ConsistOf("1"))
	var i int
	Expect(result.Next()).To(BeTrue())
	Expect(result.Scan(&i)).To(Succeed())
	Expect(i).To(Equal(1))
	Expect(result.Next()).To(BeFalse())
})

var _ = It("Errors if query output not defined", func() {
	db, err := mock.NewOpener()(driver, dsn)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Query("select 2;")
	Expect(err).To(MatchError("output not defined for the query (\"select 2;\") with input of len (0)"))
})

var _ = It("Query errors as expected", func() {
	err := fmt.Errorf("fail to select")
	db, err := mock.NewOpener(mock.ReturnErr(err).For("select 1;"))(driver, dsn)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Query("select 1;")
	Expect(err).To(MatchError("fail to select"))
})

var _ = It("Begin errors if tx is not defined", func() {
	db, err := mock.NewOpener()(driver, dsn)
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Begin()
	Expect(err).To(MatchError("transaction to execute is not defined"))
})

var _ = It("Executes inside a transaction", func() {
	db, err := mock.NewOpener(mock.Begins(
		mock.NewTx(
			mock.ReturnTx(
				mock.Result{AffectedRows: 3, LastInsertID: 42},
			).ForTx("insert into vals(v) values (1, 2, 3);"),
		),
	))(driver, dsn)
	Expect(err).NotTo(HaveOccurred())
	tx, err := db.Begin()
	Expect(err).NotTo(HaveOccurred())
	result, err := tx.Exec("insert into vals(v) values (1, 2, 3);")
	Expect(err).NotTo(HaveOccurred())
	rowsAffected, err := result.RowsAffected()
	Expect(err).NotTo(HaveOccurred())
	Expect(rowsAffected).To(BeEquivalentTo(3))
	id, err := result.LastInsertId()
	Expect(err).NotTo(HaveOccurred())
	Expect(id).To(BeEquivalentTo(42))
})

var _ = It("Tx errors if exec is not defined", func() {
	tx := mock.NewTx()
	_, err := tx.Exec("insert into vals(v) values (1);")
	Expect(err).To(MatchError("output not defined for the query (\"insert into vals(v) values (1);\") with input of len (0)"))
})

var _ = It("Tx does nothing on commit", func() {
	tx := mock.NewTx()
	Expect(tx.Commit()).To(Succeed())
})

var _ = It("Tx does nothing on rollback", func() {
	tx := mock.NewTx()
	Expect(tx.Rollback()).To(Succeed())
})

var _ = It("Tx errors as expected", func() {
	err := fmt.Errorf("fail to execute")
	tx := mock.NewTx(mock.ReturnTxErr(err).ForTx("insert into vals(v) values (1);"))
	_, err = tx.Exec("insert into vals(v) values (1);")
	Expect(err).To(MatchError("fail to execute"))
})
