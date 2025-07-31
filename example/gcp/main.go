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
