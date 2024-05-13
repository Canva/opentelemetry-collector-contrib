// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package compress_test

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awskinesisexporter/internal/compress"
)

func GzipDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	zr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	out := bytes.Buffer{}
	if _, err = io.CopyN(&out, zr, 1024); err != nil && !errors.Is(err, io.EOF) {
		zr.Close()
		return nil, err
	}
	zr.Close()
	return out.Bytes(), nil
}

func NoopDecompress(data []byte) ([]byte, error) {
	return data, nil
}

func ZlibDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	zr, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}

	out := bytes.Buffer{}
	if _, err = io.CopyN(&out, zr, 1024); err != nil && !errors.Is(err, io.EOF) {
		zr.Close()
		return nil, err
	}
	zr.Close()
	return out.Bytes(), nil
}

func FlateDecompress(data []byte) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(data)
	zr := flate.NewReader(buf)
	out := bytes.Buffer{}
	if _, err = io.CopyN(&out, zr, 1024); err != nil && !errors.Is(err, io.EOF) {
		zr.Close()
		return nil, err
	}
	zr.Close()
	return out.Bytes(), nil
}

func TestCompressorFormats(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		format     string
		decompress func(data []byte) ([]byte, error)
	}{
		{format: "none", decompress: NoopDecompress},
		{format: "noop", decompress: NoopDecompress},
		{format: "gzip", decompress: GzipDecompress},
		{format: "zlib", decompress: ZlibDecompress},
		{format: "flate", decompress: FlateDecompress},
	}

	const data = "You know nothing Jon Snow"

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("format_%s", tc.format), func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			c, err := compress.NewCompressor(tc.format, logger)
			require.NoError(t, err, "Must have a valid compression format")
			require.NotNil(t, c, "Must have a valid compressor")

			out, err := c([]byte(data))
			assert.NoError(t, err, "Must not error when processing data")
			assert.NotNil(t, out, "Must have a valid record")
			outDecompress, err := tc.decompress(out)
			assert.NoError(t, err, "Decompression must have no errors")
			assert.Equal(t, []byte(data), outDecompress, "Data input should be the same after compression and decompression")
		})
	}
	_, err := compress.NewCompressor("invalid-format", zaptest.NewLogger(t))
	assert.Error(t, err, "Must error when an invalid compression format is given")
}

func BenchmarkNoopCompressor_1000Bytes(b *testing.B) {
	benchmarkCompressor(b, "none", 1000)
}

func BenchmarkNoopCompressor_1Mb(b *testing.B) {
	benchmarkCompressor(b, "noop", 131072)
}

func BenchmarkZlibCompressor_1000Bytes(b *testing.B) {
	benchmarkCompressor(b, "zlib", 1000)
}

func BenchmarkZlibCompressor_1Mb(b *testing.B) {
	benchmarkCompressor(b, "zlib", 131072)
}

func BenchmarkFlateCompressor_1000Bytes(b *testing.B) {
	benchmarkCompressor(b, "flate", 1000)
}

func BenchmarkFlateCompressor_1Mb(b *testing.B) {
	benchmarkCompressor(b, "flate", 131072)
}

func BenchmarkGzipCompressor_1000Bytes(b *testing.B) {
	benchmarkCompressor(b, "gzip", 1000)
}

func BenchmarkGzipCompressor_1Mb(b *testing.B) {
	benchmarkCompressor(b, "gzip", 131072)
}

func benchmarkCompressor(b *testing.B, format string, length int) {
	b.Helper()

	source := rand.NewSource(time.Now().UnixMilli())
	genRand := rand.New(source)

	compressor, err := compress.NewCompressor(format, zaptest.NewLogger(b))
	require.NoError(b, err, "Must not error when given a valid format")
	require.NotNil(b, compressor, "Must have a valid compressor")

	data := make([]byte, length)
	for i := 0; i < length; i++ {
		data[i] = byte(genRand.Int31())
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		out, err := compressor(data)
		assert.NoError(b, err, "Must not error when processing data")
		assert.NotNil(b, out, "Must have a valid byte array after")
	}
}

// an issue encountered in the past was a crash due race condition in the compressor, so the
// current implementation creates a new context on each compression request
// this is a test to check no exceptions are raised for executing concurrent compressions
func TestCompressorConcurrent(t *testing.T) {

	timeout := time.After(15 * time.Second)
	done := make(chan bool)
	go func() {
		// do your testing
		concurrentCompressFunc(t)
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}

}

func concurrentCompressFunc(t *testing.T) {
	// this value should be way higher to make this test more valuable, but the make of this project uses
	// max 4 workers, so we had to set this value here
	numWorkers := 4

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	errCh := make(chan error, numWorkers)
	var errMutex sync.Mutex

	// any single format would do it here, since each exporter can be set to use only one at a time
	// and the concurrent issue that was present in the past was independent of the format
	compressFunc, err := compress.NewCompressor("gzip")

	if err != nil {
		errCh <- err
		return
	}

	// it is important for the data length to be on the higher side of a record
	// since it is where the chances of having race conditions are bigger
	dataLength := 131072

	for j := 0; j < numWorkers; j++ {
		go func() {
			defer wg.Done()

			source := rand.NewSource(time.Now().UnixMilli())
			genRand := rand.New(source)

			data := make([]byte, dataLength)
			for i := 0; i < dataLength; i++ {
				data[i] = byte(genRand.Int31())
			}

			result, localErr := compressFunc(data)
			if localErr != nil {
				errMutex.Lock()
				errCh <- localErr
				errMutex.Unlock()
				return
			}

			_ = result
		}()
	}

	wg.Wait()

	close(errCh)

	for err := range errCh {
		t.Errorf("Error encountered on concurrent compression: %v", err)
	}
}
