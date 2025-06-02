package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Fields represents a map of field names to values
type Fields map[string]interface{}

// Logger interface defines methods for logging at different levels
type Logger interface {
	Debug(msg string, fields ...Fields)
	Info(msg string, fields ...Fields)
	Warn(msg string, fields ...Fields)
	Error(msg string, fields ...Fields)
	Fatal(msg string, fields ...Fields)
	With(fields Fields) Logger
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
	fields Fields
}

// New creates a new logger with the specified level
func New(level string) Logger {
	l := &LoggerImpl{
		level:  parseLevel(level),
		logger: log.New(os.Stdout, "", 0),
		fields: make(Fields),
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

// formatFields formats fields as a string
func formatFields(fields Fields) string {
	if len(fields) == 0 {
		return ""
	}

	result := "{"
	first := true
	for k, v := range fields {
		if !first {
			result += ", "
		}
		first = false
		result += fmt.Sprintf("%s: %v", k, v)
	}
	result += "}"
	return result
}

// mergeFields merges multiple Fields objects
func mergeFields(fieldsList ...Fields) Fields {
	result := make(Fields)
	
	// Add base fields
	for k, v := range fieldsList[0] {
		result[k] = v
	}
	
	// Add additional fields
	if len(fieldsList) > 1 {
		for _, fields := range fieldsList[1:] {
			for k, v := range fields {
				result[k] = v
			}
		}
	}
	
	return result
}

// log logs a message with the specified level
func (l *LoggerImpl) log(level Level, msg string, fields ...Fields) {
	if level < l.level {
		return
	}

	// Merge instance fields with the provided fields
	mergedFields := l.fields
	if len(fields) > 0 {
		mergedFields = mergeFields(l.fields, fields[0])
	}

	fieldsStr := formatFields(mergedFields)
	prefix := fmt.Sprintf("[%s] [%s] ", levelNames[level], time.Now().Format("2006-01-02 15:04:05"))
	if fieldsStr != "" {
		l.logger.Printf("%s%s %s", prefix, msg, fieldsStr)
	} else {
		l.logger.Printf("%s%s", prefix, msg)
	}
}

// With returns a new logger with the specified fields added
func (l *LoggerImpl) With(fields Fields) Logger {
	newLogger := &LoggerImpl{
		level:  l.level,
		logger: l.logger,
		fields: mergeFields(l.fields, fields),
	}
	return newLogger
}

// Debug logs a debug message
func (l *LoggerImpl) Debug(msg string, fields ...Fields) {
	l.log(DEBUG, msg, fields...)
}

// Info logs an info message
func (l *LoggerImpl) Info(msg string, fields ...Fields) {
	l.log(INFO, msg, fields...)
}

// Warn logs a warning message
func (l *LoggerImpl) Warn(msg string, fields ...Fields) {
	l.log(WARN, msg, fields...)
}

// Error logs an error message
func (l *LoggerImpl) Error(msg string, fields ...Fields) {
	l.log(ERROR, msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *LoggerImpl) Fatal(msg string, fields ...Fields) {
	l.log(FATAL, msg, fields...)
	os.Exit(1)
}
