package handlers

import (
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...logger.Fields)  {}
func (m *mockLogger) Error(msg string, fields ...logger.Fields) {}
func (m *mockLogger) Debug(msg string, fields ...logger.Fields) {}
func (m *mockLogger) Warn(msg string, fields ...logger.Fields)  {}
func (m *mockLogger) Fatal(msg string, fields ...logger.Fields) {}
func (m *mockLogger) With(fields logger.Fields) logger.Logger   { return m } 