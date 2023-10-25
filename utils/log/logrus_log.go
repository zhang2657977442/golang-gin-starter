package log

import (
	"fmt"
	"os"
	"path"
	"runtime"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogrusConfig struct {
	Filename       string
	Level          string
	MaxFileSizeMB  int
	MaxRetainDays  int
	MaxRetainFiles int
	ModuleName     string
}

func (c *LogrusConfig) Validate() error {
	if c.Filename == "" {
		return fmt.Errorf("invalid log filename: %v", c.Filename)
	}
	return nil
}

type logrusLogger struct {
	logger     *logrus.Logger
	moduleName string
}

func NewLogrusLogger(c *LogrusConfig) (Log, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := createLogDir(path.Dir(c.Filename)); err != nil {
		return nil, err
	}

	// parse level
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid zap logger level: %v", err)
	}

	// set rotate Writer
	rotateWriter := &lumberjack.Logger{
		Filename:   c.Filename,
		MaxAge:     c.MaxRetainDays,
		MaxSize:    c.MaxFileSizeMB,
		MaxBackups: c.MaxRetainFiles,
		LocalTime:  true,
		Compress:   false,
	}

	formatter := &logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02T15:04:05.000Z07:00", // 时间戳格式
		DisableTimestamp:  false,                           // 启用时间戳输出
		DisableHTMLEscape: true,                            // 禁用 HTML 转义
		DataKey:           "",                              // 不将日志参数放入嵌套字典
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp", // 时间戳字段名
			logrus.FieldKeyLevel: "level",     // 日志级别字段名
			logrus.FieldKeyMsg:   "message",   // 消息字段名
			logrus.FieldKeyFunc:  "caller",    // 调用者信息字段名
			logrus.FieldKeyFile:  "file",      // 文件名字段名
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", "" // 不修改调用者信息的内容，即不添加函数名和文件名
		},
		PrettyPrint: false, // 不缩进 JSON 输出
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetReportCaller(true)
	logger.SetFormatter(formatter)
	logger.SetOutput(rotateWriter)
	logger.SetOutput(os.Stdout)

	log := &logrusLogger{
		logger: logger,
	}

	return log, nil
}

func (l *logrusLogger) IsDebugEnabled() bool {
	return l.logger.IsLevelEnabled(logrus.DebugLevel)
}

func (l *logrusLogger) IsInfoEnabled() bool {
	return l.logger.IsLevelEnabled(logrus.InfoLevel)
}

func (l *logrusLogger) IsWarnEnabled() bool {
	return l.logger.IsLevelEnabled(logrus.WarnLevel)
}

func (l *logrusLogger) Error(description string, keywords string, format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"desc":     description,
		"keywords": keywords,
	}).Errorf(format, args...)
}

func (l *logrusLogger) Info(description string, keywords string, format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"desc":     description,
		"keywords": keywords,
	}).Infof(format, args...)

}

func (l *logrusLogger) Debug(description string, keywords string, format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"desc":     description,
		"keywords": keywords,
	}).Debugf(format, args...)
}

func (l *logrusLogger) Warn(description string, keywords string, format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"desc":     description,
		"keywords": keywords,
	}).Warnf(format, args...)
}

func (l *logrusLogger) Flush() {

}

func (l *logrusLogger) SetLogLevel(levelStr string) {
	level, _ := logrus.ParseLevel(levelStr)
	l.logger.SetLevel(level)
}

func (l *logrusLogger) GetLogLevel() string {
	text, err := l.logger.GetLevel().MarshalText()
	if err != nil {
		return "ERR:" + err.Error()
	} else {
		return string(text)
	}
}
