// +build db

package sqlclient_test

import (
	"database/sql"

	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = It("Pings", func() {
	opener := sqlclient.NewOpener()
	db, err := opener("mysql", "root:root@tcp(127.0.0.1:3306)/information_schema")
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Ping()).To(Succeed())
})

var _ = It("Selects a", func() {
	opener := sqlclient.NewOpener()
	db, err := opener("mysql", "root:root@tcp(127.0.0.1:3306)/information_schema")
	Expect(err).NotTo(HaveOccurred())
	r, err := db.Query(`select "a";`)
	Expect(err).NotTo(HaveOccurred())
	t, err := sqlclient.ConvertRows(r)
	Expect(err).NotTo(HaveOccurred())
	Expect(t).To(Equal(sqlclient.Table{
		Columns: []string{"a"},
		Rows: []map[string]interface{}{
			{"a": sql.RawBytes("a")},
		},
	}))
})
