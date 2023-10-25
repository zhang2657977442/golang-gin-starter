package log

import (
	"fmt"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapConfig struct {
	Filename       string
	Level          string
	MaxFileSizeMB  int
	MaxRetainDays  int
	MaxRetainFiles int
	ModuleName     string
}

func (c *ZapConfig) Validate() error {
	if c.Filename == "" {
		return fmt.Errorf("invalid log filename: %v", c.Filename)
	}
	return nil
}

type zapLogger struct {
	atomLevel  zap.AtomicLevel
	logger     *zap.SugaredLogger
	moduleName string
}

func zapTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.00000"))
}

func createLogDir(dir string) error {
	if dir == "" {
		return nil
	}

	stat, err := os.Stat(dir)
	if (err != nil && !os.IsNotExist(err)) || (err == nil && !stat.IsDir()) {
		return fmt.Errorf("log dir[%s] already exists but not a dir", dir)
	}
	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create log dir[%s] failed: %v", dir, err)
		}
	}
	return nil
}

func NewZapLogger(c *ZapConfig) (Log, error) {
	// check config
	if err := c.Validate(); err != nil {
		return nil, err
	}
	// init log directory
	if err := createLogDir(path.Dir(c.Filename)); err != nil {
		return nil, err
	}

	// parse level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(c.Level)); err != nil {
		return nil, fmt.Errorf("invalid zap logger level: %v", err)
	}
	atomLevel := zap.NewAtomicLevelAt(level)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapTimeEncoder
	encoderCfg.CallerKey = "caller"
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	rotateWriter := &lumberjack.Logger{
		Filename:   c.Filename,
		MaxAge:     c.MaxRetainDays,
		MaxSize:    c.MaxFileSizeMB,
		MaxBackups: c.MaxRetainFiles,
		LocalTime:  true,
		Compress:   false,
	}

	log := &zapLogger{
		atomLevel: atomLevel,
		logger: zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(rotateWriter),
			atomLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(2)).Sugar(),
		moduleName: c.ModuleName,
	}
	return log, nil
}

func (l *zapLogger) IsDebugEnabled() bool {
	return l.atomLevel.Enabled(zap.DebugLevel)
}

func (l *zapLogger) IsInfoEnabled() bool {
	return l.atomLevel.Enabled(zap.InfoLevel)
}

func (l *zapLogger) IsWarnEnabled() bool {
	return l.atomLevel.Enabled(zap.WarnLevel)
}

func (l *zapLogger) Error(description string, keywords string, format string, args ...interface{}) {
	logEntry := fmt.Sprintf("[desc:%s] [keywords:%s] [%s]", description, keywords, format)
	l.logger.Errorf(logEntry, args...)
}

func (l *zapLogger) Info(description string, keywords string, format string, args ...interface{}) {
	logEntry := fmt.Sprintf("[desc:%s] [keywords:%s] [%s]", description, keywords, format)
	l.logger.Infof(logEntry, args...)
}

func (l *zapLogger) Debug(description string, keywords string, format string, args ...interface{}) {
	logEntry := fmt.Sprintf("[desc:%s] [keywords:%s] [%s]", description, keywords, format)
	l.logger.Debugf(logEntry, args...)
}

func (l *zapLogger) Warn(description string, keywords string, format string, args ...interface{}) {
	logEntry := fmt.Sprintf("[desc:%s] [keywords:%s] [%s]", description, keywords, format)
	l.logger.Warnf(logEntry, args...)
}

func (l *zapLogger) Flush() {
	_ = l.logger.Sync()
}

func (l *zapLogger) SetLogLevel(levelStr string) {
	_ = l.atomLevel.UnmarshalText([]byte(levelStr))
}

func (l *zapLogger) GetLogLevel() string {
	text, err := l.atomLevel.MarshalText()
	if err != nil {
		return "ERR:" + err.Error()
	} else {
		return string(text)
	}
}
