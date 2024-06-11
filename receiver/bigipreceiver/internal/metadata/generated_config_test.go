// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestMetricsBuilderConfig(t *testing.T) {
	tests := []struct {
		name string
		want MetricsBuilderConfig
	}{
		{
			name: "default",
			want: DefaultMetricsBuilderConfig(),
		},
		{
			name: "all_set",
			want: MetricsBuilderConfig{
				Metrics: MetricsConfig{
					BigipNodeAvailability:             MetricConfig{Enabled: true},
					BigipNodeConnectionCount:          MetricConfig{Enabled: true},
					BigipNodeDataTransmitted:          MetricConfig{Enabled: true},
					BigipNodeEnabled:                  MetricConfig{Enabled: true},
					BigipNodePacketCount:              MetricConfig{Enabled: true},
					BigipNodeRequestCount:             MetricConfig{Enabled: true},
					BigipNodeSessionCount:             MetricConfig{Enabled: true},
					BigipPoolAvailability:             MetricConfig{Enabled: true},
					BigipPoolConnectionCount:          MetricConfig{Enabled: true},
					BigipPoolDataTransmitted:          MetricConfig{Enabled: true},
					BigipPoolEnabled:                  MetricConfig{Enabled: true},
					BigipPoolMemberCount:              MetricConfig{Enabled: true},
					BigipPoolPacketCount:              MetricConfig{Enabled: true},
					BigipPoolRequestCount:             MetricConfig{Enabled: true},
					BigipPoolMemberAvailability:       MetricConfig{Enabled: true},
					BigipPoolMemberConnectionCount:    MetricConfig{Enabled: true},
					BigipPoolMemberDataTransmitted:    MetricConfig{Enabled: true},
					BigipPoolMemberEnabled:            MetricConfig{Enabled: true},
					BigipPoolMemberPacketCount:        MetricConfig{Enabled: true},
					BigipPoolMemberRequestCount:       MetricConfig{Enabled: true},
					BigipPoolMemberSessionCount:       MetricConfig{Enabled: true},
					BigipVirtualServerAvailability:    MetricConfig{Enabled: true},
					BigipVirtualServerConnectionCount: MetricConfig{Enabled: true},
					BigipVirtualServerDataTransmitted: MetricConfig{Enabled: true},
					BigipVirtualServerEnabled:         MetricConfig{Enabled: true},
					BigipVirtualServerPacketCount:     MetricConfig{Enabled: true},
					BigipVirtualServerRequestCount:    MetricConfig{Enabled: true},
				},
				ResourceAttributes: ResourceAttributesConfig{
					BigipNodeIPAddress:            ResourceAttributeConfig{Enabled: true},
					BigipNodeName:                 ResourceAttributeConfig{Enabled: true},
					BigipPoolName:                 ResourceAttributeConfig{Enabled: true},
					BigipPoolMemberIPAddress:      ResourceAttributeConfig{Enabled: true},
					BigipPoolMemberName:           ResourceAttributeConfig{Enabled: true},
					BigipVirtualServerDestination: ResourceAttributeConfig{Enabled: true},
					BigipVirtualServerName:        ResourceAttributeConfig{Enabled: true},
				},
			},
		},
		{
			name: "none_set",
			want: MetricsBuilderConfig{
				Metrics: MetricsConfig{
					BigipNodeAvailability:             MetricConfig{Enabled: false},
					BigipNodeConnectionCount:          MetricConfig{Enabled: false},
					BigipNodeDataTransmitted:          MetricConfig{Enabled: false},
					BigipNodeEnabled:                  MetricConfig{Enabled: false},
					BigipNodePacketCount:              MetricConfig{Enabled: false},
					BigipNodeRequestCount:             MetricConfig{Enabled: false},
					BigipNodeSessionCount:             MetricConfig{Enabled: false},
					BigipPoolAvailability:             MetricConfig{Enabled: false},
					BigipPoolConnectionCount:          MetricConfig{Enabled: false},
					BigipPoolDataTransmitted:          MetricConfig{Enabled: false},
					BigipPoolEnabled:                  MetricConfig{Enabled: false},
					BigipPoolMemberCount:              MetricConfig{Enabled: false},
					BigipPoolPacketCount:              MetricConfig{Enabled: false},
					BigipPoolRequestCount:             MetricConfig{Enabled: false},
					BigipPoolMemberAvailability:       MetricConfig{Enabled: false},
					BigipPoolMemberConnectionCount:    MetricConfig{Enabled: false},
					BigipPoolMemberDataTransmitted:    MetricConfig{Enabled: false},
					BigipPoolMemberEnabled:            MetricConfig{Enabled: false},
					BigipPoolMemberPacketCount:        MetricConfig{Enabled: false},
					BigipPoolMemberRequestCount:       MetricConfig{Enabled: false},
					BigipPoolMemberSessionCount:       MetricConfig{Enabled: false},
					BigipVirtualServerAvailability:    MetricConfig{Enabled: false},
					BigipVirtualServerConnectionCount: MetricConfig{Enabled: false},
					BigipVirtualServerDataTransmitted: MetricConfig{Enabled: false},
					BigipVirtualServerEnabled:         MetricConfig{Enabled: false},
					BigipVirtualServerPacketCount:     MetricConfig{Enabled: false},
					BigipVirtualServerRequestCount:    MetricConfig{Enabled: false},
				},
				ResourceAttributes: ResourceAttributesConfig{
					BigipNodeIPAddress:            ResourceAttributeConfig{Enabled: false},
					BigipNodeName:                 ResourceAttributeConfig{Enabled: false},
					BigipPoolName:                 ResourceAttributeConfig{Enabled: false},
					BigipPoolMemberIPAddress:      ResourceAttributeConfig{Enabled: false},
					BigipPoolMemberName:           ResourceAttributeConfig{Enabled: false},
					BigipVirtualServerDestination: ResourceAttributeConfig{Enabled: false},
					BigipVirtualServerName:        ResourceAttributeConfig{Enabled: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := loadMetricsBuilderConfig(t, tt.name)
			if diff := cmp.Diff(tt.want, cfg, cmpopts.IgnoreUnexported(MetricConfig{}, ResourceAttributeConfig{})); diff != "" {
				t.Errorf("Config mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}

func loadMetricsBuilderConfig(t *testing.T, name string) MetricsBuilderConfig {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	sub, err := cm.Sub(name)
	require.NoError(t, err)
	cfg := DefaultMetricsBuilderConfig()
	require.NoError(t, sub.Unmarshal(&cfg))
	return cfg
}

func TestResourceAttributesConfig(t *testing.T) {
	tests := []struct {
		name string
		want ResourceAttributesConfig
	}{
		{
			name: "default",
			want: DefaultResourceAttributesConfig(),
		},
		{
			name: "all_set",
			want: ResourceAttributesConfig{
				BigipNodeIPAddress:            ResourceAttributeConfig{Enabled: true},
				BigipNodeName:                 ResourceAttributeConfig{Enabled: true},
				BigipPoolName:                 ResourceAttributeConfig{Enabled: true},
				BigipPoolMemberIPAddress:      ResourceAttributeConfig{Enabled: true},
				BigipPoolMemberName:           ResourceAttributeConfig{Enabled: true},
				BigipVirtualServerDestination: ResourceAttributeConfig{Enabled: true},
				BigipVirtualServerName:        ResourceAttributeConfig{Enabled: true},
			},
		},
		{
			name: "none_set",
			want: ResourceAttributesConfig{
				BigipNodeIPAddress:            ResourceAttributeConfig{Enabled: false},
				BigipNodeName:                 ResourceAttributeConfig{Enabled: false},
				BigipPoolName:                 ResourceAttributeConfig{Enabled: false},
				BigipPoolMemberIPAddress:      ResourceAttributeConfig{Enabled: false},
				BigipPoolMemberName:           ResourceAttributeConfig{Enabled: false},
				BigipVirtualServerDestination: ResourceAttributeConfig{Enabled: false},
				BigipVirtualServerName:        ResourceAttributeConfig{Enabled: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := loadResourceAttributesConfig(t, tt.name)
			if diff := cmp.Diff(tt.want, cfg, cmpopts.IgnoreUnexported(ResourceAttributeConfig{})); diff != "" {
				t.Errorf("Config mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}

func loadResourceAttributesConfig(t *testing.T, name string) ResourceAttributesConfig {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	sub, err := cm.Sub(name)
	require.NoError(t, err)
	sub, err = sub.Sub("resource_attributes")
	require.NoError(t, err)
	cfg := DefaultResourceAttributesConfig()
	require.NoError(t, sub.Unmarshal(&cfg))
	return cfg
}
