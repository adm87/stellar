package logging

import "log/slog"

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

func (l LogLevel) IsValid() bool {
	switch l {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel:
		return true
	default:
		return false
	}
}

type Logger struct {
	slogger *slog.Logger
}

func NewLogger(args ...any) *Logger {
	return &Logger{
		slogger: slog.Default().With(args...),
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.slogger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.slogger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.slogger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.slogger.Error(msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		slogger: l.slogger.With(args...),
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	switch level {
	case DebugLevel:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case InfoLevel:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case WarnLevel:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case ErrorLevel:
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}
