// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package logging provides a comprehensive structured logging framework powered by Zap
// with OpenTelemetry integration for distributed tracing. It offers configurable
// output modes including OTLP export for observability platforms, standard output
// for local development, and no-operation mode for testing.
//
// The package implements a builder pattern allowing you to easily configure
// logging based on your environment and needs. It integrates with the configs
// package to use environment-specific settings and log levels.
package logging

import (
	"github.com/goxkit/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/goxkit/logging/noop"
	"github.com/goxkit/logging/otlp"
)

type (
	// Logger defines the interface for logging operations within the application.
	// It wraps zap.Logger to provide standardized logging methods across the application
	// while allowing for structured context and different log levels.
	Logger interface {
		// With adds structured context to the logger.
		// Returns a new zap.Logger instance with the added fields.
		With(fields ...zapcore.Field) *zap.Logger

		// Debug logs a message at Debug level with optional structured fields.
		// Debug logs are typically used for verbose information useful during development.
		Debug(msg string, fields ...zap.Field)

		// Info logs a message at Info level with optional structured fields.
		// Info logs are used for general operational information about the application.
		Info(msg string, fields ...zap.Field)

		// Warn logs a message at Warn level with optional structured fields.
		// Warn logs indicate potential issues that don't prevent the application from working.
		Warn(msg string, fields ...zap.Field)

		// Error logs a message at Error level with optional structured fields.
		// Error logs indicate issues that may require attention but don't stop the application.
		Error(msg string, fields ...zap.Field)

		// Fatal logs a message at Fatal level with optional structured fields,
		// then calls os.Exit(1), terminating the application immediately.
		// Use Fatal sparingly, only for errors that truly require immediate shutdown.
		Fatal(msg string, fields ...zap.Field)
	}
)

// NewLogger creates a configured logger based on the provided configurations.
// If OTLP configurations are enabled in the provided configs, it will set up
// a logger that exports to an OpenTelemetry collector. Otherwise, it will
// create a no-operation logger that still provides the Logger interface
// but with minimal functionality.
//
// Parameters:
//   - cfgs: Application configurations including logging settings
//
// Returns:
//   - A configured Logger implementation
//   - An error if logger initialization fails
func NewLogger(cfgs *configs.Configs) (Logger, error) {
	if cfgs.OTLPConfigs.Enabled {
		return otlp.Install(cfgs)
	}

	return noop.Install(cfgs)
}
