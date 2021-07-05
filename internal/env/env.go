package env

// Env contains information from the global variables.
type Env struct {
	Telegram `validate:"required"`
	DB       `validate:"required"`
}

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
