package flow

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus FieldLogger
type Logger interface {
	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Debug(...interface{})
	Info(...interface{})
	Print(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
}

var _ Logger = defaultLogger{}

type defaultLogger struct {
	logrus.FieldLogger
}

// WithField Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (l defaultLogger) WithField(s string, i interface{}) Logger {
	return defaultLogger{l.FieldLogger.WithField(s, i)}
}

// WithFields Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (l defaultLogger) WithFields(m map[string]interface{}) Logger {
	return defaultLogger{l.FieldLogger.WithFields(m)}
}

// NewLogger based on the specified log level.
// This logger will log to the STDOUT in a human readable,
// but parseable form.
func NewLogger(level string) Logger {
	l := logrus.New()
	l.Level, _ = logrus.ParseLevel(level)
	l.Formatter = &logrus.JSONFormatter{}

	return defaultLogger{l}
}

// NewLoggerWithFormatter  creates logger instance
// based on the specified log level and formatter
//
//This logger will log to the STDOUT in a human readable,
// but parseable form.
// Supported formatters are `text` and `json`
func NewLoggerWithFormatter(level, formatter string) Logger {
	l := logrus.New()
	l.Level, _ = logrus.ParseLevel(level)
	switch strings.ToLower(formatter) {
	case "text":
		l.Formatter = &logrus.TextFormatter{}
	case "json":
		l.Formatter = &logrus.JSONFormatter{}
	default:
		panic(fmt.Sprintf("Unsupported Logger formatter `%s`. Supported formaters are `text` and `json` ", formatter))
	}

	return defaultLogger{l}
}
