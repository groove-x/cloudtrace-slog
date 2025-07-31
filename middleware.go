package cloudtrace

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/compute/metadata"
)

func WithCloudTraceContextMiddleware(h http.Handler) http.Handler {
	if !metadata.OnGCE() {
		return h
	}

	projectID, err := metadata.ProjectIDWithContext(context.Background())
	if err != nil {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, spanID, sampled := extractTraceInfo(r)
		if traceID == "" {
			h.ServeHTTP(w, r)
			return
		}

		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
		ctx := withTraceContext(r.Context(), trace, spanID, sampled)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
