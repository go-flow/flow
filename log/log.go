package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger -
type Logger interface {
	With(fields ...zap.Field) Logger

	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type zapWrapLogger struct {
	*zap.Logger
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l zapWrapLogger) With(fields ...zap.Field) Logger {
	return zapWrapLogger{l.Logger.With(fields...)}
}

// New creates new Logger instance for given logLevel nad mode
//
// logLevel can be: `debug`, `info`, `warn`, `error`, `dpanic`, `panic`, `fatal`
// mode can be `production` or `development`
func New(logLevel, mode string) Logger {
	if mode == "development" {
		return NewDevelopment(logLevel)
	}
	return NewProduction(logLevel)
}

// NewProduction logger instance
//
// logLevel can be: `debug`, `info`, `warn`, `error`, `dpanic`, `panic`, `fatal`
func NewProduction(logLevel string) Logger {
	cfg := zap.NewProductionConfig()
	return buildLogger(cfg, logLevel)
}

// NewDevelopment logger instance
//
// logLevel can be: `debug`, `info`, `warn`, `error`, `dpanic`, `panic`, `fatal`
func NewDevelopment(logLevel string) Logger {
	cfg := zap.NewDevelopmentConfig()
	return buildLogger(cfg, logLevel)
}

// buuldLogger instance for given logLevel and zap Configuration
func buildLogger(cfg zap.Config, logLevel string) Logger {
	lvl := parseLevel(logLevel)
	cfg.Level.SetLevel(lvl)

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return zapWrapLogger{logger}
}

// parseLevel parses string level to zap.Level object
func parseLevel(logLevel string) zapcore.Level {
	var lvl zapcore.Level
	err := lvl.UnmarshalText([]byte(logLevel))
	if err != nil {
		panic(err)
	}
	return lvl
}
