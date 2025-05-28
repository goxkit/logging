// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

package noop

import (
	"github.com/goxkit/configs"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"

	zapInstance "github.com/goxkit/logging/zap"
)

func Install(cfgs *configs.Configs) (*zap.Logger, error) {
	provider := sdklog.NewLoggerProvider()
	cfgs.LoggerProvider = provider
	return zapInstance.NewStdoutZapLogger(cfgs)
}
