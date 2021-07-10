package logger

import (
	"log"
	"os"

	"github.com/nestoroprysk/TelegramBots/internal/responder"
)

// logger implements the Responder interface.
type logger struct {
	success *log.Logger
	fail    *log.Logger
	err     *log.Logger
}

var _ responder.Responder = &logger{}

// New creates a new responder.
func New() responder.Responder {
	return &logger{
		success: log.New(os.Stdout, "[SUCCESS]", log.LstdFlags),
		fail:    log.New(os.Stdout, "[FAIL]", log.LstdFlags),
		err:     log.New(os.Stderr, "[ERROR]", log.LstdFlags),
	}
}

// Succeed logs success.
func (l logger) Succeed(b interface{}) error {
	l.success.Println(b)
	return nil
}

// Fail logs fail.
func (l logger) Fail(err error) error {
	l.fail.Println(err)
	return nil
}

// Error logs error.
func (l logger) Error(err error) error {
	l.err.Println(err)
	return nil
}

// Close does nothing.
func (l logger) Close() error {
	return nil
}
