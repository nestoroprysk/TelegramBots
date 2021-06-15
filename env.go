package bot

import (
	"fmt"
	"os"
)

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

// MakeEnv initializes the environment from environmental variables.
func MakeEnv() (Env, error) {
	result := Env{}
	for res, key := range map[*string]string{
		&result.Telegram.Token:            "TOKEN",
		&result.DB.Name:                   "DB_NAME",
		&result.DB.User:                   "DB_USER",
		&result.DB.Password:               "DB_PASS",
		&result.DB.InstanceConnectionName: "INSTANCE_CONNECTION_NAME",
	} {
		val, ok := os.LookupEnv(key)
		if !ok || val == "" {
			// TODO: capture
			return Env{}, fmt.Errorf("environmental variable %q should be set to a non-empty value", key)
		}

		*res = val
	}

	return result, nil
}
