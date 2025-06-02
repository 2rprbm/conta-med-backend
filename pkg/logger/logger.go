package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

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

// Logger represents a logger instance
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new logger with the specified level
func New(level string) *Logger {
	l := &Logger{
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
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	prefix := fmt.Sprintf("[%s] [%s] ", levelNames[level], time.Now().Format("2006-01-02 15:04:05"))
	l.logger.Printf(prefix+format, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}
