package log

var defaultLogger Log = newStdGoLog()

type Log interface {
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool

	Error(description string, keywords string, format string, args ...interface{})
	Info(description string, keywords string, format string, args ...interface{})
	Debug(description string, keywords string, format string, args ...interface{})
	Warn(description string, keywords string, format string, args ...interface{})

	//when system exit you should use it
	Flush()

	GetLogLevel() string
	SetLogLevel(levelStr string)
}

func SetDefaultLogger(l Log) {
	defaultLogger = l
}

func IsDebugEnabled() bool {
	return defaultLogger.IsDebugEnabled()
}

func IsInfoEnabled() bool {
	return defaultLogger.IsInfoEnabled()
}

func IsWarnEnabled() bool {
	return defaultLogger.IsWarnEnabled()
}

func Error(description string, keywords string, format string, args ...interface{}) {
	defaultLogger.Error(description, keywords, format, args...)
}

func Info(description string, keywords string, format string, args ...interface{}) {
	defaultLogger.Info(description, keywords, format, args...)
}

func Debug(description string, keywords string, format string, args ...interface{}) {
	defaultLogger.Debug(description, keywords, format, args...)
}

func Warn(description string, keywords string, format string, args ...interface{}) {
	defaultLogger.Warn(description, keywords, format, args...)
}

func Flush() {
	defaultLogger.Flush()
}

func SetLogLevel(level string) {
	defaultLogger.SetLogLevel(level)
}

func GetLogLevel() string {
	return defaultLogger.GetLogLevel()
}
