// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package testbed // import "github.com/open-telemetry/opentelemetry-collector-contrib/testbed/testbed"

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// MultiReceiverTestCase defines a running test case.
type MultiReceiverTestCase struct {
	T *testing.T

	// Directory where test case results and logs will be written.
	resultDir string

	// Resource spec for agent.
	resourceSpec ResourceSpec

	// Agent process.
	agentProc OtelcolRunner

	Sender    DataSender
	Receivers []DataReceiver

	LoadGenerator *ProviderSender
	MockBackends  map[DataReceiver]*MockBackend
	validator     *MultiReceiverTestCaseValidator

	startTime time.Time

	// errorSignal indicates an error in the test case execution, e.g. process execution
	// failure or exceeding resource consumption, etc. The actual error message is already
	// logged, this is only an indicator on which you can wait to be informed.
	errorSignal chan struct{}
	// Duration is the requested duration of the tests. Configured via TESTBED_DURATION
	// env variable and defaults to 15 seconds if env variable is unspecified.
	Duration       time.Duration
	doneSignal     chan struct{}
	errorCause     string
	ResultsSummary TestResultsSummary
}

// NewMultiReceiverTestCase creates a new MultiReceiverTestCase. It expects agent-config.yaml in the specified directory.
func NewMultiReceiverTestCase(
	t *testing.T,
	dataProvider DataProvider,
	sender DataSender,
	receivers []DataReceiver,
	agentProc OtelcolRunner,
	validator *MultiReceiverTestCaseValidator,
	resultsSummary TestResultsSummary,
	opts ...MultiReceiverTestCaseOption,
) *MultiReceiverTestCase {
	tc := MultiReceiverTestCase{
		T:              t,
		errorSignal:    make(chan struct{}),
		doneSignal:     make(chan struct{}),
		startTime:      time.Now(),
		Sender:         sender,
		Receivers:      receivers,
		agentProc:      agentProc,
		validator:      validator,
		ResultsSummary: resultsSummary,
		MockBackends:   make(map[DataReceiver]*MockBackend),
	}

	// Get requested test case duration from env variable.
	duration := os.Getenv(testcaseDurationVar)
	if duration == "" {
		duration = "15s"
	}
	var err error
	tc.Duration, err = time.ParseDuration(duration)
	if err != nil {
		log.Fatalf("Invalid "+testcaseDurationVar+": %v. Expecting a valid duration string.", duration)
	}

	// Apply all provided options.
	for _, opt := range opts {
		opt(&tc)
	}

	// Prepare directory for results.
	tc.resultDir, err = filepath.Abs(path.Join("results", t.Name()))
	require.NoErrorf(t, err, "Cannot resolve %s", t.Name())
	require.NoErrorf(t, os.MkdirAll(tc.resultDir, os.ModePerm), "Cannot create directory %s", tc.resultDir)

	// Set default resource check period.
	tc.resourceSpec.ResourceCheckPeriod = 3 * time.Second
	if tc.Duration < tc.resourceSpec.ResourceCheckPeriod {
		// Resource check period should not be longer than entire test duration.
		tc.resourceSpec.ResourceCheckPeriod = tc.Duration
	}

	lg, err := NewLoadGenerator(dataProvider, sender)
	require.NoError(t, err, "Cannot create generator")
	tc.LoadGenerator = lg.(*ProviderSender)

	for _, receiver := range receivers {
		tc.MockBackends[receiver] = NewMockBackend(tc.composeTestResultFileName("backend.log"), receiver)
	}

	go tc.logStats()

	return &tc
}

func (tc *MultiReceiverTestCase) composeTestResultFileName(fileName string) string {
	fileName, err := filepath.Abs(path.Join(tc.resultDir, fileName))
	require.NoError(tc.T, err, "Cannot resolve %s", fileName)
	return fileName
}

// StartAgent starts the agent and redirects its standard output and standard error
// to "agent.log" file located in the test directory.
func (tc *MultiReceiverTestCase) StartAgent(args ...string) {
	logFileName := tc.composeTestResultFileName("agent.log")

	startParams := StartParams{
		Name:         "Agent",
		LogFilePath:  logFileName,
		CmdArgs:      args,
		resourceSpec: &tc.resourceSpec,
	}
	if err := tc.agentProc.Start(startParams); err != nil {
		tc.indicateError(err)
		return
	}

	// Start watching resource consumption.
	go func() {
		if err := tc.agentProc.WatchResourceConsumption(); err != nil {
			tc.indicateError(err)
		}
	}()

	endpoint := tc.LoadGenerator.Sender.GetEndpoint()
	if endpoint != nil {
		// Wait for agent to start. We consider the agent started when we can
		// connect to the port to which we intend to send load. We only do this
		// if the endpoint is not-empty, i.e. the sender does use network (some senders
		// like text log writers don't).
		tc.WaitFor(func() bool {
			conn, err := net.Dial(tc.LoadGenerator.Sender.GetEndpoint().Network(), tc.LoadGenerator.Sender.GetEndpoint().String())
			if err == nil && conn != nil {
				conn.Close()
				return true
			}
			return false
		}, fmt.Sprintf("connection to %s:%s", tc.LoadGenerator.Sender.GetEndpoint().Network(), tc.LoadGenerator.Sender.GetEndpoint().String()))
	}
}

// StopAgent stops agent process.
func (tc *MultiReceiverTestCase) StopAgent() {
	if _, err := tc.agentProc.Stop(); err != nil {
		tc.indicateError(err)
	}
}

// StartLoad starts the load generator and redirects its standard output and standard error
// to "load-generator.log" file located in the test directory.
func (tc *MultiReceiverTestCase) StartLoad(options LoadOptions) {
	tc.LoadGenerator.Start(options)
}

// StopLoad stops load generator.
func (tc *MultiReceiverTestCase) StopLoad() {
	tc.LoadGenerator.Stop()
}

// StartBackend starts the specified backend type.
func (tc *MultiReceiverTestCase) StartBackends() {

	for _, backend := range tc.MockBackends {
		require.NoError(tc.T, backend.Start(), "Cannot start backend for: "+backend.receiver.ProtocolName())
	}
}

// StopBackend stops the backend.
func (tc *MultiReceiverTestCase) StopBackend() {
	for _, backend := range tc.MockBackends {
		backend.Stop()
	}
}

// EnableRecording enables recording of all data received by MockBackend.
func (tc *MultiReceiverTestCase) EnableRecording() {
	for _, backend := range tc.MockBackends {
		backend.EnableRecording()
	}
}

// AgentMemoryInfo returns raw memory info struct about the agent
// as returned by github.com/shirou/gopsutil/process
func (tc *MultiReceiverTestCase) AgentMemoryInfo() (uint32, uint32, error) {
	stat, err := tc.agentProc.GetProcessMon().MemoryInfo()
	if err != nil {
		return 0, 0, err
	}
	return uint32(stat.RSS / mibibyte), uint32(stat.VMS / mibibyte), nil
}

// Stop stops the load generator, the agent and the backend.
func (tc *MultiReceiverTestCase) Stop() {
	// Stop monitoring the agent
	close(tc.doneSignal)

	// Stop all components
	tc.StopLoad()
	tc.StopAgent()
	tc.StopBackend()

	// Report test results
	tc.validator.RecordResults(tc)
}

// ValidateData validates data received by mock backend against what was generated and sent to the collector
// instance(s) under test by the LoadGenerator.
func (tc *MultiReceiverTestCase) ValidateData() {
	select {
	case <-tc.errorSignal:
		// Error is already signaled and recorded. Validating data is pointless.
		return
	default:
	}

	tc.validator.Validate(tc)
}

// Sleep for specified duration or until error is signaled.
func (tc *MultiReceiverTestCase) Sleep(d time.Duration) {
	select {
	case <-time.After(d):
	case <-tc.errorSignal:
	}
}

// WaitForN the specific condition for up to a specified duration. Records a test error
// if time is out and condition does not become true. If error is signaled
// while waiting the function will return false, but will not record additional
// test error (we assume that signaled error is already recorded in indicateError()).
func (tc *MultiReceiverTestCase) WaitForN(cond func() bool, duration time.Duration, errMsg any) bool {
	startTime := time.Now()

	// Start with 5 ms waiting interval between condition re-evaluation.
	waitInterval := time.Millisecond * 5

	for {
		if cond() {
			return true
		}

		select {
		case <-time.After(waitInterval):
		case <-tc.errorSignal:
			return false
		}

		// Increase waiting interval exponentially up to 500 ms.
		if waitInterval < time.Millisecond*500 {
			waitInterval *= 2
		}

		if time.Since(startTime) > duration {
			// Waited too long
			tc.T.Error("Time out waiting for", errMsg)
			return false
		}
	}
}

// WaitFor is like WaitForN but with a fixed duration of 10 seconds
func (tc *MultiReceiverTestCase) WaitFor(cond func() bool, errMsg any) bool {
	return tc.WaitForN(cond, time.Second*60, errMsg)
}

func (tc *MultiReceiverTestCase) indicateError(err error) {
	// Print to log for visibility
	log.Print(err.Error())

	// Indicate error for the test
	tc.T.Error(err.Error())

	tc.errorCause = err.Error()

	// Signal the error via channel
	close(tc.errorSignal)
}

func (tc *MultiReceiverTestCase) logStats() {
	t := time.NewTicker(tc.resourceSpec.ResourceCheckPeriod)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			tc.logStatsOnce()
		case <-tc.doneSignal:
			return
		}
	}
}

func (tc *MultiReceiverTestCase) logStatsOnce() {

	log.Printf("%s | %s | ",
		tc.agentProc.GetResourceConsumption(),
		tc.LoadGenerator.GetStats())

	for _, backend := range tc.MockBackends {
		log.Printf("%s | ",
			backend.GetStats())
	}
}
