package logger

import (
	"github.com/labstack/gommon/log"
)

// Logger singleton logger
// make sure InitializeLogger must be called at startup
var Logger *log.Logger

func InitializeLogger() {
	Logger = log.New("")
	Logger.SetHeader("${level} | ${time_rfc3339} | ${short_file}:${line} | ${message}")
}
