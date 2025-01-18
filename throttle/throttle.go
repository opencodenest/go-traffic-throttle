package throttle

import (
	"io"
	"time"
)

// Throttle wraps an io.Reader to limit its bandwidth to the specified kbps (kilobits per second).
func Throttle(reader io.Reader, kbps float64) io.Reader {
	
	kb := 1.0
	sleepDuration := time.Duration(kb / kbps * float64(time.Second))
	return &readerWrapper{reader: reader, sleepDuration: sleepDuration}
}

type readerWrapper struct {
	reader        io.Reader
	sleepDuration time.Duration
}

// Read implements the io.Reader interface with throttling.
func (w *readerWrapper) Read(p []byte) (n int, err error) {
	start := time.Now()
		
	readSize := 125
	if len(p) < readSize {
		readSize = len(p)
	}
	// Read at most 125 bytes (1 kilobit)
	n, err = w.reader.Read(p[:readSize])


	elapsedTime := time.Since(start)
	if elapsedTime < w.sleepDuration {
		time.Sleep(w.sleepDuration - elapsedTime)
	}

	return n, err
}
