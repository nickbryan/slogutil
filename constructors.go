package slogutil

import (
	"log/slog"

	"github.com/nickbryan/slogutil/slogctx"
	"github.com/nickbryan/slogutil/slogmem"
)

// NewJSONLogger creates a new [slog.Logger] configured with a
// [slogctx.Handler] which wraps a [slog.JSONHandler].
func NewJSONLogger(options ...Option) *slog.Logger {
	opts := mapOptionsToDefaults(options)

	jsonHandler := slog.NewJSONHandler(opts.writer, &slog.HandlerOptions{
		AddSource: opts.addSource,
		Level:     opts.level,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if opts.now != nil && attr.Key == slog.TimeKey {
				attr.Value = slog.TimeValue(opts.now())
			}

			return attr
		},
	})

	return slog.New(slogctx.NewHandler(jsonHandler))
}

// NewInMemoryLogger creates a new [slog.Logger] configured with a
// [slogmem.Handler] to capture logged records in-memory for testing.
//
// A [slogmem.LoggedRecords] will also be returned containing the
// records created by the returned [slog.Logger].
func NewInMemoryLogger(level slog.Leveler) (*slog.Logger, *slogmem.LoggedRecords) {
	handler := slogmem.NewHandler(level)

	return slog.New(slogctx.NewHandler(handler)), handler.Records()
}
