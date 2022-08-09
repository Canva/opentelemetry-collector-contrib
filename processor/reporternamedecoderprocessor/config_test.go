package reporternamedecoderprocessor

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factories.Processors[typeStr] = NewFactory()

	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestLoadConfigEmpty(t *testing.T) {
	factories, err := componenttest.NopFactories()
	require.NoError(t, err)

	factory := NewFactory()
	factories.Processors[typeStr] = factory

	cfg, err := servicetest.LoadConfigAndValidate(path.Join(".", "testdata", "empty.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	p0 := cfg.Processors[config.NewComponentID(typeStr)]
	assert.Equal(t, p0, createDefaultConfig())
}
