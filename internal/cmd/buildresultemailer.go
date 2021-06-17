package cmd

import (
	"context"
	"log"
	"net/smtp"
	"os"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func BuildResultEmail(ctx context.Context, m PubSubMessage) error {
	// TODO: consider mailgun

	mustGetenv := func(e string) string {
		result := os.Getenv(e)
		if result == "" {
			panic(e + " should be set")
		}

		return result
	}

	from := mustGetenv("FROM")
	pass := mustGetenv("PASSWORD")
	to := mustGetenv("TO")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Build Result\n\n" + // TODO: inform about success/fail in subject
		string(m.Data) // TODO: trim and format data

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		// TODO: capture
		log.Printf("smtp error: %s", err)
		return err
	}

	return nil
}
