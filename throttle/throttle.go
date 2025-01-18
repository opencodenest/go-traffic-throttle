// Package throttle provides a way to limit the bandwidth of reading from an io.Reader.
package throttle

import (
	"io"
	"time"
)

// Throttle wraps an io.Reader to limit its bandwidth to the specified kbps (kilobits per second).
// It returns a new io.Reader that will read at the specified rate (kbps).
//
// For example, if you want to throttle an io.Reader to 100 kilobits per second,
// you would call:
//   throttledReader := throttle.Throttle(reader, 100)
//
// This will wrap the original reader and ensure data is read at a rate of 100 kilobits per second.
func Throttle(reader io.Reader, kbps float64) io.Reader {
	// Convert the kbps to the duration to sleep between reads
	kb := 1.0
	sleepDuration := time.Duration(kb / kbps * float64(time.Second))
	return &readerWrapper{reader: reader, sleepDuration: sleepDuration}
}

// readerWrapper is a wrapper type that implements the io.Reader interface.
// It adds throttling by sleeping between reads to control the bandwidth.
type readerWrapper struct {
	reader        io.Reader
	sleepDuration time.Duration
}

// Read implements the io.Reader interface. It reads data from the wrapped reader,
// but enforces a bandwidth limit by sleeping between reads.
func (w *readerWrapper) Read(p []byte) (n int, err error) {
	start := time.Now()

	// Limit the number of bytes to read at a time to 125 (1 kilobit)
	readSize := 125
	if len(p) < readSize {
		readSize = len(p)
	}

	// Read from the original reader
	n, err = w.reader.Read(p[:readSize])

	// Calculate elapsed time and sleep to enforce throttling
	elapsedTime := time.Since(start)
	if elapsedTime < w.sleepDuration {
		time.Sleep(w.sleepDuration - elapsedTime)
	}

	return n, err
}
