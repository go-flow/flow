package flow

import (
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus FieldLogger
type Logger struct {
	logrus.FieldLogger
}

// WithField Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (l Logger) WithField(s string, i interface{}) Logger {
	return Logger{l.FieldLogger.WithField(s, i)}
}

// WithFields Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (l Logger) WithFields(m map[string]interface{}) Logger {
	return Logger{l.FieldLogger.WithFields(m)}
}

// NewLogger based on the specified log level.
// This logger will log to the STDOUT in a human readable,
// but parseable form.
/*
	Example: time="2016-12-01T21:02:07-05:00" level=info duration=225.283µs human_size="106 B" method=GET path="/" render=199.79µs request_id=2265736089 size=106 status=200
*/
func NewLogger(level string) *Logger {
	l := logrus.New()
	l.Level, _ = logrus.ParseLevel(level)
	l.Formatter = &logrus.JSONFormatter{}

	return &Logger{l}
}
