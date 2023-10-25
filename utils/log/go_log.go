package log

import (
	"fmt"
	golog "log"
	"os"
	"strings"
)

type Level int8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR

	debug = "[DEBUG] "
	info  = "[INFO] "
	warn  = "[WARN] "
	err   = "[ERROR] "
)

func newStdGoLog() Log {
	return &GoLog{Logger: golog.New(os.Stderr, "", golog.Lshortfile|golog.LstdFlags), L: DEBUG, LevelCode: debug}
}

type GoLog struct {
	*golog.Logger
	L         Level
	LevelCode string
}

func (gl *GoLog) GetLogLevel() string {
	switch gl.L {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "err"
	default:
		return "<invalid>"
	}
}

func (gl *GoLog) SetLogLevel(levelStr string) {
	level := strings.ToUpper(levelStr)

	switch level {
	case "TRACE":
		gl.L = DEBUG
		gl.LevelCode = debug
	case "DEBUG":
		gl.L = DEBUG
		gl.LevelCode = debug
	case "INFO":
		gl.L = INFO
		gl.LevelCode = info
	case "WARN":
		gl.L = WARN
		gl.LevelCode = warn
	case "ERROR":
		gl.L = ERROR
		gl.LevelCode = err
	default:
		gl.L = INFO
		gl.LevelCode = info
	}
}

func (gl *GoLog) Flush() {
}

func NewGoLog(lg *golog.Logger, l Level) Log {
	return &GoLog{Logger: lg, L: l}
}

func (gl *GoLog) IsDebugEnabled() bool {
	return gl.L == 0
}

func (gl *GoLog) IsInfoEnabled() bool {
	return gl.L <= 1
}

func (gl *GoLog) IsWarnEnabled() bool {
	return gl.L <= 2
}

func (gl *GoLog) Debug(description string, keywords string, format string, args ...interface{}) {
	if gl.IsDebugEnabled() {
		gl.LevelCode = debug
		gl.write(description, keywords, format, args...)
	}
}

func (gl *GoLog) Info(description string, keywords string, format string, args ...interface{}) {
	if gl.IsInfoEnabled() {
		gl.LevelCode = info
		gl.write(description, keywords, format, args...)
	}
}

func (gl *GoLog) Warn(description string, keywords string, format string, args ...interface{}) {
	if gl.IsWarnEnabled() {
		gl.LevelCode = warn
		gl.write(description, keywords, format, args...)
	}
}

func (gl *GoLog) Error(description string, keywords string, format string, args ...interface{}) {
	gl.LevelCode = err
	gl.write(description, keywords, format, args...)
}

func (gl *GoLog) write(description string, keywords string, format string, args ...interface{}) {
	logEntry := fmt.Sprintf("[desc:%s] [keywords:%s] [%s]", description, keywords, format)
	if len(args) == 0 {
		_ = gl.Output(4, gl.LevelCode+logEntry)
	} else {
		_ = gl.Output(4, fmt.Sprintf(gl.LevelCode+logEntry, args...))
	}

}
