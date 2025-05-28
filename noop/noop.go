// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package noop provides a no-operation logger implementation that satisfies
// the logging interface but with minimal overhead. This is useful for testing,
// benchmarking, or when logging needs to be disabled without changing code structure.
package noop

import (
	"github.com/goxkit/configs"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"

	zapInstance "github.com/goxkit/logging/zap"
)

// Install initializes and returns a no-operation logger that still outputs to stdout
// but doesn't connect to any OpenTelemetry collector. This provides a lightweight
// logging solution when full observability isn't required.
//
// It creates a basic log provider in the configs object and initializes a standard
// Zap logger configured for local output.
//
// Parameters:
//   - cfgs: Application configurations to use and update with the logger provider
//
// Returns:
//   - A configured zap.Logger instance
//   - An error if logger initialization fails
func Install(cfgs *configs.Configs) (*zap.Logger, error) {
	provider := sdklog.NewLoggerProvider()
	cfgs.LoggerProvider = provider
	return zapInstance.NewStdoutZapLogger(cfgs)
}
