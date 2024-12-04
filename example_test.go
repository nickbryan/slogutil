package slogutil_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/nickbryan/slogutil"
	"github.com/nickbryan/slogutil/slogctx"
	"github.com/nickbryan/slogutil/slogmem"
)

func constantTimeFactory() time.Time {
	t, err := time.Parse(time.RFC3339, "2024-03-05T12:00:00Z")
	if err != nil {
		panic("unable to parse time value: " + err.Error())
	}

	return t
}

func ExampleNewJSONLogger() {
	ctx := context.Background()

	logger := slogutil.NewJSONLogger(
		slogutil.WithLevel(slog.LevelInfo),
		slogutil.WithWriter(os.Stdout),
		slogutil.WithSourceAdded(false),
		slogutil.WithTimeFactory(constantTimeFactory),
	)
	logger = logger.With(slog.Int("my_root_attribute", 123))
	logger = logger.WithGroup("my_group")

	logger.DebugContext(ctx, "Debug log message") // Not logged due to the level set on the logger.
	logger.InfoContext(ctx, "Info log message", slog.String("my_grouped_attribute", "my_value"))

	// Output:
	// {"time":"2024-03-05T12:00:00Z","level":"INFO","msg":"Info log message","my_root_attribute":123,"my_group":{"my_grouped_attribute":"my_value"}}
}

func ExampleNewJSONLogger_context() {
	ctx := slogctx.WithRootAttrs(context.Background(), slog.String("prepend_attribute", "prepend_value"))
	ctx = slogctx.WithAttrs(ctx, slog.String("append_attribute", "append_value"))

	logger := slogutil.NewJSONLogger(
		slogutil.WithLevel(slog.LevelInfo),
		slogutil.WithWriter(os.Stdout),
		slogutil.WithSourceAdded(false),
		slogutil.WithTimeFactory(constantTimeFactory),
	)
	logger = logger.With(slog.Int("my_root_attribute", 123))
	logger = logger.WithGroup("my_group")

	logger.DebugContext(ctx, "Debug log message") // Not logged due to the level set on the logger.
	logger.InfoContext(ctx, "Info log message", slog.String("my_grouped_attribute", "my_value"))

	// Output:
	// {"time":"2024-03-05T12:00:00Z","level":"INFO","msg":"Info log message","prepend_attribute":"prepend_value","my_root_attribute":123,"my_group":{"my_grouped_attribute":"my_value","append_attribute":"append_value"}}
}

func ExampleNewInMemoryLogger() {
	ctx := context.Background()

	ctx = slogctx.WithRootAttrs(ctx, slog.String("prepend_attribute", "prepend_value"))
	ctx = slogctx.WithAttrs(ctx, slog.String("append_attribute", "append_value"))

	logger, logs := slogutil.NewInMemoryLogger(slog.LevelInfo)
	logger = logger.With(slog.Int("my_root_attribute", 123))
	logger = logger.WithGroup("my_group")

	logger.DebugContext(ctx, "Debug log message") // Not logged due to the level set on the logger.
	logger.InfoContext(ctx, "Info log message", slog.String("my_grouped_attribute", "my_value"))

	if ok, diff := logs.Contains(slogmem.RecordQuery{
		Level:   slog.LevelInfo,
		Message: "Info log message",
		Attrs: map[string]slog.Value{
			"prepend_attribute":             slog.StringValue("prepend_value"),
			"my_root_attribute":             slog.IntValue(123),
			"my_group.my_grouped_attribute": slog.StringValue("my_value"),
			"my_group.append_attribute":     slog.StringValue("append_value"),
		},
	}); !ok {
		fmt.Print(diff)
	} else {
		fmt.Print("Record contains query")
	}

	// Output: Record contains query
}
