package util_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAdminSQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util Test Suite")
}

var _ = It("True", func() {
	Expect(true).To(BeTrue())
})
