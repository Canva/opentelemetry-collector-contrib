// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

<<<<<<<< HEAD:exporter/signalfxexporter/internal/dimensions/package_test.go
package dimensions
|||||||| 6b1d3dd2c0c:exporter/instanaexporter/internal/converter/package_test.go
package converter
========
package cadvisor
>>>>>>>> v0.137.0:receiver/awscontainerinsightreceiver/internal/cadvisor/package_test.go

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
