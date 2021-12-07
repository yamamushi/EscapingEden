package logfile

import (
	"fmt"
	"github.com/yamamushi/EscapingEden/edenutil"
	"github.com/yamamushi/EscapingEden/logging"
	"os"
	"strings"
	"time"
)

type FileLogger struct {
	logging.Logger
	FileName string
}

func NewFileLogger(path string) (*FileLogger, error) {

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// create path if not exists
	err := edenutil.CreatePathIfNotExists(path)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("01-02-2006-150405-")
	filename := path + timestamp + "server.log"
	// create file if not exists
	err = edenutil.CreateFileIfNotExists(filename)

	logger := &FileLogger{FileName: filename}
	logger.TypeID = logging.LoggerTypeID_Console
	return logger, nil
}

func (fl *FileLogger) Println(level logging.LogLevel, message string, v ...interface{}) {
	fl.LogMutex.Lock()
	defer fl.LogMutex.Unlock()
	// We're dealing with open files here and not taking any chances

	err := edenutil.CreateFileIfNotExists(fl.FileName)
	if err != nil {
		fmt.Println("Could not create log file: ", err)
	}

	file, err := os.OpenFile(fl.FileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Could not open log file: ", err)
	}
	defer file.Close()

	timestamp := time.Now().Format("2006/01/02 15:04:05")

	// Write to file
	if v == nil {
		_, err = fmt.Fprintf(file, "[%s] %s: %s\n", timestamp, level.String(), message)
		if err != nil {
			fmt.Println("Could not write to log file: ", err)
		}
	} else {
		_, err = fmt.Fprintf(file, "[%s] %s: %s\n", timestamp, level.String(), fmt.Sprintf(message, v...))
		if err != nil {
			fmt.Println("Could not write to log file: ", err)
		}
	}

	if level == logging.LogFatal {
		os.Exit(1)
	}
}
