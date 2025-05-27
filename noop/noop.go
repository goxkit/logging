package noop

import (
	"github.com/goxkit/configs"
	"go.uber.org/zap"

	zapInstance "github.com/goxkit/logging/zap"
)

func Install(cfgs *configs.Configs) (*zap.Logger, error) {
	return zapInstance.NewStdoutZapLogger(cfgs)
}
