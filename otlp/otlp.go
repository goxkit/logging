package otlp

import (
	"context"

	"github.com/goxkit/configs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	zapInstance "github.com/goxkit/logging/zap"
)

func Install(cfgs *configs.Configs) (*zap.Logger, error) {
	ctx := context.Background()

	exp, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(cfgs.OTLPConfigs.Endpoint),
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithTimeout(0),
		otlploggrpc.WithCompressor("gzip"),
		otlploggrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, err
	}

	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(
		log.WithProcessor(processor),
		log.WithResource(resource.NewWithAttributes(
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
