package errorreporter

import (
	"context"
	"log"

	"cloud.google.com/go/errorreporting"
)

// ErrorReporter reports error.
type ErrorReporter interface {
	// Error reports errors.
	Error(error)
	// Close flushes all that should be reported.
	Close() error
}

type errorReporter struct {
	*errorreporting.Client
}

var _ ErrorReporter = &errorReporter{}

// Config defines the error reporter client.
type Config struct {
	// ProjectID links the reporter error with the project.
	ProjectID string `validate: "required"`
	// ServiceName labels the reported error as the one that belongs to the service.
	ServiceName string `validate: "required"`
}

// New creates a new error reporter.
func New(conf Config) (ErrorReporter, error) {
	client, err := errorreporting.NewClient(context.TODO(), conf.ProjectID, errorreporting.Config{
		ServiceName: conf.ServiceName,
		OnError: func(err error) {
			// TODO: Inject our err client
			log.Printf("could not report the error error (%+v)", err)
		},
	})
	if err != nil {
		return nil, err
	}

	return errorReporter{Client: client}, nil
}

// Error reports errors.
func (er errorReporter) Error(err error) {
	er.Report(errorreporting.Entry{
		// TODO: Add user
		Error: err,
	})
}

// Close flushes all that should be reported.
func (er errorReporter) Close() error {
	return er.Client.Close()
}
