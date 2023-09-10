// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type testConfigCollection int

const (
	testSetDefault testConfigCollection = iota
	testSetAll
	testSetNone
)

func TestMetricsBuilder(t *testing.T) {
	tests := []struct {
		name      string
		configSet testConfigCollection
	}{
		{
			name:      "default",
			configSet: testSetDefault,
		},
		{
			name:      "all_set",
			configSet: testSetAll,
		},
		{
			name:      "none_set",
			configSet: testSetNone,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := pcommon.Timestamp(1_000_000_000)
			ts := pcommon.Timestamp(1_000_001_000)
			observedZapCore, observedLogs := observer.New(zap.WarnLevel)
			settings := receivertest.NewNopCreateSettings()
			settings.Logger = zap.New(observedZapCore)
			mb := NewMetricsBuilder(loadMetricsBuilderConfig(t, test.name), settings, WithStartTime(start))

			expectedWarnings := 0
			if test.configSet == testSetDefault {
				assert.Equal(t, "[WARNING] Please set `enabled` field explicitly for `default.metric`: This metric will be disabled by default soon.", observedLogs.All()[expectedWarnings].Message)
				expectedWarnings++
			}
			if test.configSet == testSetDefault || test.configSet == testSetAll {
				assert.Equal(t, "[WARNING] `default.metric.to_be_removed` should not be enabled: This metric is deprecated and will be removed soon.", observedLogs.All()[expectedWarnings].Message)
				expectedWarnings++
			}
			if test.configSet == testSetAll || test.configSet == testSetNone {
				assert.Equal(t, "[WARNING] `optional.metric` should not be configured: This metric is deprecated and will be removed soon.", observedLogs.All()[expectedWarnings].Message)
				expectedWarnings++
			}
			assert.Equal(t, expectedWarnings, observedLogs.Len())

			defaultMetricsCount := 0
			allMetricsCount := 0

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordDefaultMetricDataPoint(ts, 1, "string_attr-val", 19, AttributeEnumAttrRed, []any{"slice_attr-item1", "slice_attr-item2"}, map[string]any{"key1": "map_attr-val1", "key2": "map_attr-val2"})

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordDefaultMetricToBeRemovedDataPoint(ts, 1)

			allMetricsCount++
			mb.RecordOptionalMetricDataPoint(ts, 1, "string_attr-val", true)

			rb := mb.NewResourceBuilder()
			rb.SetMapResourceAttr(map[string]any{"key1": "map.resource.attr-val1", "key2": "map.resource.attr-val2"})
			rb.SetOptionalResourceAttr("optional.resource.attr-val")
			rb.SetSliceResourceAttr([]any{"slice.resource.attr-item1", "slice.resource.attr-item2"})
			rb.SetStringEnumResourceAttrOne()
			rb.SetStringResourceAttr("string.resource.attr-val")
			res := rb.Emit()
			metrics := mb.Emit(WithResource(res))

			if test.configSet == testSetNone {
				assert.Equal(t, 0, metrics.ResourceMetrics().Len())
				return
			}

			assert.Equal(t, 1, metrics.ResourceMetrics().Len())
			rm := metrics.ResourceMetrics().At(0)
			assert.Equal(t, res, rm.Resource())
			assert.Equal(t, 1, rm.ScopeMetrics().Len())
			ms := rm.ScopeMetrics().At(0).Metrics()
			if test.configSet == testSetDefault {
				assert.Equal(t, defaultMetricsCount, ms.Len())
			}
			if test.configSet == testSetAll {
				assert.Equal(t, allMetricsCount, ms.Len())
			}
			validatedMetrics := make(map[string]bool)
			for i := 0; i < ms.Len(); i++ {
				switch ms.At(i).Name() {
				case "default.metric":
					assert.False(t, validatedMetrics["default.metric"], "Found a duplicate in the metrics slice: default.metric")
					validatedMetrics["default.metric"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Monotonic cumulative sum int metric enabled by default.", ms.At(i).Description())
					assert.Equal(t, "s", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("string_attr")
					assert.True(t, ok)
					assert.EqualValues(t, "string_attr-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("state")
					assert.True(t, ok)
					assert.EqualValues(t, 19, attrVal.Int())
					attrVal, ok = dp.Attributes().Get("enum_attr")
					assert.True(t, ok)
					assert.EqualValues(t, "red", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("slice_attr")
					assert.True(t, ok)
					assert.EqualValues(t, []any{"slice_attr-item1", "slice_attr-item2"}, attrVal.Slice().AsRaw())
					attrVal, ok = dp.Attributes().Get("map_attr")
					assert.True(t, ok)
					assert.EqualValues(t, map[string]any{"key1": "map_attr-val1", "key2": "map_attr-val2"}, attrVal.Map().AsRaw())
				case "default.metric.to_be_removed":
					assert.False(t, validatedMetrics["default.metric.to_be_removed"], "Found a duplicate in the metrics slice: default.metric.to_be_removed")
					validatedMetrics["default.metric.to_be_removed"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "[DEPRECATED] Non-monotonic delta sum double metric enabled by default.", ms.At(i).Description())
					assert.Equal(t, "s", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityDelta, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, dp.ValueType())
					assert.Equal(t, float64(1), dp.DoubleValue())
				case "optional.metric":
					assert.False(t, validatedMetrics["optional.metric"], "Found a duplicate in the metrics slice: optional.metric")
					validatedMetrics["optional.metric"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "[DEPRECATED] Gauge double metric disabled by default.", ms.At(i).Description())
					assert.Equal(t, "1", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, dp.ValueType())
					assert.Equal(t, float64(1), dp.DoubleValue())
					attrVal, ok := dp.Attributes().Get("string_attr")
					assert.True(t, ok)
					assert.EqualValues(t, "string_attr-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("boolean_attr")
					assert.True(t, ok)
					assert.EqualValues(t, true, attrVal.Bool())
				}
			}
		})
	}
}
