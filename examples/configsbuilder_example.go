// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package examples provides usage examples for the Goxkit logging module.
package examples

import (
	"context"
	"os"
	"time"

	"github.com/goxkit/configs"
	configsBuilder "github.com/goxkit/configs_builder"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ConfigsBuilderBasicExample demonstrates how to set up logging with the ConfigsBuilder
// and use the resulting logger for basic operations.
func ConfigsBuilderBasicExample() {
	// Set up environment variables (optional - typically done in deployment)
	os.Setenv("APP_NAME", "ExampleService")
	os.Setenv("APP_NAMESPACE", "examples")
	os.Setenv("GO_ENV", "development") // development, staging, production
	os.Setenv("LOG_LEVEL", "debug")

	// Use the ConfigsBuilder to set up your application with OTLP enabled
	cfg, err := configsBuilder.NewConfigsBuilder().
		Otlp(). // Enable OpenTelemetry logging
		HTTP(). // Add HTTP configuration if needed
		Build()
	if err != nil {
		panic(err)
	}

	// The logger is automatically configured and available in cfg.Logger
	cfg.Logger.Info("Application initialized",
		zap.String("version", "1.0.0"),
		zap.Int("port", cfg.AppConfigs.Port),
	)

	// Use the logger in your application code
	cfg.Logger.Debug("Debug information", zap.Int("workers", 4))
}

// ConfigsBuilderTracingExample demonstrates how to use the logger with tracing
// for complete observability.
func ConfigsBuilderTracingExample() {
	// Set up environment variables for OTLP
	os.Setenv("APP_NAME", "TracedService")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")

	// Create configs with both OTLP and HTTP enabled
	cfg, err := configsBuilder.NewConfigsBuilder().
		Otlp().     // Enable OpenTelemetry for tracing, metrics, and logging
		HTTP().     // HTTP configuration
		Postgres(). // Database configuration
		Build()
	if err != nil {
		panic(err)
	}

	// Create a root context
	ctx := context.Background()

	// Start a trace span
	tracer := otel.GetTracerProvider().Tracer("example-tracer")
	ctx, span := tracer.Start(ctx, "process-request")
	defer span.End()

	// Log with trace context for correlation in observability platforms
	cfg.Logger.Info("Processing started",
		zap.String("operation", "data-fetch"),
		zap.Any("context", ctx),
	)

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Create a child span for a sub-operation
	ctx, childSpan := tracer.Start(ctx, "database-query")
	defer childSpan.End()

	// Log in the context of the child span
	cfg.Logger.Debug("Executing database query",
		zap.String("query", "SELECT * FROM users"),
		zap.Any("context", ctx),
	)

	// Add span events for important operations
	childSpan.AddEvent("query completed")

	// Log results
	cfg.Logger.Info("Operation completed successfully",
		zap.Int("results", 42),
		zap.Any("context", ctx),
	)
}

// HandleRequest shows how to use the logger in a typical HTTP request handler
func HandleRequest(ctx context.Context, cfg *configs.Configs) {
	// Extract the current span from context (assuming it was created by middleware)
	span := trace.SpanFromContext(ctx)

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("user.id", "user-123"),
		attribute.String("request.type", "GET"),
	)

	// Log with span context
	cfg.Logger.Info("Request processing started",
		zap.String("path", "/api/users"),
		zap.String("method", "GET"),
		zap.Any("context", ctx),
	)

	// Database operation example
	cfg.Logger.Debug("Database query executed",
		zap.String("query", "SELECT * FROM users WHERE id = ?"),
		zap.String("user_id", "user-123"),
		zap.Any("context", ctx),
	)

	// Error handling example
	if err := configsBuilderPerformOperation(); err != nil {
		cfg.Logger.Error("Operation failed",
			zap.Error(err),
			zap.String("operation", "user_lookup"),
			zap.Any("context", ctx),
		)
	}

	cfg.Logger.Info("Request completed",
		zap.Int("status_code", 200),
		zap.Duration("latency", time.Millisecond*45),
		zap.Any("context", ctx),
	)
}

func configsBuilderPerformOperation() error {
	// This is just a placeholder
	return nil
}
