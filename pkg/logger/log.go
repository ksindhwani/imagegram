package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFormatJSON    = "json"
	logFormatText    = "text"
	logFormatConsole = "console"
)

const (
	LogLevelError = "error"
	LogLevelWarn  = "warn"
	LogLevelFatal = "fatal"
	LogLevelPanic = "panic"
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
)

// These are the names of the log fields in accordance with Elastic-APM:
// https://www.elastic.co/blog/how-to-instrument-your-go-app-with-the-elastic-apm-go-agent#logs
const (
	logKeyMessage  = "message"
	logKeyTime     = "@timestamp"
	logKeyRevision = "revision"
)

// Logger is a small wrapper around a zap.SugaredLogger.
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new Logger with given logLevel, logFormat and revision as part
// of a permanent field of the logger.
func New(logLevel, logFormat, revision string) (*Logger, error) {
	if logFormat == logFormatText {
		logFormat = logFormatConsole
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Encoding = logFormat

	var level zapcore.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, err
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	zapConfig.EncoderConfig.MessageKey = logKeyMessage
	zapConfig.EncoderConfig.TimeKey = logKeyTime
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapConfig.InitialFields = map[string]interface{}{
		logKeyRevision: revision,
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("could not build logger: %w", err)
	}

	zap.ReplaceGlobals(logger)

	return &Logger{SugaredLogger: logger.Sugar()}, nil
}

// Printf prints the given message on info level.s
func (l *Logger) Infof(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

// Printf prints the given message on info level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}

// Printf prints the given message on info level.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.SugaredLogger.Fatalf(format, args...)
}
