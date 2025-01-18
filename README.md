# go-traffic-throttle

`go-traffic-throttle` is a lightweight Go package for limiting upload speeds in HTTP requests or other data streaming scenarios. It provides an `io.Reader` wrapper that throttles bandwidth to a specified rate, ensuring smooth and controlled data transmission.

## Features

- Bandwidth throttling in kilobits per second (kbps).
- Simple integration with HTTP requests or any `io.Reader`.
- Customizable speed limits for controlled data flow.

## Installation

```bash
go get github.com/OpenCodeNest/go-traffic-throttle
```

## Usage

Here's how you can use `go-traffic-throttle` to limit the upload speed in an HTTP POST request:

```go
package main

import (
	"bytes"
	"net/http"
	"time"

	"github.com/OpenCodeNest/go-traffic-throttle/throttle"
)

func main() {
	// Configuration
	const SERVER_URL = "https://example.com/upload"
	const BANDWIDTH_SPEED_KBPS = 64.0 // Bandwidth limit in kilobits per second

	// Sample body to upload
	bodyData := []byte("This is a sample payload for the POST request.")
	body := bytes.NewReader(bodyData)

	// Wrap the body reader with the Throttle function
	throttledBody := throttle.Throttle(body, BANDWIDTH_SPEED_KBPS)

	// Create the HTTP request with the throttled reader
	req, err := http.NewRequest("POST", SERVER_URL, throttledBody)
	if err != nil {
		panic(err)
	}

	// Add necessary headers
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print response status
	println("Response status:", resp.Status)
}
```

## API

### `Throttle`

```go
func Throttle(reader io.Reader, kbps float64) io.Reader
```

Wraps an `io.Reader` to limit its bandwidth to the specified `kbps` (kilobits per second).

**Parameters:**
- `reader`: The original `io.Reader` to be wrapped.
- `kbps`: The bandwidth limit in kilobits per second.

**Returns:**
- A wrapped `io.Reader` that enforces the bandwidth limit.

## Example Use Cases

- Limiting upload speeds in HTTP file uploads.
- Throttling data streams in distributed systems.
- Simulating low-bandwidth scenarios for testing.

## Contributing

Contributions are welcome! If you have ideas, improvements, or bug fixes, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.