// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

<<<<<<<< HEAD:pkg/translator/azurelogs/package_test.go
package azurelogs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/azurelogs"
|||||||| 6b1d3dd2c0c:cmd/configschema/cfgmetadatagen/cfgmetadatagen/package_test.go
package cfgmetadatagen
========
package ottlprofilesample
>>>>>>>> v0.137.0:pkg/ottl/contexts/ottlprofilesample/package_test.go

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
