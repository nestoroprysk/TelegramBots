package env

import "fmt"

// Env contains information from the global variables.
type Env struct {
	Telegram
	DB
}

// Telegram is the Telegram environment.
type Telegram struct {
	// Token is a Telegram bot token.
	Token string
}

// DB is the SQL environment.
type DB struct {
	// Name is a database name to connect.
	Name string
	// User is an admin username.
	User string
	// Password is a password to the DBUser.
	Password string
	// InstanceConnectionName connects to the cloud SQL instance.
	InstanceConnectionName string
}

func (e Env) IsValid() error {
	for _, field := range []string{
		e.Telegram.Token,
		e.DB.Name,
		e.DB.User,
		e.DB.Password,
		e.DB.InstanceConnectionName,
	} {
		if field == "" {
			return fmt.Errorf("%q should be set, yet it's empty", field)
		}
	}

	return nil
}
