package log

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rautaruukkipalich/prettyslog"
)

type (
	ctxLogger struct{}
	ctxKey    struct{}
	ctxTrace  struct{}
)

const (
	attrKey  = "attrs"
	traceKey = "trace"
)

var (
	globalLogger *slog.Logger
)

func init() {
	globalLogger = prettyslog.NewPrettyLogger("\t")
}

func CreateGlobalLogger(logger *slog.Logger) {
	globalLogger = logger
}

func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger := ctx.Value(ctxLogger{})
	attrs := getFields(ctx)
	trace := getTrace(ctx)

	if logger == nil {
		logger = globalLogger
	}

	lenAttrs := len(attrs)
	lenTrace := len(trace)

	if lenAttrs == 0 && lenTrace == 0 {
		return logger.(*slog.Logger)
	}

	if lenAttrs == 0 && lenTrace > 0 {
		return logger.(*slog.Logger).With(slog.Any(traceKey, trace))
	}

	if lenAttrs > 0 && lenTrace == 0 {
		return logger.(*slog.Logger).With(slog.Any(attrKey, attrs))
	}

	return logger.(*slog.Logger).With(
		slog.Any(attrKey, attrs),
		slog.Any(traceKey, trace), 
	)
}

func AddAttr(ctx context.Context, key string, value any) context.Context {
	attrs := getFields(ctx)
	attrs = append(attrs, slog.Any(key, value))
	return context.WithValue(ctx, ctxKey{}, attrs)
}

func AddGroup(ctx context.Context, group slog.Attr) context.Context {
	attrs := getFields(ctx)
	attrs = append(attrs, group)
	return context.WithValue(ctx, ctxKey{}, attrs)
}

func AddTrace(ctx context.Context, trace slog.Attr) context.Context {
	traces := getTrace(ctx)
	traces = append(traces, slog.Attr{
		Key:   fmt.Sprintf("%d: %s", len(traces)+1, trace.Key),
		Value: trace.Value,
	})
	return context.WithValue(ctx, ctxTrace{}, traces)
}

func getFields(ctx context.Context) []slog.Attr {
	attrs := ctx.Value(ctxKey{})
	if attrs == nil {
		return []slog.Attr{}
	}
	return attrs.([]slog.Attr)
}

func getTrace(ctx context.Context) []slog.Attr {
	trace := ctx.Value(ctxTrace{})
	if trace == nil {
		return []slog.Attr{}
	}
	return trace.([]slog.Attr)
}
