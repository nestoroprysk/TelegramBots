package env

import "github.com/nestoroprysk/TelegramBots/internal/errorreporter"

// Env contains information from the global variables.
type Env struct {
	Telegram      `validate:"required"`
	DB            `validate:"required"`
	ErrorReporter errorreporter.Config `validate:"required"`
}

// TODO: Move structs to proper packages

// Telegram is the Telegram environment.
type Telegram struct {
	// Token is a Telegram bot token.
	Token string `validate:"required"`
	// AdminID is an ID of the admin Telegram user.
	AdminID int `validate:"gt=0"`
}

// DB is the SQL environment.
type DB struct {
	// Name is a database name to connect.
	Name string `validate:"required"`
	// User is an admin username.
	User string `validate:"required"`
	// Password is a password to the DBUser.
	Password string `validate:"required"`
	// InstanceConnectionName connects to the cloud SQL instance.
	InstanceConnectionName string `validate:"required"`
}
