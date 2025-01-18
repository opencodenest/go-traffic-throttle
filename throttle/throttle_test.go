package throttle

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	// Test input data
	data := []byte("This is a test payload for throttled reading.")
	reader := bytes.NewReader(data)

	// Throttle the reader to 8 kilobits per second (1 KB/sec or 1024 bytes/sec)
	const testKbps = 8.0
	throttledReader := Throttle(reader, testKbps)

	// Measure the time taken to read the full data
	start := time.Now()
	buffer := make([]byte, len(data))
	n, err := io.ReadFull(throttledReader, buffer)
	elapsed := time.Since(start)

	// Check if the throttling worked as expected
	if err != nil && err != io.EOF {
		t.Fatalf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Read incorrect number of bytes: expected %d, got %d", len(data), n)
	}

	// Expected time: data size in bytes * (1 sec / speed in bytes/sec)
	expectedDuration := time.Duration(float64(len(data)*8) / testKbps * float64(time.Second))

	// Allow a small margin for processing overhead
	margin := 50 * time.Millisecond
	if elapsed < expectedDuration-margin || elapsed > expectedDuration+margin {
		t.Errorf("Throttle duration mismatch: expected ~%v, got %v", expectedDuration, elapsed)
	}
}

func TestThrottlePartialReads(t *testing.T) {
	// Input data for testing partial reads
	data := []byte("Testing partial reads with throttled reader.")
	reader := bytes.NewReader(data)

	// Limit bandwidth to 16 kilobits per second (2 KB/sec or 2048 bytes/sec)
	const testKbps = 16.0
	throttledReader := Throttle(reader, testKbps)

	// Read data in chunks
	chunkSize := 10
	buffer := make([]byte, chunkSize)
	totalRead := 0

	for {
		n, err := throttledReader.Read(buffer)
		totalRead += n

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	// Ensure all data was read
	if totalRead != len(data) {
		t.Errorf("Partial reads did not match: expected %d bytes, got %d bytes", len(data), totalRead)
	}
}

func TestThrottleEmptyReader(t *testing.T) {
	// Test an empty reader
	reader := bytes.NewReader(nil)

	// Limit bandwidth (value doesn't matter for an empty reader)
	const testKbps = 1.0
	throttledReader := Throttle(reader, testKbps)

	// Attempt to read from the empty throttled reader
	buffer := make([]byte, 10)
	n, err := throttledReader.Read(buffer)

	if n != 0 {
		t.Errorf("Expected 0 bytes read from empty reader, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("Expected EOF error for empty reader, got %v", err)
	}
}
