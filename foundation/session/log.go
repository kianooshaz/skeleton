package session

import (
	"context"
	"log/slog"
)

type logKet struct{}

func GetLogAttributes(ctx context.Context) []slog.Attr {
	v := ctx.Value(logKet{})
	if v == nil {
		return make([]slog.Attr, 0)
	}

	return v.([]slog.Attr)
}

func SetLogAttributes(ctx context.Context, attrs ...slog.Attr) context.Context {
	v := ctx.Value(logKet{})
	if v == nil {
		return context.WithValue(ctx, logKet{}, attrs)
	}

	return context.WithValue(ctx, logKet{}, append(v.([]slog.Attr), attrs...))
}
