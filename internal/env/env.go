package env

import (
	"github.com/nestoroprysk/TelegramBots/internal/errorreporter"
	"github.com/nestoroprysk/TelegramBots/internal/sqlclient"
	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
)

// Env contains information from the global variables.
type Env struct {
	Telegram      telegramclient.Config `validate:"required"`
	SQL           sqlclient.Config      `validate:"required"`
	ErrorReporter errorreporter.Config  `validate:"required"`
}
