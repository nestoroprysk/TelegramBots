package sql_test

import (
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQL Suite")
}

var _ = DescribeTable("Formats", func(t sql.Table, expectedResult string) {
	result := sql.FormatTable(t)
	Expect(result).To(Equal(expectedResult))
},
	Entry("Formats response",
		sql.Table{
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
		`+------+-----+
| NAME | AGE |
+------+-----+
| John |  20 |
| Bart |  15 |
+------+-----+`,
	),
	Entry("Formats empty response",
		sql.Table{},
		"",
	),
	Entry("Formats column only response",
		sql.Table{Columns: []string{"a", "b"}},
		`+---+---+
| A | B |
+---+---+
+---+---+`,
	),
	Entry("Formats response with not enough columns in row",
		sql.Table{
			Columns: []string{"name", "age", "date"},
			Rows: []map[string]interface{}{
				{
					"name": "John",
					"age":  20,
				},
			},
		},
		`+------+-----+------+
| NAME | AGE | DATE |
+------+-----+------+
| John |  20 |      |
+------+-----+------+`,
	),
	Entry("Formats response with too many columns in row",
		sql.Table{
			Columns: []string{"name"},
			Rows: []map[string]interface{}{
				{
					"name": "John",
					"age":  20,
				},
			},
		},
		`+------+
| NAME |
+------+
| John |
+------+`,
	),
)
