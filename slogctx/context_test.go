package slogctx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/nickbryan/slogutil/slogctx"
	"github.com/nickbryan/slogutil/slogmem"
)

func TestWithAttrs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ctx  context.Context
		log  func(ctx context.Context, logger *slog.Logger)
		want slogmem.RecordQuery
	}{
		"appending attrs to a log entry with no additional log attrs adds the attrs": {
			ctx: slogctx.WithAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1")},
			},
		},
		"appending attrs to a log entry with additional log attrs appends the attrs": {
			ctx: slogctx.WithAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message", slog.Int("e1", 123))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"e1": slog.IntValue(123), "p1": slog.StringValue("v1")},
			},
		},
		"appending grouped attrs to a log entry appends the grouped attrs": {
			ctx: slogctx.WithRootAttrs(context.Background(), slog.Group("g1", slog.String("p1", "v1"), slog.String("p2", "v2"))),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message", slog.Int("e1", 123))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"e1": slog.IntValue(123), "g1.p1": slog.StringValue("v1"), "g1.p2": slog.StringValue("v2")},
			},
		},
		"appending attrs to a log entry that contains groups appends attrs": {
			ctx: slogctx.WithAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message", slog.Group("g1", slog.Int("e1", 123)))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"g1.e1": slog.IntValue(123), "p1": slog.StringValue("v1")},
			},
		},
		"appending attrs to a log entry that is nested in a group appends the attr to the current group": {
			ctx: slogctx.WithAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1")},
			},
		},
		"appending attrs to a nil ctx returns a ctx with the given attrs": {
			ctx: slogctx.WithAttrs(nil, slog.String("p1", "v1")), //nolint:staticcheck // Staticcheck warns on the use of nil ctx.
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").InfoContext(ctx, "Test message", slog.Int("e1", 123))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"g1.e1": slog.IntValue(123), "g1.p1": slog.StringValue("v1")},
			},
		},
		"appending attrs to a ctx with existing attrs adds the attrs": {
			ctx: slogctx.WithAttrs(slogctx.WithAttrs(context.Background(), slog.String("p1", "v1")), slog.String("p2", "v2")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1"), "p2": slog.StringValue("v2")},
			},
		},
		"appending duplicate attrs to a ctx with existing attrs adds the attrs": {
			ctx: slogctx.WithAttrs(slogctx.WithAttrs(context.Background(), slog.String("p1", "v1")), slog.String("p1", "v2")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1"), "p1#01": slog.StringValue("v2")},
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			handler := slogmem.NewHandler(slog.LevelDebug)
			logger := slog.New(slogctx.NewHandler(handler))

			testCase.log(testCase.ctx, logger)

			records := handler.Records()
			if ok, diff := records.ContainsExact(testCase.want); !ok {
				t.Errorf("expected logged records to contain: %+v, got: %s", testCase.want, diff)
			}
		})
	}
}

func TestWithRootAttrs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ctx  context.Context
		log  func(ctx context.Context, logger *slog.Logger)
		want slogmem.RecordQuery
	}{
		"prepending attrs to a log entry with no additional log attrs adds the attrs": {
			ctx: slogctx.WithRootAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1")},
			},
		},
		"prepending attrs to a log entry with additional log attrs prepends the attrs": {
			ctx: slogctx.WithRootAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message", slog.Int("e1", 123))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1"), "e1": slog.IntValue(123)},
			},
		},
		"prepending grouped attrs to a log entry prepends the grouped attrs": {
			ctx: slogctx.WithRootAttrs(context.Background(), slog.Group("g1", slog.String("p1", "v1"), slog.String("p2", "v2"))),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message", slog.Int("e1", 123))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"g1.p1": slog.StringValue("v1"), "g1.p2": slog.StringValue("v2"), "e1": slog.IntValue(123)},
			},
		},
		"prepending attrs to a log entry that contains groups prepends attrs": {
			ctx: slogctx.WithRootAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message", slog.Group("g1", slog.Int("e1", 123)))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1"), "g1.e1": slog.IntValue(123)},
			},
		},
		"prepending attrs to a log entry that is nested in a group prepends the attrs to the root": {
			ctx: slogctx.WithRootAttrs(context.Background(), slog.String("p1", "v1")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").InfoContext(ctx, "Test message", slog.Int("e1", 123))
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1"), "g1.e1": slog.IntValue(123)},
			},
		},
		"prepending attrs to a nil ctx returns a ctx with the given attrs": {
			ctx: slogctx.WithRootAttrs(nil, slog.String("p1", "v1")), //nolint:staticcheck // Staticcheck warns on the use of nil ctx.
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1")},
			},
		},
		"prepending attrs to a ctx with existing attrs adds the attrs": {
			ctx: slogctx.WithRootAttrs(slogctx.WithRootAttrs(context.Background(), slog.String("p1", "v1")), slog.String("p2", "v2")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1"), "p2": slog.StringValue("v2")},
			},
		},
		"prepending duplicate attrs to a ctx with existing attrs adds the attrs": {
			ctx: slogctx.WithRootAttrs(slogctx.WithRootAttrs(context.Background(), slog.String("p1", "v1")), slog.String("p1", "v2")),
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "Test message")
			},
			want: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "Test message",
				Attrs:   map[string]slog.Value{"p1": slog.StringValue("v1"), "p1#01": slog.StringValue("v2")},
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			handler := slogmem.NewHandler(slog.LevelDebug)
			logger := slog.New(slogctx.NewHandler(handler))

			testCase.log(testCase.ctx, logger)

			records := handler.Records()
			if ok, diff := records.ContainsExact(testCase.want); !ok {
				t.Errorf("expected logged records to contain: %+v, got: %s", testCase.want, diff)
			}
		})
	}
}
