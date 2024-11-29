package slogctx

import (
	"context"
	"log/slog"
)

// An Extractor extracts [slog.Attr] values from a [context.Context].
type Extractor interface {
	Extract(ctx context.Context) []slog.Attr
}

// Ensure that [ExtractorFunc] implements [Extractor].
var _ Extractor = ExtractorFunc(func(ctx context.Context) []slog.Attr { return nil })

// ExtractorFunc allows a function to be used as an [Extractor].
type ExtractorFunc func(ctx context.Context) []slog.Attr

// Extract calls the underlying function to implement the [Extractor] interface.
func (f ExtractorFunc) Extract(ctx context.Context) []slog.Attr {
	return f(ctx)
}

// newCtxExtractor creates an [ExtractorFunc] that uses one of the allowed keys to extract
// [slog.Attr] values from the given [context.Context].
func newCtxExtractor[K ctxKeyWithAttrs | ctxKeyWithRootAttrs](key K) ExtractorFunc {
	return func(ctx context.Context) []slog.Attr {
		if attrs, ok := ctx.Value(key).([]slog.Attr); ok {
			return attrs
		}

		return nil
	}
}
