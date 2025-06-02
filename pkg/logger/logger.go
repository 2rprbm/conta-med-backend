package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger interface defines methods for logging at different levels
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

// Level represents the logging level
type Level int

const (
	// DEBUG level
	DEBUG Level = iota
	// INFO level
	INFO
	// WARN level
	WARN
	// ERROR level
	ERROR
	// FATAL level
	FATAL
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// LoggerImpl represents a logger instance
type LoggerImpl struct {
	level  Level
	logger *log.Logger
}

// New creates a new logger with the specified level
func New(level string) Logger {
	l := &LoggerImpl{
		level:  parseLevel(level),
		logger: log.New(os.Stdout, "", 0),
	}
	return l
}

// parseLevel parses a string level to Level
func parseLevel(level string) Level {
	switch level {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// log logs a message with the specified level
func (l *LoggerImpl) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	prefix := fmt.Sprintf("[%s] [%s] ", levelNames[level], time.Now().Format("2006-01-02 15:04:05"))
	l.logger.Printf(prefix+format, args...)
}

// Debug logs a debug message
func (l *LoggerImpl) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *LoggerImpl) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *LoggerImpl) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *LoggerImpl) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *LoggerImpl) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}
