// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package otlp provides OpenTelemetry Protocol (OTLP) integration for logging,
// enabling logs to be exported to observability platforms like Jaeger, Prometheus,
// or other compatible collectors. This package connects the application logging
// to distributed tracing and metrics systems for comprehensive observability.
package otlp

import (
	"context"
	"time"

	"github.com/goxkit/configs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	zapInstance "github.com/goxkit/logging/zap"
)

// Install configures and initializes an OpenTelemetry-enabled logger that exports
// logs to an OTLP collector. It sets up the connection to the OTLP endpoint specified
// in the configuration and configures the logger with proper service and environment
// attributes for better observability context.
//
// The function handles the complete setup of:
// - OTLP exporter with gRPC transport
// - Batch processing for efficient log export
// - Resource attributes for service identification
// - Global logger provider registration
// - Integration with Zap for structured logging
//
// Parameters:
//   - cfgs: Application configurations including OTLP endpoint and service information
//
// Returns:
//   - A configured zap.Logger instance with OTLP export capabilities
//   - An error if the OTLP exporter or logger initialization fails
func Install(cfgs *configs.Configs) (*zap.Logger, error) {
	ctx := context.Background()

	exp, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(cfgs.OTLPConfigs.Endpoint),
		otlploggrpc.WithReconnectionPeriod(cfgs.OTLPConfigs.ExporterReconnectionPeriod),
		otlploggrpc.WithTimeout(cfgs.OTLPConfigs.ExporterTimeout),
		otlploggrpc.WithCompressor("gzip"),
		otlploggrpc.WithDialOption(
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1 * time.Second,
					Multiplier: 1.6,
					MaxDelay:   15 * time.Second,
				},
				MinConnectTimeout: 0,
			}),
		),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	processor := sdklog.NewBatchProcessor(exp)
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfgs.AppConfigs.Name),
			semconv.ServiceNamespaceKey.String(cfgs.AppConfigs.Namespace),
			attribute.String("service.environment", cfgs.AppConfigs.Environment.String()),
			semconv.DeploymentEnvironmentKey.String(cfgs.AppConfigs.Environment.String()),
			semconv.TelemetrySDKLanguageKey.String("go"),
			semconv.TelemetrySDKLanguageGo.Key.Bool(true),
		)),
	)

	global.SetLoggerProvider(provider)
	cfgs.LoggerProvider = provider

	return zapInstance.NewZapLogger(cfgs, provider)
}
