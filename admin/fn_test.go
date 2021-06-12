package admin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAdminSQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AdminSQL")
}

var _ = It("False", func() {
	Expect(true).To(BeTrue())
})
