package benchmarks

import (
	"io"
	"log/slog"

	"github.com/nickbryan/slogutil"
	"github.com/nickbryan/slogutil/slogctx"
	"github.com/nickbryan/slogutil/slogmem"
)

func newSlogUtilInMem(fields ...slog.Attr) *slog.Logger {
	return slog.New(slogmem.NewHandler(slog.LevelDebug).WithAttrs(fields))
}

func newDisabledSlogUtilInMem(fields ...slog.Attr) *slog.Logger {
	return slog.New(slogmem.NewHandler(slog.LevelError).WithAttrs(fields))
}

func newSlogUtilCtx(fields ...slog.Attr) *slog.Logger {
	return slog.New(slogctx.NewHandler(slog.NewJSONHandler(io.Discard, nil).WithAttrs(fields)))
}

func newDisabledSlogUtilCtx(fields ...slog.Attr) *slog.Logger {
	return slog.New(slogctx.NewHandler(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}).WithAttrs(fields)))
}

func newSlogUtilJSONLogger(fields ...slog.Attr) *slog.Logger {
	logger := slogutil.NewJSONLogger(slogutil.WithWriter(io.Discard))
	logger.Handler().WithAttrs(fields)
	return logger
}

func newDisabledSlogUtilJSONLogger(fields ...slog.Attr) *slog.Logger {
	logger := slogutil.NewJSONLogger(slogutil.WithWriter(io.Discard), slogutil.WithLevel(slog.LevelError))
	logger.Handler().WithAttrs(fields)
	return logger
}
