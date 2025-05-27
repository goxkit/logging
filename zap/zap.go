package zap

import (
	"os"

	"github.com/goxkit/configs"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
// Zap log level. It defaults to InfoLevel if the level is not recognized.
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
