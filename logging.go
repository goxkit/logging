// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package logging provides structured logging capabilities powered by Zap
// with various configuration options for different environments.
package logging

import (
	"github.com/goxkit/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/goxkit/logging/otlp"
)

type (
	// Logger defines the interface for logging operations within the application.
	// It provides methods for different log levels and the ability to add context fields.
	Logger interface {
		// With adds structured context to the logger.
		With(fields ...zapcore.Field) *zap.Logger

		// Debug logs a message at Debug level with optional fields.
		Debug(msg string, fields ...zap.Field)

		// Info logs a message at Info level with optional fields.
		Info(msg string, fields ...zap.Field)

		// Warn logs a message at Warn level with optional fields.
		Warn(msg string, fields ...zap.Field)

		// Error logs a message at Error level with optional fields.
		Error(msg string, fields ...zap.Field)

		// Fatal logs a message at Fatal level with optional fields,
		// then calls os.Exit(1).
		Fatal(msg string, fields ...zap.Field)
	}
)

// NewDefaultLogger creates a new logger that outputs to stdout.
// It configures the logger based on the environment:
// - Production/Staging: Uses JSON encoder
// - Development: Uses colored console output
//
// The log level is determined by the configuration provided.
func NewLogger(cfgs *configs.Configs) (Logger, error) {
	if cfgs.OTLPConfigs.Enabled {
		return otlp.Install(cfgs)
	}

	return nil, nil
}
