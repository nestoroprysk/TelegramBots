package mock_test

import (
	"github.com/nestoroprysk/TelegramBots/internal/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = It("ReadCloser is open by default", func() {
	rc := mock.NewReadCloser(nil)
	Expect(rc.IsClosed).To(BeFalse())
})

var _ = It("ReadCloser is closed correctly", func() {
	rc := mock.NewReadCloser(nil)
	Expect(rc.Close()).To(Succeed())
	Expect(rc.IsClosed).To(BeTrue())
})

var _ = It("ReadCloser reads correctly", func() {
	rc := mock.NewReadCloser([]byte("abc"))
	result := make([]byte, 3, 3)
	n, err := rc.Read(result)
	Expect(err).To(BeNil())
	Expect(n).To(Equal(3))
	Expect(result).To(BeEquivalentTo("abc"))
})
