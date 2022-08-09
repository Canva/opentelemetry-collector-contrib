package reporternamedecoderprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/reporternameprocessor"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "reporternamedecoder"
)

// NewFactory creates a factory for the reporternamedecoder processor.
func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesProcessor(createTracesProcessor),
	)
}

func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
	}
}

// createTracesProcessor creates an instance of reporternamedecoder for processing traces
func createTracesProcessor(
	ctx context.Context,
	params component.ProcessorCreateSettings,
	cfg config.Processor,
	next consumer.Traces,
) (component.TracesProcessor, error) {
	oCfg := cfg.(*Config)

	reporterName, err := newReporterNameDecoder(ctx, oCfg, params.Logger, next)
	if err != nil {
		// TODO: Placeholder for an error metric in the next PR
		return nil, fmt.Errorf("error creating a reportername processor: %w", err)
	}

	return processorhelper.NewTracesProcessor(
		cfg,
		next,
		reporterName.processTraces,
		processorhelper.WithCapabilities(reporterName.Capabilities()),
		processorhelper.WithStart(reporterName.Start),
		processorhelper.WithShutdown(reporterName.Shutdown))
}
