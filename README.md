# Gokit Logging

<p align="center">
  <a href="https://github.com/goxkit/logging/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  </a>
  <a href="https://pkg.go.dev/github.com/goxkit/logging">
    <img src="https://godoc.org/github.com/goxkit/logging?status.svg" alt="Go Doc">
  </a>
  <a href="https://goreportcard.com/report/github.com/goxkit/logging">
    <img src="https://goreportcard.com/badge/github.com/goxkit/logging" alt="Go Report Card">
  </a>
  <a href="https://github.com/goxkit/logging/actions">
    <img src="https://github.com/goxkit/logging/actions/workflows/action.yml/badge.svg?branch=main" alt="Build Status">
  </a>
</p>

A comprehensive Go logging package built on [Zap](https://github.com/uber-go/zap) with OpenTelemetry integration, supporting observability-first application development. This package enables high-performance structured logging with multiple output options including OTLP export to observability platforms.

## Features

- **Multiple Output Modes**:
  - OpenTelemetry Protocol (OTLP) export for observability platforms
  - Standard output with environment-specific formatting
  - No-operation mode for testing scenarios

- **Environment-Aware Configuration**:
  - Development: Colored, human-readable console output
  - Production/Staging: JSON formatted logs for better machine parsing

- **OpenTelemetry Integration**:
  - Correlation between logs, traces, and metrics
  - Service and environment context automatically added
  - Batch processing for efficient export

- **Performance-Focused**:
  - Zero allocations during logging (via Zap)
  - Minimal CPU overhead
  - Efficient batching and export

- **Testing Support**:
  - Mock logger implementation
  - Easy integration with testify

## Installation

```bash
go get github.com/goxkit/logging
```

## Usage

### Using with ConfigsBuilder (Recommended)

The easiest way to set up logging is through the ConfigsBuilder, which handles all the necessary configuration:

```go
package main

import (
	"github.com/goxkit/configs_builder"
	"go.uber.org/zap"
)

func main() {
	// Create configurations with OTLP enabled
	cfgs, err := configsBuilder.NewConfigsBuilder().Otlp().Build()
	if err != nil {
		panic(err)
	}

	// Logger is already configured and available in cfgs.Logger
	cfgs.Logger.Info("Application started",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8080))

	// Use the logger throughout your application
	// For example, in an HTTP handler
	cfgs.Logger.Debug("Processing request",
		zap.String("path", "/api/users"),
		zap.String("method", "GET"))
}
```

### Manual Setup

If you need more control over the setup:

```go
package main

import (
	"github.com/goxkit/configs"
	"github.com/goxkit/logging"
	"go.uber.org/zap"
)

func main() {
	// Create app configurations
	appConfigs := &configs.Configs{
		AppConfigs: &configs.AppConfigs{
			Name:        "MyService",
			Namespace:   "my-namespace",
			Environment: configs.DevelopmentEnv,
			LogLevel:    configs.DEBUG,
		},
		OTLPConfigs: &configs.OTLPConfigs{
			Enabled:  true,
			Endpoint: "localhost:4317",
		},
	}

	// Initialize the logger
	logger, err := logging.NewLogger(appConfigs)
	if err != nil {
		panic(err)
	}

	// Use structured logging with context fields
	logger.Info("Service initialized",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8080))
}
```

### Log Levels

The package supports multiple log levels:

```go
// Debug: Verbose information for development
logger.Debug("Connection details",
	zap.String("host", server.Host),
	zap.Int("port", server.Port))

// Info: General operational information
logger.Info("User registered",
	zap.String("user_id", user.ID),
	zap.String("email", user.Email))

// Warn: Potential issues that don't prevent operation
logger.Warn("Database connection slow",
	zap.Duration("latency", dbLatency),
	zap.String("query", query))

// Error: Issues that require attention
logger.Error("Failed to process payment",
	zap.Error(err),
	zap.String("transaction_id", txID))

// Fatal: Critical errors that halt execution
logger.Fatal("Failed to start server",
	zap.Error(err),
	zap.Int("port", config.Port))
```

### Logging with Traces

When using the OTLP exporter, logs are automatically correlated with traces when used in a traced context:

```go
import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func HandleRequest(ctx context.Context) {
	// Assuming trace is already started in the context
	span := trace.SpanFromContext(ctx)

	// Log with trace context - will be correlated in observability platforms
	logger.Info("Processing request",
    zap.Any("context", ctx))

	// Rest of your handler logic...
}
```

### Testing with MockLogger

For unit testing code that uses the logger:

```go
func TestMyHandler(t *testing.T) {
	// Create a mock logger for testing
	mockLogger := logging.NewMockLogger()

	// Configure expected calls if needed
	mockLogger.On("Info", "User updated", mock.Anything).Return()

	// Create the handler with the mock
	handler := NewUserHandler(mockLogger)

	// Test the handler
	result := handler.UpdateUser(userID, userData)

	// Verify logger was called as expected
	mockLogger.AssertExpectations(t)
}
```

## Configuration Options

### OpenTelemetry (OTLP) Configuration

When using OTLP, configure these settings in your application:

| Setting | Environment Variable | Description |
|---------|---------------------|-------------|
| Endpoint | `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP endpoint (default: `localhost:4317`) |
| Insecure | `OTEL_EXPORTER_OTLP_INSECURE` | Whether to use insecure connection (default: `true`) |
| Timeout | `OTEL_EXPORTER_OTLP_TIMEOUT` | Timeout for export operations (default: `10s`) |
| Headers | `OTEL_EXPORTER_OTLP_HEADERS` | Headers for authentication (format: `key1=value1,key2=value2`) |

### Application Configuration

| Setting | Environment Variable | Description |
|---------|---------------------|-------------|
| Name | `APP_NAME` | Service name for identification |
| Namespace | `APP_NAMESPACE` | Service namespace for grouping |
| Environment | `GO_ENV` | Application environment (`development`, `staging`, `production`) |
| LogLevel | `LOG_LEVEL` | Minimum log level (`debug`, `info`, `warn`, `error`, `panic`) |

## Best Practices

1. **Always use structured logging**:
   ```go
   // Good
   logger.Info("User login", zap.String("user_id", userID), zap.String("source_ip", ip))

   // Avoid
   logger.Info(fmt.Sprintf("User %s logged in from %s", userID, ip))
   ```

2. **Add enough context, but not too much**:
   - Include relevant information that would help troubleshooting
   - Avoid logging sensitive information (passwords, tokens, etc.)
   - Avoid excessively large payloads

3. **Use appropriate log levels**:
   - Debug: Detailed information for development/troubleshooting
   - Info: Normal application behavior, key events
   - Warn: Unexpected but handled conditions
   - Error: Issues that require attention
   - Fatal: Critical errors that prevent operation

4. **Enable OTLP in production environments** to leverage observability platforms

## License

MIT

## References

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Uber Zap](https://github.com/uber-go/zap)
- [otelzap Bridge](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/bridges/otelzap)
