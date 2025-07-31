package cloudtrace

import (
	"context"
	"log/slog"
	"os"
)

type CloudLoggingHandler struct{ handler slog.Handler }

func NewCloudLoggingHandler() slog.Handler {
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.MessageKey:
				a.Key = "message"
			case slog.LevelKey:
				a.Key = "severity"
			}
			return a
		},
	})

	return &CloudLoggingHandler{handler: h}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CloudLoggingHandler) Handle(ctx context.Context, r slog.Record) error {
	trace := traceFromContext(ctx)
	if trace != nil && trace.TraceID != "" {
		r = r.Clone()
		// https://cloud.google.com/trace/docs/trace-log-integration
		r.Add("logging.googleapis.com/trace", slog.StringValue(trace.TraceID))
		r.Add("logging.googleapis.com/trace_sampled", slog.BoolValue(trace.Sampled))
		if trace.SpanID != "" {
			r.Add("logging.googleapis.com/spanId", slog.StringValue(trace.SpanID))
		}
	}
	return h.handler.Handle(ctx, r)
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}
