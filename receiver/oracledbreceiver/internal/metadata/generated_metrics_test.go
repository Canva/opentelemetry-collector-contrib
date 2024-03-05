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

			assert.Equal(t, expectedWarnings, observedLogs.Len())

			defaultMetricsCount := 0
			allMetricsCount := 0

			allMetricsCount++
			mb.RecordOracledbConsistentGetsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbCPUTimeDataPoint(ts, 1)

			allMetricsCount++
			mb.RecordOracledbDbBlockGetsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbDmlLocksLimitDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbDmlLocksUsageDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbEnqueueDeadlocksDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbEnqueueLocksLimitDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbEnqueueLocksUsageDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbEnqueueResourcesLimitDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbEnqueueResourcesUsageDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbExchangeDeadlocksDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbExecutionsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbHardParsesDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbLogicalReadsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbParseCallsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbPgaMemoryDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbPhysicalReadsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbProcessesLimitDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbProcessesUsageDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbSessionsLimitDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbSessionsUsageDataPoint(ts, "1", "session_type-val", "session_status-val")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbTablespaceSizeLimitDataPoint(ts, 1, "tablespace_name-val")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbTablespaceSizeUsageDataPoint(ts, "1", "tablespace_name-val")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbTransactionsLimitDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbTransactionsUsageDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbUserCommitsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordOracledbUserRollbacksDataPoint(ts, "1")

			rb := mb.NewResourceBuilder()
			rb.SetOracledbInstanceName("oracledb.instance.name-val")
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
				case "oracledb.consistent_gets":
					assert.False(t, validatedMetrics["oracledb.consistent_gets"], "Found a duplicate in the metrics slice: oracledb.consistent_gets")
					validatedMetrics["oracledb.consistent_gets"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of times a consistent read was requested for a block from the buffer cache.", ms.At(i).Description())
					assert.Equal(t, "{gets}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.cpu_time":
					assert.False(t, validatedMetrics["oracledb.cpu_time"], "Found a duplicate in the metrics slice: oracledb.cpu_time")
					validatedMetrics["oracledb.cpu_time"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Cumulative CPU time, in seconds", ms.At(i).Description())
					assert.Equal(t, "s", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, dp.ValueType())
					assert.Equal(t, float64(1), dp.DoubleValue())
				case "oracledb.db_block_gets":
					assert.False(t, validatedMetrics["oracledb.db_block_gets"], "Found a duplicate in the metrics slice: oracledb.db_block_gets")
					validatedMetrics["oracledb.db_block_gets"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of times a current block was requested from the buffer cache.", ms.At(i).Description())
					assert.Equal(t, "{gets}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.dml_locks.limit":
					assert.False(t, validatedMetrics["oracledb.dml_locks.limit"], "Found a duplicate in the metrics slice: oracledb.dml_locks.limit")
					validatedMetrics["oracledb.dml_locks.limit"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Maximum limit of active DML (Data Manipulation Language) locks, -1 if unlimited.", ms.At(i).Description())
					assert.Equal(t, "{locks}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.dml_locks.usage":
					assert.False(t, validatedMetrics["oracledb.dml_locks.usage"], "Found a duplicate in the metrics slice: oracledb.dml_locks.usage")
					validatedMetrics["oracledb.dml_locks.usage"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Current count of active DML (Data Manipulation Language) locks.", ms.At(i).Description())
					assert.Equal(t, "{locks}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.enqueue_deadlocks":
					assert.False(t, validatedMetrics["oracledb.enqueue_deadlocks"], "Found a duplicate in the metrics slice: oracledb.enqueue_deadlocks")
					validatedMetrics["oracledb.enqueue_deadlocks"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Total number of deadlocks between table or row locks in different sessions.", ms.At(i).Description())
					assert.Equal(t, "{deadlocks}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.enqueue_locks.limit":
					assert.False(t, validatedMetrics["oracledb.enqueue_locks.limit"], "Found a duplicate in the metrics slice: oracledb.enqueue_locks.limit")
					validatedMetrics["oracledb.enqueue_locks.limit"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Maximum limit of active enqueue locks, -1 if unlimited.", ms.At(i).Description())
					assert.Equal(t, "{locks}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.enqueue_locks.usage":
					assert.False(t, validatedMetrics["oracledb.enqueue_locks.usage"], "Found a duplicate in the metrics slice: oracledb.enqueue_locks.usage")
					validatedMetrics["oracledb.enqueue_locks.usage"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Current count of active enqueue locks.", ms.At(i).Description())
					assert.Equal(t, "{locks}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.enqueue_resources.limit":
					assert.False(t, validatedMetrics["oracledb.enqueue_resources.limit"], "Found a duplicate in the metrics slice: oracledb.enqueue_resources.limit")
					validatedMetrics["oracledb.enqueue_resources.limit"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Maximum limit of active enqueue resources, -1 if unlimited.", ms.At(i).Description())
					assert.Equal(t, "{resources}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.enqueue_resources.usage":
					assert.False(t, validatedMetrics["oracledb.enqueue_resources.usage"], "Found a duplicate in the metrics slice: oracledb.enqueue_resources.usage")
					validatedMetrics["oracledb.enqueue_resources.usage"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Current count of active enqueue resources.", ms.At(i).Description())
					assert.Equal(t, "{resources}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.exchange_deadlocks":
					assert.False(t, validatedMetrics["oracledb.exchange_deadlocks"], "Found a duplicate in the metrics slice: oracledb.exchange_deadlocks")
					validatedMetrics["oracledb.exchange_deadlocks"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of times that a process detected a potential deadlock when exchanging two buffers and raised an internal, restartable error. Index scans are the only operations that perform exchanges.", ms.At(i).Description())
					assert.Equal(t, "{deadlocks}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.executions":
					assert.False(t, validatedMetrics["oracledb.executions"], "Found a duplicate in the metrics slice: oracledb.executions")
					validatedMetrics["oracledb.executions"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Total number of calls (user and recursive) that executed SQL statements", ms.At(i).Description())
					assert.Equal(t, "{executions}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.hard_parses":
					assert.False(t, validatedMetrics["oracledb.hard_parses"], "Found a duplicate in the metrics slice: oracledb.hard_parses")
					validatedMetrics["oracledb.hard_parses"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of hard parses", ms.At(i).Description())
					assert.Equal(t, "{parses}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.logical_reads":
					assert.False(t, validatedMetrics["oracledb.logical_reads"], "Found a duplicate in the metrics slice: oracledb.logical_reads")
					validatedMetrics["oracledb.logical_reads"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of logical reads", ms.At(i).Description())
					assert.Equal(t, "{reads}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.parse_calls":
					assert.False(t, validatedMetrics["oracledb.parse_calls"], "Found a duplicate in the metrics slice: oracledb.parse_calls")
					validatedMetrics["oracledb.parse_calls"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Total number of parse calls.", ms.At(i).Description())
					assert.Equal(t, "{parses}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.pga_memory":
					assert.False(t, validatedMetrics["oracledb.pga_memory"], "Found a duplicate in the metrics slice: oracledb.pga_memory")
					validatedMetrics["oracledb.pga_memory"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Session PGA (Program Global Area) memory", ms.At(i).Description())
					assert.Equal(t, "By", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.physical_reads":
					assert.False(t, validatedMetrics["oracledb.physical_reads"], "Found a duplicate in the metrics slice: oracledb.physical_reads")
					validatedMetrics["oracledb.physical_reads"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of physical reads", ms.At(i).Description())
					assert.Equal(t, "{reads}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.processes.limit":
					assert.False(t, validatedMetrics["oracledb.processes.limit"], "Found a duplicate in the metrics slice: oracledb.processes.limit")
					validatedMetrics["oracledb.processes.limit"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Maximum limit of active processes, -1 if unlimited.", ms.At(i).Description())
					assert.Equal(t, "{processes}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.processes.usage":
					assert.False(t, validatedMetrics["oracledb.processes.usage"], "Found a duplicate in the metrics slice: oracledb.processes.usage")
					validatedMetrics["oracledb.processes.usage"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Current count of active processes.", ms.At(i).Description())
					assert.Equal(t, "{processes}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.sessions.limit":
					assert.False(t, validatedMetrics["oracledb.sessions.limit"], "Found a duplicate in the metrics slice: oracledb.sessions.limit")
					validatedMetrics["oracledb.sessions.limit"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Maximum limit of active sessions, -1 if unlimited.", ms.At(i).Description())
					assert.Equal(t, "{sessions}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.sessions.usage":
					assert.False(t, validatedMetrics["oracledb.sessions.usage"], "Found a duplicate in the metrics slice: oracledb.sessions.usage")
					validatedMetrics["oracledb.sessions.usage"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Count of active sessions.", ms.At(i).Description())
					assert.Equal(t, "{sessions}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("session_type")
					assert.True(t, ok)
					assert.EqualValues(t, "session_type-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("session_status")
					assert.True(t, ok)
					assert.EqualValues(t, "session_status-val", attrVal.Str())
				case "oracledb.tablespace_size.limit":
					assert.False(t, validatedMetrics["oracledb.tablespace_size.limit"], "Found a duplicate in the metrics slice: oracledb.tablespace_size.limit")
					validatedMetrics["oracledb.tablespace_size.limit"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Maximum size of tablespace in bytes, -1 if unlimited.", ms.At(i).Description())
					assert.Equal(t, "By", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("tablespace_name")
					assert.True(t, ok)
					assert.EqualValues(t, "tablespace_name-val", attrVal.Str())
				case "oracledb.tablespace_size.usage":
					assert.False(t, validatedMetrics["oracledb.tablespace_size.usage"], "Found a duplicate in the metrics slice: oracledb.tablespace_size.usage")
					validatedMetrics["oracledb.tablespace_size.usage"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Used tablespace in bytes.", ms.At(i).Description())
					assert.Equal(t, "By", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("tablespace_name")
					assert.True(t, ok)
					assert.EqualValues(t, "tablespace_name-val", attrVal.Str())
				case "oracledb.transactions.limit":
					assert.False(t, validatedMetrics["oracledb.transactions.limit"], "Found a duplicate in the metrics slice: oracledb.transactions.limit")
					validatedMetrics["oracledb.transactions.limit"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Maximum limit of active transactions, -1 if unlimited.", ms.At(i).Description())
					assert.Equal(t, "{transactions}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.transactions.usage":
					assert.False(t, validatedMetrics["oracledb.transactions.usage"], "Found a duplicate in the metrics slice: oracledb.transactions.usage")
					validatedMetrics["oracledb.transactions.usage"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Current count of active transactions.", ms.At(i).Description())
					assert.Equal(t, "{transactions}", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.user_commits":
					assert.False(t, validatedMetrics["oracledb.user_commits"], "Found a duplicate in the metrics slice: oracledb.user_commits")
					validatedMetrics["oracledb.user_commits"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of user commits. When a user commits a transaction, the redo generated that reflects the changes made to database blocks must be written to disk. Commits often represent the closest thing to a user transaction rate.", ms.At(i).Description())
					assert.Equal(t, "{commits}", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "oracledb.user_rollbacks":
					assert.False(t, validatedMetrics["oracledb.user_rollbacks"], "Found a duplicate in the metrics slice: oracledb.user_rollbacks")
					validatedMetrics["oracledb.user_rollbacks"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Number of times users manually issue the ROLLBACK statement or an error occurs during a user's transactions", ms.At(i).Description())
					assert.Equal(t, "1", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				}
			}
		})
	}
}
