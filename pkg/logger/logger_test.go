package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockWriter is a mock io.Writer for testing
type mockWriter struct {
	buffer bytes.Buffer
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *mockWriter) String() string {
	return m.buffer.String()
}

func (m *mockWriter) Reset() {
	m.buffer.Reset()
}

func TestLoggerLevels(t *testing.T) {
	t.Run("should respect log levels", func(t *testing.T) {
		// arrange
		writer := &mockWriter{}
		logger := &LoggerImpl{
			level:  INFO,
			logger: log.New(writer, "", 0),
			fields: make(Fields),
		}

		// act & assert
		// Debug should not be logged at INFO level
		logger.Debug("This is a debug message")
		assert.Empty(t, writer.String())
		writer.Reset()

		// Info should be logged at INFO level
		logger.Info("This is an info message")
		assert.Contains(t, writer.String(), "INFO")
		assert.Contains(t, writer.String(), "This is an info message")
		writer.Reset()

		// Warn should be logged at INFO level
		logger.Warn("This is a warning message")
		assert.Contains(t, writer.String(), "WARN")
		assert.Contains(t, writer.String(), "This is a warning message")
		writer.Reset()

		// Error should be logged at INFO level
		logger.Error("This is an error message")
		assert.Contains(t, writer.String(), "ERROR")
		assert.Contains(t, writer.String(), "This is an error message")
		writer.Reset()
	})
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Level
	}{
		{"should parse debug level", "debug", DEBUG},
		{"should parse info level", "info", INFO},
		{"should parse warn level", "warn", WARN},
		{"should parse error level", "error", ERROR},
		{"should parse fatal level", "fatal", FATAL},
		{"should default to info for unknown level", "unknown", INFO},
		{"should default to info for empty level", "", INFO},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange & act
			got := parseLevel(tt.input)

			// assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Run("should create a new logger with the specified level", func(t *testing.T) {
		// arrange & act
		logger := New("debug")

		// assert
		assert.NotNil(t, logger)
		loggerImpl, ok := logger.(*LoggerImpl)
		assert.True(t, ok, "logger should be a *LoggerImpl")
		assert.Equal(t, DEBUG, loggerImpl.level)
	})

	t.Run("should create a new logger with default level when level is invalid", func(t *testing.T) {
		// arrange & act
		logger := New("invalid")

		// assert
		assert.NotNil(t, logger)
		loggerImpl, ok := logger.(*LoggerImpl)
		assert.True(t, ok, "logger should be a *LoggerImpl")
		assert.Equal(t, INFO, loggerImpl.level)
	})
}

func TestLogFormatting(t *testing.T) {
	t.Run("should format log messages correctly", func(t *testing.T) {
		// arrange
		writer := &mockWriter{}
		logger := &LoggerImpl{
			level:  DEBUG,
			logger: log.New(writer, "", 0),
			fields: make(Fields),
		}

		// act
		logger.Info("Test message with parameter", Fields{
			"param": "value",
		})

		// assert
		logged := writer.String()
		assert.Contains(t, logged, "[INFO]")
		assert.Contains(t, logged, "Test message with parameter")
		assert.Contains(t, logged, "param: value")

		// Verify date format (YYYY-MM-DD)
		datePart := strings.Split(logged, " ")[1]
		assert.Regexp(t, `\[\d{4}-\d{2}-\d{2}`, datePart)
	})
}

func TestWithFields(t *testing.T) {
	t.Run("should add fields to logger", func(t *testing.T) {
		// arrange
		writer := &mockWriter{}
		logger := &LoggerImpl{
			level:  DEBUG,
			logger: log.New(writer, "", 0),
			fields: make(Fields),
		}

		// act
		loggerWithFields := logger.With(Fields{
			"service": "test",
		})
		
		loggerWithFields.Info("Message with fields")

		// assert
		logged := writer.String()
		assert.Contains(t, logged, "Message with fields")
		assert.Contains(t, logged, "service: test")
	})
	
	t.Run("should merge fields when logging with additional fields", func(t *testing.T) {
		// arrange
		writer := &mockWriter{}
		logger := &LoggerImpl{
			level:  DEBUG,
			logger: log.New(writer, "", 0),
			fields: make(Fields),
		}

		// act
		loggerWithFields := logger.With(Fields{
			"service": "test",
		})
		
		loggerWithFields.Info("Message with merged fields", Fields{
			"request_id": "123",
		})

		// assert
		logged := writer.String()
		assert.Contains(t, logged, "Message with merged fields")
		assert.Contains(t, logged, "service: test")
		assert.Contains(t, logged, "request_id: 123")
	})
}
