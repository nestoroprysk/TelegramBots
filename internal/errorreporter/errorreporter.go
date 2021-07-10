package errorreporter

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/errorreporting"
	"github.com/nestoroprysk/TelegramBots/internal/responder"
)

type errorReporter struct {
	*errorreporting.Client
	responder.Responder
}

var _ responder.Responder = &errorReporter{}

// Config defines the error reporter client.
type Config struct {
	// ProjectID links the reporter error with the project.
	ProjectID string `validate: "required"`
	// ServiceName labels the reported error as the one that belongs to the service.
	ServiceName string `validate: "required"`
}

// Wrap adds an error reporter to the responder.
func Wrap(r responder.Responder, conf Config) (responder.Responder, bool) {
	client, err := errorreporting.NewClient(context.TODO(), conf.ProjectID, errorreporting.Config{
		ServiceName: conf.ServiceName,
		OnError: func(err error) {
			log.Printf("could not report the error (%+v)", err)
		},
	})
	if err != nil {
		r.Error(fmt.Errorf("failed to initialize the error reporter: %w", err))
		return nil, false
	}

	return errorReporter{Client: client, Responder: r}, false
}

func (er errorReporter) Succeed(b interface{}) error {
	return er.Responder.Succeed(b)
}

func (er errorReporter) Fail(err error) error {
	return er.Responder.Fail(err)
}

// Error reports errors.
func (er errorReporter) Error(err error) error {
	er.Report(errorreporting.Entry{
		// TODO: Add user
		Error: err,
	})

	return er.Responder.Error(err)
}

// Close flushes all that should be reported.
func (er errorReporter) Close() error {
	err := er.Client.Close()
	log.Printf("failed to close the error reporter: %+v", err)
	return err
}
