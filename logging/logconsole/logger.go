package logconsole

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/logging"
	"os"
	"time"
	"github.com/yamamushi/EscapingEden/edenutil"
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
		colorCode = SHWhite
	case logging.LogInfo:
		colorCode = SHPurple
	case logging.LogWarn:
		colorCode = SHYellow
	case logging.LogError:
		colorCode = SHRed
	case logging.LogFatal:
		colorCode = SHPurple
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
