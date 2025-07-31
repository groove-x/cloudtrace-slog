package cloudtrace

import "context"

type contextKey string

const traceKey contextKey = "trace"

type Trace struct {
	TraceID string
	SpanID  string
	Sampled bool
}

func traceFromContext(ctx context.Context) *Trace {
	trace := ctx.Value(traceKey)
	if trace == nil {
		return nil
	}
	traceInfo, ok := trace.(*Trace)
	if !ok {
		return nil
	}
	return traceInfo
}

func withTraceContext(ctx context.Context, traceID string, spanID string, sampled bool) context.Context {
	if traceID == "" {
		return ctx
	}
	return context.WithValue(ctx, traceKey, &Trace{
		TraceID: traceID,
		SpanID:  spanID,
		Sampled: sampled,
	})
}
