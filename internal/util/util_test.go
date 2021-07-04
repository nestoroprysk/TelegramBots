package util_test

import (
	"github.com/nestoroprysk/TelegramBots/internal/util"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Formats", func(str string, count int, expectedResult string) {
	result := util.Pluralize(str, count)
	Expect(result).To(Equal(expectedResult))
},
	Entry("Singular if one",
		"row", 1, "row",
	),
	Entry("Plural if many",
		"row", 2, "rows",
	),
	Entry("Plural if zero",
		"row", 0, "rows",
	),
	Entry("Plural if negative",
		"row", -1, "rows",
	),
)
