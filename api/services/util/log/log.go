package log

import (
	"api/services/util/logger"
	"os"
)

func init() {
	err := logger.SetLogger("./config/log.json")
	if err != nil {
		dir, err := os.Getwd()
		println(dir, err)
	}
}

func Debug(format string, v ...interface{}) {

	logger.Debug(format, v)
}

func Info(format string, v ...interface{}) {
	logger.Info(format, v)
}

func Warning(format string, v ...interface{}) {
	logger.Warn(format, v)
}

func Trace(format string, v ...interface{}) {
	logger.Trace(format, v)
}

func Error(format string, v ...interface{}) {
	logger.Error(format, v)
}

func Fatal(format string, v ...interface{}) {
	logger.Alert(format, v)
}
