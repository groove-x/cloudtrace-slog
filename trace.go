package cloudtrace

import (
	"net/http"
	"strconv"
	"strings"
)

func extractTraceInfo(r *http.Request) (string, string, bool) {
	// W3C Trace Context (traceparent)
	// https://www.w3.org/TR/trace-context/#trace-context-http-headers-format
	// Format: VERSION-TRACE_ID-PARENT_ID-TRACE_FLAGS
	traceparent := r.Header.Get("traceparent")
	if traceparent != "" {
		parts := strings.Split(traceparent, "-")
		if len(parts) >= 4 && parts[0] == "00" {
			traceID := parts[1]
			spanID := parts[2]
			if traceFlags, err := strconv.ParseUint(parts[3], 16, 8); err == nil {
				sampled := (traceFlags & 0x01) == 0x01
				return traceID, spanID, sampled
			}
			return traceID, spanID, false
		}
	}

	// Google Cloud Trace Context
	// https://cloud.google.com/trace/docs/trace-context?hl=ja
	// X-Cloud-Trace-Context: TRACE_ID/SPAN_ID;o=OPTIONS
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	if traceHeader != "" {
		parts := strings.Split(traceHeader, "/")
		if len(parts) >= 2 && parts[0] != "" {
			traceID := parts[0]
			spanInfo := parts[1]

			spanParts := strings.Split(spanInfo, ";")
			spanID := spanParts[0]

			sampled := false
			if len(spanParts) > 1 {
				for _, option := range spanParts[1:] {
					if option == "o=1" {
						sampled = true
						break
					}
				}
			}
			return traceID, spanID, sampled
		}
	}

	return "", "", false
}
