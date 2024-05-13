// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ucal // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/processscraper/ucal"

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var normalizeProcessCPUUtilizationFeatureGate = featuregate.GlobalRegistry().MustRegister(
	"receiver.hostmetrics.normalizeProcessCPUUtilization",
	featuregate.StageBeta,
	featuregate.WithRegisterDescription("When enabled, normalizes the process.cpu.utilization metric onto the interval [0-1] by dividing the value by the number of logical processors."),
	featuregate.WithRegisterFromVersion("v0.97.0"),
	featuregate.WithRegisterReferenceURL("https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/31368"),
)

// CPUUtilization stores the utilization percents [0-1] for the different cpu states
type CPUUtilization struct {
	User   float64
	System float64
	Iowait float64
}

// CPUUtilizationCalculator calculates the cpu utilization percents for the different cpu states
// It requires 2 []cpu.TimesStat and spend time to be able to calculate the difference
type CPUUtilizationCalculator struct {
	previousCPUStats *cpu.TimesStat
	previousReadTime pcommon.Timestamp
}

// CalculateAndRecord calculates the cpu utilization for the different cpu states comparing previously
// stored []cpu.TimesStat and time.Time and current []cpu.TimesStat and current time.Time
// If no previous data is stored it will return empty slice of CPUUtilization and no error
func (c *CPUUtilizationCalculator) CalculateAndRecord(now pcommon.Timestamp, logicalCores int, currentCPUStats *cpu.TimesStat, recorder func(pcommon.Timestamp, CPUUtilization)) error {
	if c.previousCPUStats != nil {
		recorder(now, cpuUtilization(logicalCores, c.previousCPUStats, c.previousReadTime, currentCPUStats, now))
	}
	c.previousCPUStats = currentCPUStats
	c.previousReadTime = now

	return nil
}

// cpuUtilization calculates the difference between 2 cpu.TimesStat using spent time between them
func cpuUtilization(logicalCores int, startStats *cpu.TimesStat, startTime pcommon.Timestamp, endStats *cpu.TimesStat, endTime pcommon.Timestamp) CPUUtilization {
	elapsedTime := time.Duration(endTime - startTime).Seconds()
	if elapsedTime <= 0 {
		return CPUUtilization{}
	}

	userUtilization := (endStats.User - startStats.User) / elapsedTime
	systemUtilization := (endStats.System - startStats.System) / elapsedTime
	ioWaitUtilization := (endStats.Iowait - startStats.Iowait) / elapsedTime

	if normalizeProcessCPUUtilizationFeatureGate.IsEnabled() && logicalCores > 0 {
		// Normalize onto the [0-1] interval by dividing by the number of logical cores
		userUtilization /= float64(logicalCores)
		systemUtilization /= float64(logicalCores)
		ioWaitUtilization /= float64(logicalCores)
	}

	return CPUUtilization{
		User:   userUtilization,
		System: systemUtilization,
		Iowait: ioWaitUtilization,
	}
}
