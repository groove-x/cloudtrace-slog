package cloudtrace

import (
	"net/http"
	"testing"
)

func TestExtractTraceInfo(t *testing.T) {
	tests := []struct {
		name            string
		headers         map[string]string
		expectedTraceID string
		expectedSpanID  string
		expectedSampled bool
	}{
		// W3C Trace Context tests
		{
			name: "Valid W3C traceparent with sampling",
			headers: map[string]string{
				"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			},
			expectedTraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
			expectedSpanID:  "00f067aa0ba902b7",
			expectedSampled: true,
		},
		{
			name: "Valid W3C traceparent without sampling",
			headers: map[string]string{
				"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			},
			expectedTraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
			expectedSpanID:  "00f067aa0ba902b7",
			expectedSampled: false,
		},
		{
			name: "W3C traceparent with invalid trace flags",
			headers: map[string]string{
				"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-xx",
			},
			expectedTraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
			expectedSpanID:  "00f067aa0ba902b7",
			expectedSampled: false,
		},
		{
			name: "Invalid W3C traceparent - wrong version",
			headers: map[string]string{
				"traceparent": "01-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			},
			expectedTraceID: "",
			expectedSpanID:  "",
			expectedSampled: false,
		},
		{
			name: "Invalid W3C traceparent - insufficient parts",
			headers: map[string]string{
				"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7",
			},
			expectedTraceID: "",
			expectedSpanID:  "",
			expectedSampled: false,
		},
		{
			name: "Invalid W3C traceparent - empty",
			headers: map[string]string{
				"traceparent": "",
			},
			expectedTraceID: "",
			expectedSpanID:  "",
			expectedSampled: false,
		},

		// Google Cloud Trace Context tests
		{
			name: "Valid GCP trace context with sampling",
			headers: map[string]string{
				"X-Cloud-Trace-Context": "105445aa7843bc8bf206b120001000/1;o=1",
			},
			expectedTraceID: "105445aa7843bc8bf206b120001000",
			expectedSpanID:  "1",
			expectedSampled: true,
		},
		{
			name: "Valid GCP trace context without sampling",
			headers: map[string]string{
				"X-Cloud-Trace-Context": "105445aa7843bc8bf206b120001000/1;o=0",
			},
			expectedTraceID: "105445aa7843bc8bf206b120001000",
			expectedSpanID:  "1",
			expectedSampled: false,
		},
		{
			name: "Valid GCP trace context without options",
			headers: map[string]string{
				"X-Cloud-Trace-Context": "105445aa7843bc8bf206b120001000/1",
			},
			expectedTraceID: "105445aa7843bc8bf206b120001000",
			expectedSpanID:  "1",
			expectedSampled: false,
		},
		{
			name: "Invalid GCP trace context - missing span ID",
			headers: map[string]string{
				"X-Cloud-Trace-Context": "105445aa7843bc8bf206b120001000",
			},
			expectedTraceID: "",
			expectedSpanID:  "",
			expectedSampled: false,
		},
		{
			name: "Invalid GCP trace context - empty trace ID",
			headers: map[string]string{
				"X-Cloud-Trace-Context": "/1;o=1",
			},
			expectedTraceID: "",
			expectedSpanID:  "",
			expectedSampled: false,
		},
		{
			name: "Invalid GCP trace context - empty",
			headers: map[string]string{
				"X-Cloud-Trace-Context": "",
			},
			expectedTraceID: "",
			expectedSpanID:  "",
			expectedSampled: false,
		},

		// Priority tests (W3C takes precedence over GCP)
		{
			name: "W3C traceparent takes precedence over GCP",
			headers: map[string]string{
				"traceparent":           "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
				"X-Cloud-Trace-Context": "105445aa7843bc8bf206b120001000/1;o=0",
			},
			expectedTraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
			expectedSpanID:  "00f067aa0ba902b7",
			expectedSampled: true,
		},

		// No trace headers
		{
			name:            "No trace headers",
			headers:         map[string]string{},
			expectedTraceID: "",
			expectedSpanID:  "",
			expectedSampled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Set headers
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			traceID, spanID, sampled := extractTraceInfo(req)

			if traceID != tt.expectedTraceID {
				t.Errorf("Expected traceID %q, got %q", tt.expectedTraceID, traceID)
			}
			if spanID != tt.expectedSpanID {
				t.Errorf("Expected spanID %q, got %q", tt.expectedSpanID, spanID)
			}
			if sampled != tt.expectedSampled {
				t.Errorf("Expected sampled %t, got %t", tt.expectedSampled, sampled)
			}
		})
	}
}
