package slogctx

import (
	"context"
	"log/slog"
	"slices"
)

type (
	ctxKeyWithAttrs     struct{}
	ctxKeyWithRootAttrs struct{}
)

// WithAttrs will add attrs to the [context.Context] so that they can be
// appended to the log attributes a the log is written with the given
// [context.Context]. This is helpful when you want to ensure that the attrs
// are placed at the end of the current group.
//
// Making subsequent calls to this on the same [context.Context] will result in
// attrs being appended to the set. This is safe to do.
func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return addToContext(ctx, ctxKeyWithAttrs{}, attrs)
}

// WithRootAttrs will add attrs to the [context.Context] so that they can be
// prepended to the log attributes when the log is written with the given
// [context.Context]. This is helpful when you want to ensure that the attrs are
// placed at the root of the log and not within the current group.
//
// Making subsequent calls to this on the same [context.Context] will result in
// attrs being appended to the set. This is safe to do.
func WithRootAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return addToContext(ctx, ctxKeyWithRootAttrs{}, attrs)
}

func addToContext[K ctxKeyWithAttrs | ctxKeyWithRootAttrs](ctx context.Context, key K, attrs []slog.Attr) context.Context {
	if existingAttrs, ok := ctx.Value(key).([]slog.Attr); ok {
		return context.WithValue(ctx, key, append(slices.Clip(existingAttrs), slices.Clip(attrs)...))
	}

	return context.WithValue(ctx, key, attrs)
}
