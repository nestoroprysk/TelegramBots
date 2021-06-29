package mock_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nestoroprysk/TelegramBots/internal/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = It("Poster returns empty output by default", func() {
	p := mock.NewPoster()
	result, err := p.PostForm("eh", nil)
	Expect(err).NotTo(HaveOccurred())
	rc := mock.NewReadCloser(nil)
	Expect(result).To(Equal(&http.Response{Body: &rc}))
})

var _ = It("Poster errors as expected", func() {
	p := mock.NewPoster(mock.PostFormError(fmt.Errorf("foo")))
	_, err := p.PostForm("eh", nil)
	Expect(err).To(MatchError("foo"))
})

var _ = It("Poster return status code", func() {
	p := mock.NewPoster(mock.PostFormStatusCode(12))
	resp, err := p.PostForm("eh", nil)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(12))
})

var _ = It("Poster sets body", func() {
	rc := mock.NewReadCloser([]byte("abc"))
	p := mock.NewPoster(mock.PostFormBody(&rc))
	resp, err := p.PostForm("eh", nil)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(0))
	result, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).To(BeEquivalentTo("abc"))
})

var _ = It("Multitle options work", func() {
	rc := mock.NewReadCloser([]byte("abc"))
	p := mock.NewPoster(mock.PostFormBody(&rc), mock.PostFormStatusCode(10))
	resp, err := p.PostForm("eh", nil)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(10))
	result, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).To(BeEquivalentTo("abc"))
})
