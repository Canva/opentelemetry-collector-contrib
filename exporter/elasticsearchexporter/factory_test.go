// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package elasticsearchexporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter/internal/metadata"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestFactory_CreateLogs(t *testing.T) {
	for _, tc := range signalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			factory := NewFactory()
			params := exportertest.NewNopSettings(metadata.Type)

			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			params.Logger = zap.New(observedZapCore)

			cfg := createDefaultConfig()
			cm := confmap.NewFromStringMap(tc.cfg)
			require.NoError(t, cm.Unmarshal(cfg))

			exporter, err := factory.CreateLogs(t.Context(), params, cfg.(*Config))
			require.NoError(t, err)
			require.NotNil(t, exporter)
			require.Equal(t, len(tc.expectedLogs), observedLogs.Len())
			actualLogs := observedLogs.All()
			for i, expectedLog := range tc.expectedLogs {
				assert.Contains(t, actualLogs[i].Message, expectedLog)
			}

			require.NoError(t, exporter.Shutdown(t.Context()))
		})
	}
}

func TestFactory_CreateMetrics(t *testing.T) {
	for _, tc := range signalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			factory := NewFactory()
			params := exportertest.NewNopSettings(metadata.Type)

			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			params.Logger = zap.New(observedZapCore)

			cfg := createDefaultConfig()
			cm := confmap.NewFromStringMap(tc.cfg)
			require.NoError(t, cm.Unmarshal(cfg))

			exporter, err := factory.CreateMetrics(t.Context(), params, cfg.(*Config))
			require.NoError(t, err)
			require.NotNil(t, exporter)
			require.Equal(t, len(tc.expectedLogs), observedLogs.Len())
			actualLogs := observedLogs.All()
			for i, expectedLog := range tc.expectedLogs {
				assert.Contains(t, actualLogs[i].Message, expectedLog)
			}

			require.NoError(t, exporter.Shutdown(t.Context()))
		})
	}
}

func TestFactory_CreateTraces(t *testing.T) {
	for _, tc := range signalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			factory := NewFactory()
			params := exportertest.NewNopSettings(metadata.Type)

			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			params.Logger = zap.New(observedZapCore)

			cfg := createDefaultConfig()
			cm := confmap.NewFromStringMap(tc.cfg)
			require.NoError(t, cm.Unmarshal(cfg))

			exporter, err := factory.CreateTraces(t.Context(), params, cfg.(*Config))
			require.NoError(t, err)
			require.NotNil(t, exporter)
			require.Equal(t, len(tc.expectedLogs), observedLogs.Len())
			actualLogs := observedLogs.All()
			for i, expectedLog := range tc.expectedLogs {
				assert.Contains(t, actualLogs[i].Message, expectedLog)
			}

			require.NoError(t, exporter.Shutdown(t.Context()))
		})
	}
}

var signalTestCases = []struct {
	name         string
	cfg          map[string]any
	expectedLogs []string
}{
	{
		name: "default",
		cfg: map[string]any{
			"endpoints": []string{"http://test:9200"},
		},
	},
	{
		name: "with_deprecated_batcher",
		cfg: map[string]any{
			"batcher": map[string]any{
				"enabled": true,
			},
		},
		expectedLogs: []string{
			"batcher has been deprecated",
		},
	},
	{
		name: "with_sending_queue",
		cfg: map[string]any{
			"sending_queue": map[string]any{
				"enabled": true,
				"batch":   map[string]any{},
			},
		},
	},
	{
		name: "with_sending_queue_disabled_and_deprecated_batcher",
		cfg: map[string]any{
			"sending_queue": map[string]any{
				"batch": map[string]any{},
			},
			"batcher": map[string]any{
				"enabled": true,
			},
		},
		expectedLogs: []string{
			"sending_queue::batch will take preference",
		},
	},
	{
		name: "with_sending_queue_enabled_and_deprecated_batcher",
		cfg: map[string]any{
			"sending_queue": map[string]any{
				"enabled": true,
				"batch":   map[string]any{},
			},
			"batcher": map[string]any{
				"enabled": true,
			},
		},
		expectedLogs: []string{
			"sending_queue::batch will take preference",
		},
	},
}

func TestFactory_DedupDeprecated(t *testing.T) {
	factory := NewFactory()
	cfg := withDefaultConfig(func(cfg *Config) {
		dedup := false
		cfg.Endpoint = "http://testing.invalid:9200"
		cfg.Mapping.Dedup = &dedup
		cfg.Mapping.Dedot = false // avoid dedot warnings
	})

	loggerCore, logObserver := observer.New(zap.WarnLevel)
	set := exportertest.NewNopSettings()
	set.Logger = zap.New(loggerCore)

	logsExporter, err := factory.CreateLogs(context.Background(), set, cfg)
	require.NoError(t, err)
	require.NoError(t, logsExporter.Shutdown(context.Background()))

	tracesExporter, err := factory.CreateTraces(context.Background(), set, cfg)
	require.NoError(t, err)
	require.NoError(t, tracesExporter.Shutdown(context.Background()))

	metricsExporter, err := factory.CreateMetrics(context.Background(), set, cfg)
	require.NoError(t, err)
	require.NoError(t, metricsExporter.Shutdown(context.Background()))

	records := logObserver.AllUntimed()
	assert.Len(t, records, 3)
	assert.Equal(t, "dedup is deprecated, and is always enabled", records[0].Message)
	assert.Equal(t, "dedup is deprecated, and is always enabled", records[1].Message)
	assert.Equal(t, "dedup is deprecated, and is always enabled", records[2].Message)
}

func TestFactory_DedotDeprecated(t *testing.T) {
	loggerCore, logObserver := observer.New(zap.WarnLevel)
	set := exportertest.NewNopSettings()
	set.Logger = zap.New(loggerCore)

	cfgNoDedotECS := withDefaultConfig(func(cfg *Config) {
		cfg.Endpoint = "http://testing.invalid:9200"
		cfg.Mapping.Dedot = false
		cfg.Mapping.Mode = "ecs"
	})

	cfgDedotRaw := withDefaultConfig(func(cfg *Config) {
		cfg.Endpoint = "http://testing.invalid:9200"
		cfg.Mapping.Dedot = true
		cfg.Mapping.Mode = "raw"
	})

	for _, cfg := range []*Config{cfgNoDedotECS, cfgDedotRaw} {
		factory := NewFactory()
		logsExporter, err := factory.CreateLogs(context.Background(), set, cfg)
		require.NoError(t, err)
		require.NoError(t, logsExporter.Shutdown(context.Background()))

		tracesExporter, err := factory.CreateTraces(context.Background(), set, cfg)
		require.NoError(t, err)
		require.NoError(t, tracesExporter.Shutdown(context.Background()))

		metricsExporter, err := factory.CreateMetrics(context.Background(), set, cfg)
		require.NoError(t, err)
		require.NoError(t, metricsExporter.Shutdown(context.Background()))
	}

	records := logObserver.AllUntimed()
	assert.Len(t, records, 6)
	for _, record := range records {
		assert.Equal(t, "dedot has been deprecated: in the future, dedotting will always be performed in ECS mode only", record.Message)
	}
}
