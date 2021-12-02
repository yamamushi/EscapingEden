package logging

import "sync"

type LoggerTypeID int

const (
	LoggerTypeID_Unkonwn LoggerTypeID = iota
	LoggerTypeID_Console
	LoggerTypeID_File
)

type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
	LogFatal
)

func (ll LogLevel) String() string {
	switch ll {
	case LogDebug:
		return "[Debug]"
	case LogInfo:
		return "[Info]"
	case LogWarn:
		return "[Warn]"
	case LogError:
		return "[Error]"
	case LogFatal:
		return "[Fatal]"
	default:
		return "[Unknown]"
	}
}

type LoggerType interface {
	GetTypeID() LoggerTypeID
	Println(level LogLevel, format string, v ...interface{})
}

type Logger struct {
	TypeID LoggerTypeID

	LogMutex sync.Mutex
}

func (l *Logger) GetTypeID() LoggerTypeID {
	return l.TypeID
}

func (l *Logger) Println(level LogLevel, format string, v ...interface{}) {
	return
}
