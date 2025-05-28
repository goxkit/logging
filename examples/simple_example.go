// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// This example demonstrates using the Goxkit logging package with ConfigsBuilder.
package main

import (
	"github.com/goxkit/configs_builder"
	"go.uber.org/zap"
)

func main() {
	// Initialize configs with ConfigsBuilder
	// The Otlp() call enables OpenTelemetry for logging, tracing and metrics
	cfgs, err := configsBuilder.NewConfigsBuilder().Otlp().Build()
	if err != nil {
		panic(err)
	}

	// The logger is automatically configured and available in cfgs.Logger
	// Use structured logging with context fields
	cfgs.Logger.Info("Application started",
		zap.String("version", "1.0.0"),
		zap.String("env", cfgs.AppConfigs.Environment.String()))

	// Use different log levels based on the information importance
	cfgs.Logger.Debug("Configuration loaded",
		zap.Bool("otlp_enabled", cfgs.OTLPConfigs.Enabled),
		zap.String("otlp_endpoint", cfgs.OTLPConfigs.Endpoint))

	// Log operational events
	cfgs.Logger.Info("Ready to process requests")

	// Example of error logging with context
	err = performOperation()
	if err != nil {
		cfgs.Logger.Error("Operation failed",
			zap.Error(err),
			zap.String("operation", "startup_check"))
	}

	// When the application is shutting down
	cfgs.Logger.Info("Application shutting down")
}

func performOperation() error {
	// This is just a placeholder
	return nil
}
