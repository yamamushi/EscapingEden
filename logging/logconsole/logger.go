package logconsole

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/logging"
	"os"
	"time"
)

type ConsoleLogger struct {
	logging.Logger
}

func NewConsoleLogger() *ConsoleLogger {
	logger := &ConsoleLogger{}
	logger.TypeID = logging.LoggerTypeID_Console
	return logger
}

func (cl *ConsoleLogger) Println(level logging.LogLevel, message string, v ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	var colorCode string
	switch level {
	case logging.LogDebug:
		colorCode = "\033[37m"
	case logging.LogInfo:
		colorCode = "\033[36m"
	case logging.LogWarn:
		colorCode = "\033[33m"
	case logging.LogError:
		colorCode = "\033[31m"
	case logging.LogFatal:
		colorCode = "\033[35m"
	}

	if v == nil {
		fmt.Println(colorCode, timestamp, level, message, "\033[0m")
	} else {
		fmt.Println(colorCode, timestamp, level, message, v, "\033[0m")
	}

	if level == logging.LogFatal {
		os.Exit(1)
	}
}
