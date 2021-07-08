package sql_test

import (
	"github.com/nestoroprysk/TelegramBots/internal/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = It("Formats exec result", func() {
	Expect(sql.FormatResult(sql.Result{RowsAffected: 1})).To(Equal("Query OK, 1 row affected"))
})
