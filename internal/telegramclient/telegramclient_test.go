package telegramclient_test

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nestoroprysk/TelegramBots/internal/mock"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
)

func TestTelegramClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telegram Client Suite")
}

var _ = It("Errors nicely if post fails", func() {
	conf := telegramclient.Config{}
	chatID := 10
	p := mock.NewPoster(mock.PostFormError(fmt.Errorf("oh no!")))
	r := telegramclient.Wrap(mock.ForwardFailErrResponder(), conf, chatID, p)
	Expect(r.Succeed("abc")).To(MatchError(`error when posting text to the chat "10": oh no!`))
})

var _ = It("Errors nicely if the response code is not ok", func() {
	conf := telegramclient.Config{}
	chatID := 10
	p := mock.NewPoster(mock.PostFormStatusCode(401))
	r := telegramclient.Wrap(mock.ForwardFailErrResponder(), conf, chatID, p)
	Expect(r.Succeed("abc")).To(MatchError(`expecting status code 200 for the Telegram response; got 401`))
})

var _ = It("Errors nicely if fails to parse the body", func() {
	conf := telegramclient.Config{}
	chatID := 10
	p := mock.NewPoster(mock.PostFormStatusCode(http.StatusOK))
	r := telegramclient.Wrap(mock.ForwardFailErrResponder(), conf, chatID, p)
	Expect(r.Succeed("abc")).To(MatchError("could not decode an incoming response: EOF"))
})

var _ = It("Errors nicely if the response is not ok", func() {
	conf := telegramclient.Config{}
	chatID := 10
	rc := mock.NewReadCloser([]byte(`{"ok":false}`))
	p := mock.NewPoster(mock.PostFormStatusCode(http.StatusOK), mock.PostFormBody(&rc))
	r := telegramclient.Wrap(mock.ForwardFailErrResponder(), conf, chatID, p)
	Expect(r.Succeed("abc")).To(MatchError("expecting ok; got {Ok:false Result:{Text: Chat:{ID:0} From:{ID:0}} ErrorCode:0 Description:}"))
})

var _ = It("Errors nicely if the error code is not zero", func() {
	conf := telegramclient.Config{}
	chatID := 10
	rc := mock.NewReadCloser([]byte(`{"ok":true, "error_code":1}`))
	p := mock.NewPoster(mock.PostFormStatusCode(http.StatusOK), mock.PostFormBody(&rc))
	r := telegramclient.Wrap(mock.ForwardFailErrResponder(), conf, chatID, p)
	Expect(r.Succeed("abc")).To(MatchError("expecting zero exit code; got {Ok:true Result:{Text: Chat:{ID:0} From:{ID:0}} ErrorCode:1 Description:}"))
})

var _ = It("Succeeds", func() {
	conf := telegramclient.Config{}
	chatID := 10
	rc := mock.NewReadCloser([]byte(`{"ok":true, "result":{"text":"123"}}`))
	p := mock.NewPoster(mock.PostFormStatusCode(http.StatusOK), mock.PostFormBody(&rc))
	var result interface{}
	r := telegramclient.Wrap(mock.SucceedResponder(func(d interface{}) error {
		result = d
		return nil
	}), conf, chatID, p)
	Expect(r.Succeed("abc")).To(Succeed())
	Expect(result).To(Equal(telegram.Response{Ok: true, Result: telegram.Message{Text: "123"}}))
})
