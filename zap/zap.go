// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package zap provides Uber's Zap logger implementations with environment-specific
// configurations and OpenTelemetry integration. This package handles the core logging
// functionality, offering both standard output logging and OpenTelemetry-enabled logging
// with appropriate configuration based on the application environment.
package zap

import (
	"os"

	"github.com/goxkit/configs"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger creates a Zap logger configured for both local output and OpenTelemetry
// export. It sets up a combined core that routes log entries to both standard output
// and the OpenTelemetry logger provider, allowing logs to be displayed locally while
// also being sent to observability systems.
//
// The logger format is environment-sensitive:
// - Development/QA/Local: Console output with colored level encoding
// - Production/Staging: JSON output for better machine parsing
//
// Parameters:
//   - cfgs: Application configurations including environment and log level settings
//   - provider: OpenTelemetry logger provider for exporting logs
//
// Returns:
//   - A configured zap.Logger instance with both local and OTLP output
//   - An error if logger initialization fails
func NewZapLogger(cfgs *configs.Configs, provider *log.LoggerProvider) (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fmtEncoder := zapcore.NewJSONEncoder(encoderCfg)

	if cfgs.AppConfigs.Environment == configs.DevelopmentEnv ||
		cfgs.AppConfigs.Environment == configs.QaEnv ||
		cfgs.AppConfigs.Environment == configs.LocalEnv ||
		cfgs.AppConfigs.Environment == configs.UnknownEnv {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		fmtEncoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	stdout := zapcore.AddSync(os.Stdout)
	minLevel := mapZapLogLevel(cfgs.AppConfigs)
	defaultCore := zapcore.NewCore(fmtEncoder, stdout, minLevel)

	otelCore := otelzap.NewCore(
		cfgs.AppConfigs.Name,
		otelzap.WithLoggerProvider(provider),
	)

	combinedCore := zapcore.NewTee(defaultCore, otelCore)

	logger := zap.
		New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).
		Named(cfgs.AppConfigs.Name)

	return logger, nil
}

// NewStdoutZapLogger creates a Zap logger that only outputs to stdout without
// OpenTelemetry integration. This is used when OpenTelemetry is not enabled
// or not available, providing a standard logging solution with environment-specific
// formatting.
//
// The logger format is environment-sensitive:
// - Development/QA/Local: Console output with colored level encoding
// - Production/Staging: JSON output for better machine parsing
//
// Parameters:
//   - cfgs: Application configurations including environment and log level settings
//
// Returns:
//   - A configured zap.Logger instance for standard output
//   - An error if logger initialization fails
func NewStdoutZapLogger(cfgs *configs.Configs) (*zap.Logger, error) {
	zapLogLevel := mapZapLogLevel(cfgs.AppConfigs)

	if cfgs.AppConfigs.Environment == configs.ProductionEnv || cfgs.AppConfigs.Environment == configs.StagingEnv {
		logConfig := zap.NewProductionEncoderConfig()
		logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(logConfig)

		cfgs.Logger = zap.New(
			zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stdout),
				zapLogLevel,
			),
		).
			Named(cfgs.AppConfigs.Name)

		return cfgs.Logger, nil
	}

	logConfig := zap.NewDevelopmentEncoderConfig()
	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(logConfig)

	cfgs.Logger = zap.New(
		zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zapLogLevel,
		),
	).Named(cfgs.AppConfigs.Name)

	return cfgs.Logger, nil
}

// mapZapLogLevel converts the application config log level to the corresponding
// Zap log level. It provides appropriate mapping between the configs package
// log level constants and Zap's level constants.
//
// Parameters:
//   - e: Application configs containing the log level setting
//
// Returns:
//   - The corresponding zapcore.Level value, defaulting to InfoLevel if not recognized
func mapZapLogLevel(e *configs.AppConfigs) zapcore.Level {
	switch e.LogLevel {
	case configs.DEBUG:
		return zap.DebugLevel
	case configs.INFO:
		return zap.InfoLevel
	case configs.WARN:
		return zap.WarnLevel
	case configs.ERROR:
		return zap.ErrorLevel
	case configs.PANIC:
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
