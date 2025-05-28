// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package logging provides structured logging capabilities powered by Zap.
// This file contains mock implementations for testing purposes.
package logging

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLogger is a mock implementation of the Logger interface
// that can be used in unit tests to verify logging behavior
// without actually producing log output. It allows testing code
// that uses the logger without creating actual logs.
//
// It leverages the testify/mock package to provide mocking capabilities
// for assertion and verification in tests.
type MockLogger struct {
	mock.Mock
}

// With implements the Logger interface's With method for the mock.
// In a real logger, this would add structured context fields.
// For testing purposes, it returns nil instead of a real logger.
//
// Parameters:
//   - fields: The zap fields that would be added to a real logger
//
// Returns:
//   - nil (since it's a mock)
func (m *MockLogger) With(_ ...zap.Field) *zap.Logger {
	return nil
}

// Debug implements the Logger interface's Debug method for the mock.
// In test scenarios, this can be used to verify Debug level logs were attempted.
//
// Parameters:
//   - msg: The message that would be logged
//   - fields: The zap fields that would be included in the log
func (m *MockLogger) Debug(_ string, _ ...zap.Field) {
}

// Info implements the Logger interface's Info method for the mock.
// In test scenarios, this can be used to verify Info level logs were attempted.
//
// Parameters:
//   - msg: The message that would be logged
//   - fields: The zap fields that would be included in the log
func (m *MockLogger) Info(_ string, _ ...zap.Field) {
}

// Warn implements the Logger interface's Warn method for the mock.
// In test scenarios, this can be used to verify Warn level logs were attempted.
//
// Parameters:
//   - msg: The message that would be logged
//   - fields: The zap fields that would be included in the log
func (m *MockLogger) Warn(_ string, _ ...zap.Field) {
}

// Error implements the Logger interface's Error method for the mock.
// In test scenarios, this can be used to verify Error level logs were attempted.
//
// Parameters:
//   - msg: The message that would be logged
//   - fields: The zap fields that would be included in the log
func (m *MockLogger) Error(_ string, _ ...zap.Field) {
}

// Fatal implements the Logger interface's Fatal method for the mock.
// In test scenarios, this can be used to verify Fatal level logs were attempted.
// Unlike a real logger, this won't terminate the application.
//
// Parameters:
//   - msg: The message that would be logged
//   - fields: The zap fields that would be included in the log
func (m *MockLogger) Fatal(_ string, _ ...zap.Field) {
}

// NewMockLogger creates and returns a new instance of MockLogger
// that can be used in tests to verify logging behavior without
// producing actual log output.
//
// Returns:
//   - A configured MockLogger instance ready for use in tests
func NewMockLogger() *MockLogger {
	return new(MockLogger)
}
