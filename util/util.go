package util

import "os"

var (
	LogLevel       = os.Getenv("LOG_LEVEL")
	LogType        = os.Getenv("LOG_TYPE")
	LogTypeFile    = "file"
	LogTypeConsole = "console"
	LogTypeOff     = "off"
	LogTypeAll     = "all"
	LogRoot        = os.Getenv("LOG_ROOT")
	LogDirPath     = ".logs"
	LogFilePath    = os.Getenv("LOG_FILE")
	LinkFilePath   = "link.log"
	LogLevelOff    = "off"
	LogLevelDebug  = "debug"
	LogLevelInfo   = "info"
	LogLevelWarn   = "warn"
	LogLevelError  = "error"
	LogLevelFatal  = "fatal"
	ServerType     = os.Getenv("SERVER_TYPE")
	ServerTypeProd = "production"
)

func init() {
	if LogFilePath == "" {
		LogFilePath = "app.log"
	}
}
