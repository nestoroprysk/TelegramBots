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

var _ = DescribeTable("Formats", func(t sqlclient.Table, expectedResult string) {
	result := util.Format(t)
	Expect(result).To(Equal(expectedResult))
},
	Entry("Formats response",
		sqlclient.Table{
			Columns: []string{"name", "age"},
			Rows: []map[string]interface{}{
				{
					"name": "John",
					"age":  20,
				},
				{
					"name": "Bart",
					"age":  15,
				},
			},
		},
		"name,age\nJohn,20\nBart,15",
	),
	Entry("Formats empty response",
		sqlclient.Table{},
		"",
	),
	Entry("Formats column only response",
		sqlclient.Table{Columns: []string{"a", "b"}},
		"a,b",
	),
	Entry("Formats response with not enough columns in row",
		sqlclient.Table{
			Columns: []string{"name", "age", "date"},
			Rows: []map[string]interface{}{
				{
					"name": "John",
					"age":  20,
				},
			},
		},
		"name,age,date\nJohn,20,",
	),
	Entry("Formats response with too many columns in row",
		sqlclient.Table{
			Columns: []string{"name"},
			Rows: []map[string]interface{}{
				{
					"name": "John",
					"age":  20,
				},
			},
		},
		"name\nJohn",
	),
)
