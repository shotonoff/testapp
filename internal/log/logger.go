package log

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

type (
	ctxKey int
	Logger interface {
		Info(message string, args ...any)
		Warn(message string, args ...any)
		Error(message string, args ...any)
		Debug(message string, args ...any)
		Trace(message string, args ...any)
		With(args ...any) Logger
	}
	Config struct {
		Level string
	}
	OptionFunc func(cfg *Config)
)

const (
	loggerCtxKey ctxKey = iota

	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	DebugLevel = "debug"
	TraceLevel = "trace"
)

// WithLevel is an option function that sets the log level
func WithLevel(level string) OptionFunc {
	return func(cfg *Config) {
		cfg.Level = level
	}
}

type logrusLogger struct {
	fields logrus.Fields
	logger *logrus.Logger
}

// New creates a new logger
func New(opts ...OptionFunc) Logger {
	cfg := Config{
		Level: DebugLevel,
	}
	l := &logrusLogger{
		fields: logrus.Fields{},
		logger: logrus.New(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Level != "" {
		var err error
		l.logger.Level, err = logrus.ParseLevel(cfg.Level)
		if err != nil {
			panic(err)
		}
	}
	return l
}

// NewNop creates a new logger that discards all logs
func NewNop() Logger {
	logger := logrus.New()
	logger.Out = io.Discard
	return &logrusLogger{
		logger: logger,
	}
}

// With returns a new logger with the given fields
func (l *logrusLogger) With(args ...any) Logger {
	return &logrusLogger{
		fields: extendFields(l.fields, args...),
		logger: l.logger,
	}
}

func (l *logrusLogger) Log(level logrus.Level, msg string, args ...any) {
	fields := extendFields(l.fields, args...)
	l.logger.WithFields(fields).Log(level, msg)
}

func (l *logrusLogger) Info(message string, args ...any) {
	l.Log(logrus.InfoLevel, message, args...)
}

func (l *logrusLogger) Warn(message string, args ...any) {
	l.Log(logrus.WarnLevel, message, args...)
}

func (l *logrusLogger) Debug(message string, args ...any) {
	l.Log(logrus.DebugLevel, message, args...)
}

func (l *logrusLogger) Trace(message string, args ...any) {
	l.Log(logrus.TraceLevel, message, args...)
}

func (l *logrusLogger) Error(message string, args ...any) {
	l.Log(logrus.ErrorLevel, message, args...)
}

// WithContext returns a new context with the given logger
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func extendFields(origin logrus.Fields, args ...any) logrus.Fields {
	fields := make(logrus.Fields)
	for k, v := range origin {
		fields[k] = v
	}
	for i := 0; i < len(args); i += 2 {
		fields[args[i].(string)] = args[i+1]
	}
	return fields
}
