# slogutil
Package `slogutil` provides enhanced logging capabilities for the standard library `log/slog` package, focusing on
context integration and testability.

<div align="center">

[![Test](https://github.com/nickbryan/slogutil/actions/workflows/test.yml/badge.svg)](https://github.com/nickbryan/slogutil/actions)
[![Coverage](https://raw.githubusercontent.com/nickbryan/slogutil/badges/.badges/main/coverage.svg)](https://github.com/nickbryan/slogutil/actions)
[![Go Report Card](https://goreportcard.com/badge/nickbryan/slogutil)](https://goreportcard.com/report/nickbryan/slogutil)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/nickbryan/slogutil/blob/master/LICENSE)

</div>

## Features
* **Context-Aware Logging:**  Enriches log records with contextual information from `context.Context`:
    * The `slogctx` sub-package provides a handler (`slogctx.Handler`) and an `Extractor` API to extract values from the context.
    * Supports adding attributes to the root of the log context or appending them within the current log group.

* **Testability:** Enables easy testing of log output:
    * Provides an in-memory handler (`slogmem`) to capture log records during tests, allowing for assertions and verification.

* **Attribute Consistency:** Provides consistent handling of log attributes:
    * Deduplicates attributes with the same respecting groups. For example: `duplicate`, `duplicate#01`, `duplcate#02`.

## Quick Start
```go
package main

import (
	"context"
	"log/slog"

	"github.com/nickbryan/slogutil"
	"github.com/nickbryan/slogutil/slogctx"
)

func main() {
	ctx := slogctx.WithAttrs(context.Background(), slog.String("my_appended_attribute", "my_appended_value"))
	ctx = slogctx.WithRootAttrs(ctx, slog.String("my_root_attribute", "my_root_value"))

	logger := slogutil.NewJSONLogger().WithGroup("my_group")
	logger.InfoContext(ctx, "Info log message", slog.String("my_attribute", "my_value"))

	// Output:
	// {"time":"2024-11-29T10:27:41.460372Z","level":"INFO","source":{"function":"main.main","file":main.go","line":16},"msg":"Info log message","my_root_attribute":"my_root_value","my_group":{"my_attribute":"my_value","my_appended_attribute":"my_appended_value"}}
}
```
## Testing
`slogutil` provides an in-memory handler (`slogmem`) to capture log records, enabling you to test log output effectively:

```go
func TestThing(t *testing.T) {
	ctx := slogctx.WithPrependAttrs(context.Background(), slog.String("prepend_attribute", "prepend_value"))
	ctx = slogctx.WithAppendAttrs(ctx, slog.String("append_attribute", "append_value"))

	logger, logs := slogutil.NewInMemoryLogger(slog.LevelDebug)
	logger = logger.With(slog.Int("my_root_attribute", 123))
	logger = logger.WithGroup("my_group")

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
		t.Errorf("expected log not written, got: %s", diff)
	}
}
```
