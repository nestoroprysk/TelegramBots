package telegram_test

import (
	"testing"

	"github.com/nestoroprysk/TelegramBots/internal/mock"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTelegramt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telegram Suite")
}

var _ = It("Update errors nicely if fails to parse", func() {
	rc := mock.NewReadCloser(nil)
	_, err := telegram.ParseUpdate(&rc)
	Expect(err).To(MatchError("could not decode an incoming update: EOF"))
})

var _ = It("Update gets parsed correctly", func() {
	rc := mock.NewReadCloser([]byte(`{
  "message": {
    "message_id": 345,
    "from": {
      "id": 381126698,
      "is_bot": false,
      "first_name": "Nestor",
      "username": "npetro",
      "language_code": "en"
    },
    "chat": {
      "id": 381126698,
      "first_name": "Nestor",
      "username": "npetro",
      "type": "private"
    },
    "date": 1623779885,
    "text": "select 9;"
  }
	}`))
	u, err := telegram.ParseUpdate(&rc)
	Expect(err).NotTo(HaveOccurred())
	Expect(u).To(Equal(telegram.Update{
		Message: telegram.Message{
			Text: "select 9;",
			Chat: telegram.Chat{
				ID: 381126698,
			},
			From: telegram.User{
				ID: 381126698,
			},
		},
	}))
})
