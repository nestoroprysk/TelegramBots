package logger

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// TODO: add a mock logger

func New() Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}
