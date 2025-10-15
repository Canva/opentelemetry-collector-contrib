// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

<<<<<<<< HEAD:receiver/awscontainerinsightreceiver/internal/cadvisor/package_test.go
package cadvisor
|||||||| 6b1d3dd2c0c:exporter/datadogexporter/internal/hostmetadata/valid/package_test.go
package valid
========
package parser
>>>>>>>> v0.137.0:receiver/statsdreceiver/internal/parser/package_test.go

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
