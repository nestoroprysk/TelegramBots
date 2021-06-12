package lowlevel_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAdminSQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lowlevel Test Suite")
}

var _ = It("True", func() {
	Expect(true).To(BeTrue())
})
