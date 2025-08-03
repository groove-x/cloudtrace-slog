package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	cloudtrace "github.com/groove-x/cloudtrace-slog"
)

func main() {
	log.Print("starting server...")

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

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

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
