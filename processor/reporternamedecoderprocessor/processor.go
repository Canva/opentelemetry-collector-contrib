package reporternamedecoderprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/reporternamedecoderprocessor"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

var _ component.TracesProcessor = (*reportername)(nil)

type reportername struct {
	// ReporterName processor configuration
	config *Config
	// Logger
	logger *zap.Logger
	// Next trace consumer in line
	next consumer.Traces
}

// newReporterNameDecoder creates a new instance of the reporternamedecoder processor
func newReporterNameDecoder(ctx context.Context, config *Config, logger *zap.Logger, next consumer.Traces) (*reportername, error) {
	return &reportername{
		config: config,
		logger: logger,
		next:   next,
	}, nil
}

// processTraces implements ProcessMetricsFunc. It processes the incoming data
// and returns the data to be sent to the next component
func (s *reportername) processTraces(ctx context.Context, batch ptrace.Traces) (ptrace.Traces, error) {
	for i := 0; i < batch.ResourceSpans().Len(); i++ {
		rs := batch.ResourceSpans().At(i)
		s.processResourceSpan(ctx, rs)
	}
	return batch, nil
}

// processResourceSpan processes the RS and all of its spans and then returns the last
// view metric context. The context can be used for tests
func (s *reportername) processResourceSpan(ctx context.Context, rs ptrace.ResourceSpans) {
	rsAttrs := rs.Resource().Attributes()

	// Attributes can be part of a resource span
	s.processAttrs(ctx, &rsAttrs)

	for j := 0; j < rs.ScopeSpans().Len(); j++ {
		ils := rs.ScopeSpans().At(j)
		for k := 0; k < ils.Spans().Len(); k++ {
			span := ils.Spans().At(k)
			spanAttrs := span.Attributes()

			// Attributes can also be part of span
			s.processAttrs(ctx, &spanAttrs)
		}
	}
}

// processAttrs processes the attributes of a resource span or a span
func (s *reportername) processAttrs(_ context.Context, attributes *pcommon.Map) {
}

// ConsumeTraces implements the SpanProcessor interface
func (s *reportername) ConsumeTraces(ctx context.Context, batch ptrace.Traces) error {
	batch, err := s.processTraces(ctx, batch)
	if err != nil {
		return err
	}

	err = s.next.ConsumeTraces(ctx, batch)
	return err
}

const (
	debug = "debug"
	info  = "info"
)

// Capabilities specifies what this processor does, such as whether it mutates data
func (s *reportername) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start the reporternamedecoder processor
func (s *reportername) Start(_ context.Context, _ component.Host) error {
	return nil
}

// Shutdown the reporternamedecoder processor
func (s *reportername) Shutdown(context.Context) error {
	return nil
}
