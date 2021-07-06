package mock_test

import (
	"net/http"
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mock Suite")
}

var _ = It("Returns and sets header successfully", func() {
	rw := mock.NewResponseWriter()
	h := rw.Header()
	h.Set("key", "value")
	Expect(rw.Header()).To(HaveLen(1))
	Expect(rw.Header().Get("key")).To(Equal("value"))
})

var _ = It("Write adds a line to written and returns its len", func() {
	rw := mock.NewResponseWriter()
	l, err := rw.Write([]byte("abc"))
	Expect(err).To(BeNil())
	Expect(l).To(Equal(3))
	Expect(rw.Written).To(HaveLen(1))
	Expect(rw.Written[0]).To(BeEquivalentTo("abc"))
})

var _ = It("WriteHeader sets status code", func() {
	rw := mock.NewResponseWriter()
	Expect(rw.StatusCode).To(BeZero())
	rw.WriteHeader(http.StatusOK)
	Expect(rw.StatusCode).To(Equal(http.StatusOK))
})
