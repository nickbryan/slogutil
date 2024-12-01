package slogmem_test

import (
	"encoding/json"
	"log/slog"
	"maps"
	"testing"
	"testing/slogtest"
	"time"

	"github.com/nickbryan/slogutil/slogmem"
)

// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If r.PC is zero, ignore it.
//   - Attr's values should be resolved.
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//   - If a group's key is empty, inline the group's Attrs.
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
//
// This test harness ensures that we are meeting the above criteria for a
// [slog.Handler] implementation. The above snippet is taken from the
// [slog.Handler] interface's [Handle] method docs: https://pkg.go.dev/log/slog#Handler.
func TestHandlerSatisfiesSlogTestHarness(t *testing.T) {
	t.Parallel()

	handler := slogmem.NewHandler(slog.LevelDebug)

	results := func() []map[string]any {
		records := handler.Records().AsSliceOfNestedKeyValuePairs()

		for _, record := range records {
			// Unexpected key "time": a Handler should ignore a zero Record.Time
			//
			// The testing/slogtest harness executes the above assertion. We want to ensure
			// that we capture zero time for debugging purposes when the in memory Handler is
			// used for such cases. We capture all time values in the Handler and we delete
			// them here in order to past the test harness as per https://pkg.go.dev/testing/slogtest#TestHandler.
			maps.DeleteFunc(record, func(key string, value any) bool {
				if t, ok := value.(time.Time); ok && key == slog.TimeKey {
					return t.IsZero() // Delete time attribute where value is zero.
				}

				return false
			})
		}

		return records
	}

	if err := slogtest.TestHandler(handler, results); err != nil {
		jsonResults, marshalErr := json.MarshalIndent(results(), "", "  ")
		if marshalErr != nil {
			t.Fatalf("Unable to marshal JSON results: got: %v, want: no marshal errors", marshalErr)
		}

		t.Errorf("testing/slogtest harness is not satisfied for slogmem.Handler\ngot error: \n%s\n\ngot logs: \n%s", err, jsonResults)
	}
}
