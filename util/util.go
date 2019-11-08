package util

import "os"

var (
	LogLevel       = os.Getenv("LOG_LEVEL")
	LogDirPath     = ".logs"
	LogFilePath    = "app.log"
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
