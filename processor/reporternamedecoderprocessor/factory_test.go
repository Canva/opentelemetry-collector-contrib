package reporternamedecoderprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumertest"
)

func TestDefaultConfiguration(t *testing.T) {
	c := createDefaultConfig().(*Config)
	assert.Empty(t, c.S3BucketUrl)
}

func TestCreateTestProcessor(t *testing.T) {
	cfg := &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
	}

	tp, err := createTracesProcessor(context.Background(), componenttest.NewNopProcessorCreateSettings(), cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, tp)
	assert.Equal(t, true, tp.Capabilities().MutatesData)
}
