// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package emittest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/emit"
)

func TestNop(t *testing.T) {
	require.NoError(t, Nop(t.Context(), [][]byte{}, map[string]any{}, int64(0), []int64{}))
}
