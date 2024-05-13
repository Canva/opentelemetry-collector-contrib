// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	Type = component.MustNewType("receiver_creator")
)

var (
	Type = component.MustNewType("receiver_creator")
)

const (
	LogsStability    = component.StabilityLevelAlpha
	TracesStability  = component.StabilityLevelAlpha
	MetricsStability = component.StabilityLevelBeta
)

func Meter(settings component.TelemetrySettings) metric.Meter {
	return settings.MeterProvider.Meter("otelcol")
}

func Tracer(settings component.TelemetrySettings) trace.Tracer {
	return settings.TracerProvider.Tracer("otelcol")
}
