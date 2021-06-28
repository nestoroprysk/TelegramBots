package util_test

import (
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util Suite")
}

var _ = DescribeTable("Formats", func(r sqlclient.Response, expectedResult string) {
	result := util.Format(r)
	Expect(result).To(Equal(expectedResult))
},
	Entry("Formats response",
		sqlclient.Response{
			Columns: []string{"name", "age"},
			Rows: [][]string{
				{"John", "20"},
				{"Bart", "15"},
			},
		},
		"name,age\nJohn,20\nBart,15",
	),
	Entry("Formats empty response",
		sqlclient.Response{},
		"",
	),
	Entry("Formats column only response",
		sqlclient.Response{Columns: []string{"a", "b"}},
		"a,b",
	),
	Entry("Formats response with not enough columns in row",
		sqlclient.Response{
			Columns: []string{"name", "age", "date"},
			Rows: [][]string{
				{"John", "20"},
			},
		},
		"name,age,date\nJohn,20,",
	),
	Entry("Formats response with too many columns in row",
		sqlclient.Response{
			Columns: []string{"name"},
			Rows: [][]string{
				{"John", "20"},
			},
		},
		"name,\nJohn,20",
	),
)
