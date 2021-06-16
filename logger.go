package bot

type Logger interface {
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// TODO: add a mock logger
