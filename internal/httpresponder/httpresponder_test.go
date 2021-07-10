package httpresponder_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/httpresponder"
	"github.com/nestoroprysk/TelegramBots/internal/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHTTPResponder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Responder Suite")
}

var _ = It("Writes success", func() {
	rw := mock.NewResponseWriter()
	r := httpresponder.Wrap(mock.NoOpResponder(), &rw)
	Expect(r.Succeed("Hooray!")).To(Succeed())
	Expect(rw.Header()).To(HaveLen(1))
	Expect(rw.Header().Get("Content-Type")).To(Equal("application/json"))
	Expect(rw.Written).To(HaveLen(1))
	Expect(string(rw.Written[0])).To(MatchJSON(`{"status":"success","data":"Hooray!","message":"","status_code": 200}`))
	Expect(rw.StatusCode).To(Equal(http.StatusOK))
})

var _ = It("Writes fail", func() {
	rw := mock.NewResponseWriter()
	r := httpresponder.Wrap(mock.NoOpResponder(), &rw)
	Expect(r.Fail(fmt.Errorf("invalid input!!"))).To(Succeed())
	Expect(rw.Header()).To(HaveLen(1))
	Expect(rw.Header().Get("Content-Type")).To(Equal("application/json"))
	Expect(rw.Written).To(HaveLen(1))
	Expect(string(rw.Written[0])).To(MatchJSON(`{"status":"fail","data":null,"message":"invalid input!!","status_code": 400}`))
	Expect(rw.StatusCode).To(Equal(http.StatusBadRequest))
})

var _ = It("Writes error", func() {
	rw := mock.NewResponseWriter()
	r := httpresponder.Wrap(mock.NoOpResponder(), &rw)
	Expect(r.Error(fmt.Errorf("bad connection"))).To(Succeed())
	Expect(rw.Header()).To(HaveLen(1))
	Expect(rw.Header().Get("Content-Type")).To(Equal("application/json"))
	Expect(rw.Written).To(HaveLen(1))
	Expect(string(rw.Written[0])).To(MatchJSON(`{"status":"error","data":null,"message":"bad connection","status_code": 500}`))
	Expect(rw.StatusCode).To(Equal(http.StatusInternalServerError))
})
