package responder_test

import (
	"fmt"
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/mock"
	"github.com/nestoroprysk/TelegramBots/internal/responder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResponder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Responder Suite")
}

var _ = It("Writes success", func() {
	rw := mock.NewResponseWriter()
	r := responder.New(&rw)
	Expect(r.Succeed("Hooray!")).To(Succeed())
	Expect(rw.Header()).To(HaveLen(1))
	Expect(rw.Header().Get("Content-Type")).To(Equal("application/json"))
	Expect(rw.Written).To(HaveLen(1))
	Expect(string(rw.Written[0])).To(MatchJSON(`{"status":"success","data":"Hooray!","message":""}`))
})

var _ = It("Writes fail", func() {
	rw := mock.NewResponseWriter()
	r := responder.New(&rw)
	Expect(r.Fail(fmt.Errorf("invalid input!!"))).To(Succeed())
	Expect(rw.Header()).To(HaveLen(1))
	Expect(rw.Header().Get("Content-Type")).To(Equal("application/json"))
	Expect(rw.Written).To(HaveLen(1))
	Expect(string(rw.Written[0])).To(MatchJSON(`{"status":"fail","data":"","message":"invalid input!!"}`))
})

var _ = It("Writes error", func() {
	rw := mock.NewResponseWriter()
	r := responder.New(&rw)
	Expect(r.Error(fmt.Errorf("bad connection"))).To(Succeed())
	Expect(rw.Header()).To(HaveLen(1))
	Expect(rw.Header().Get("Content-Type")).To(Equal("application/json"))
	Expect(rw.Written).To(HaveLen(1))
	Expect(string(rw.Written[0])).To(MatchJSON(`{"status":"error","data":"","message":"bad connection"}`))
})
