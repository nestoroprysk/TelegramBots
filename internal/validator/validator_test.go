package validator_test

import (
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/validator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validator Suite")
}

type Test struct {
	RequiredString string `json:"text" validate:"required"`
	PositiveInt    int    `json:"id" validate:"gt=0"`
	OptionalString string `json:"text"`
}

var _ = DescribeTable("Validates the test struct", func(t Test, succeed bool) {
	v := validator.New()
	result := v.Struct(t)
	if succeed {
		Expect(result).To(Succeed())
	} else {
		Expect(result).NotTo(BeNil())
	}
},
	Entry("Returns nil when all is fine",
		Test{RequiredString: "s", PositiveInt: 1},
		true, /* succeed */
	),
	Entry("Returns nil when optional string is indicated as well",
		Test{RequiredString: "s", PositiveInt: 1, OptionalString: "s"},
		true, /* succeed */
	),
	Entry("Errors if positive int is not positive",
		Test{RequiredString: "s", PositiveInt: -1},
		false, /* succeed */
	),
	Entry("Returns if required string is empty",
		Test{RequiredString: "", PositiveInt: 1},
		false, /* succeed */
	),
)
