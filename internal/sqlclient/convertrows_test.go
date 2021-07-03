package sqlclient_test

import (
	"fmt"
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/mock"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSQLClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQL Client Suite")
}

var _ = It("Fails to convert rows if columns error", func() {
	r := mock.NewRows()
	_, err := sqlclient.ConvertRows(&r)
	Expect(err).To(MatchError("no columns"))
})

var _ = It("Fails to convert rows if scan errors", func() {
	err := fmt.Errorf("oh no!")
	r := mock.NewRows(mock.Cols("name"), mock.RowErr(err))
	_, err = sqlclient.ConvertRows(&r)
	Expect(err).To(MatchError("oh no!"))
})

var _ = It("Converts columns", func() {
	r := mock.NewRows(mock.Cols("name", "age"))
	result, err := sqlclient.ConvertRows(&r)
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Columns).To(ConsistOf("name", "age"))
	Expect(result.Rows).To(BeEmpty())
})

var _ = It("Converts rows", func() {
	r := mock.NewRows(mock.Cols("name", "age"),
		mock.Row("Claire", 44),
		mock.Row("Rossie", 34),
	)
	result, err := sqlclient.ConvertRows(&r)
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Columns).To(ConsistOf("name", "age"))
	Expect(result.Rows).To(ConsistOf(
		map[string]interface{}{
			"name": "Claire",
			"age":  44,
		},
		map[string]interface{}{
			"name": "Rossie",
			"age":  34,
		},
	))
})
