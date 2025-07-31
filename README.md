# Cloud Trace and Cloud Logging integration package for slog

A simple Go package that provides integration with Google Cloud Logging and supports structured logging and distributed tracing.

Note: This package does not require the Google Cloud Logging API 
([cloud.google.com/go/logging](https://pkg.go.dev/cloud.google.com/go/logging) package) 
because it uses standard error output for logging.

## Usage

This example can be deployed on Google Cloud Run, App Engine, or any other environment that supports Google Cloud Trace and Logging.

```go
package main

import (
	"log/slog"
	"net/http"

	cloudtrace "github.com/groove-x/cloudtrace-slog"
)

func main() {
	// Setup slog
	slog.SetDefault(slog.New(cloudtrace.NewCloudLoggingHandler()))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// logging with trace context
		slog.InfoContext(ctx, "hello world")

		w.Write([]byte("Hello, World!"))
	})

	// Wrap your handler with trace context middleware
	handler := cloudtrace.WithCloudTraceContextMiddleware(mux)

	http.ListenAndServe(":8080", handler)
}
```

## How It Works

1. The middleware extracts trace IDs from incoming HTTP requests
2. Trace context is stored in the request context
3. The logging handler automatically includes trace information in log entries
4. Output the logs to standard error, which Google Cloud Logging captures
5. This enables correlation between logs and traces in Google Cloud Console

## Requirements

- Go 1.21+
- Running on Google Cloud Platform (for trace integration)
- `cloud.google.com/go/compute/metadata` package

## References

- [Google Cloud Logging Documentation](https://cloud.google.com/logging/docs)
- [Google Cloud Trace Documentation](https://cloud.google.com/trace/docs)
- [Go Metadata Package](https://pkg.go.dev/cloud.google.com/go/compute/metadata)
- [Go slog Package](https://pkg.go.dev/log/slog)
- [Go HTTP Package](https://pkg.go.dev/net/http)
- [Go context Package](https://pkg.go.dev/context)
