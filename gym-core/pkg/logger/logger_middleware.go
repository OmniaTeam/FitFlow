package logger

import (
	"context"
	"log/slog"
	"strconv"
)

type (
	mapLogger map[string]string
	keyType   int
)

const keyMap = keyType(0)

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(keyMap).(mapLogger); ok {

		for key, value := range c {
			rec.Add(key, value)
		}

	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	key := "trace_id"
	if c, ok := ctx.Value(keyMap).(mapLogger); ok {
		c[key] = traceID
		return context.WithValue(ctx, keyMap, c)
	}
	return context.WithValue(ctx, keyMap, mapLogger{key: traceID})
}

func WithUserID(ctx context.Context, userID int) context.Context {
	key := "user_id"
	if c, ok := ctx.Value(keyMap).(mapLogger); ok {
		c[key] = strconv.Itoa(userID)
		return context.WithValue(ctx, keyMap, c)
	}
	return context.WithValue(ctx, keyMap, mapLogger{key: strconv.Itoa(userID)})
}

func WithKeyValue(ctx context.Context, key, value string) context.Context {
	if c, ok := ctx.Value(keyMap).(mapLogger); ok {
		c[key] = value
		return context.WithValue(ctx, keyMap, c)
	}
	return context.WithValue(ctx, keyMap, mapLogger{key: value})
}
