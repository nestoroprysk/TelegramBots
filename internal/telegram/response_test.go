package telegram_test

import (
	"github.com/nestoroprysk/TelegramBots/internal/mock"
	"github.com/nestoroprysk/TelegramBots/internal/telegram"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = It("Response errors nicely if fails to parse", func() {
	rc := mock.NewReadCloser(nil)
	_, err := telegram.ParseResponse(&rc)
	Expect(err).To(MatchError("could not decode an incoming response: EOF"))
})

var _ = It("Resonse gets parsed correctly", func() {
	rc := mock.NewReadCloser([]byte(`{
  "ok": true,
  "result": {
    "message_id": 348,
    "from": {
      "id": 1680972219,
      "is_bot": true,
      "first_name": "Expenses",
      "username": "ExpensesNewSQLBot"
    },
    "chat": {
      "id": 381126698,
      "first_name": "Nestor",
      "username": "npetro",
      "type": "private"
    },
    "date": 1623781374,
    "text": "9\n9"
  }
	}`))
	r, err := telegram.ParseResponse(&rc)
	Expect(err).NotTo(HaveOccurred())
	Expect(r).To(Equal(telegram.Response{
		Ok: true,
		Result: telegram.Message{
			Text: "9\n9",
			Chat: telegram.Chat{
				ID: 381126698,
			},
			From: telegram.User{
				ID: 1680972219,
			},
		},
	}))
})

var _ = It("Error resonse gets parsed correctly as well", func() {
	rc := mock.NewReadCloser([]byte(`{"ok":false,"error_code":404,"description":"Not Found"}`))
	r, err := telegram.ParseResponse(&rc)
	Expect(err).NotTo(HaveOccurred())
	Expect(r).To(Equal(telegram.Response{
		Ok:          false,
		ErrorCode:   404,
		Description: "Not Found",
	}))
})
