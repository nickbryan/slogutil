package slogmem_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/nickbryan/slogutil/slogmem"
)

type logValuerStubError struct {
	I int
	S string
}

func (e *logValuerStubError) Error() string { return "some error from logValuerStubError" }
func (e *logValuerStubError) LogValue() slog.Value {
	return slog.GroupValue(slog.Int("i", e.I), slog.String("s", e.S))
}

func TestLoggedRecordsAsSliceOfNestedKeyValuePairs(t *testing.T) {
	t.Parallel()

	fixedNow := time.Date(2024, 5, 28, 1, 0, 0, 0, time.UTC)

	testCases := map[string]struct {
		records []slogmem.LoggedRecord
		want    []map[string]any
	}{
		"calling AsSliceOfNestedKeyValuePairs on an empty LoggedRecords initialized with nil returns an empty slice": {
			records: nil,
			want:    []map[string]any{},
		},
		"calling AsSliceOfNestedKeyValuePairs on an empty LoggedRecords initialized with an empty slice returns an empty slice": {
			records: []slogmem.LoggedRecord{},
			want:    []map[string]any{},
		},
		"calling AsSliceOfNestedKeyValuePairs on a list of LoggedRecords returns the flattened version of them records": {
			records: []slogmem.LoggedRecord{
				{
					Time:    fixedNow,
					Level:   slog.LevelDebug,
					Message: "first debug log entry",
					Attrs:   nil,
				},
				{
					Time:    fixedNow.Add(time.Hour),
					Level:   slog.LevelDebug,
					Message: "second debug log entry",
					Attrs:   []slog.Attr{slog.String("r2ka", "r2va"), slog.String("r2kb", "r2vb")},
				},
			},
			want: []map[string]any{
				{slog.TimeKey: fixedNow, slog.LevelKey: slog.LevelDebug, slog.MessageKey: "first debug log entry"},
				{slog.TimeKey: fixedNow.Add(time.Hour), slog.LevelKey: slog.LevelDebug, slog.MessageKey: "second debug log entry", "r2ka": "r2va", "r2kb": "r2vb"},
			},
		},
		"calling AsSliceOfNestedKeyValuePairs with a LoggedRecord that has grouped attributes nests the attributes under the group": {
			records: []slogmem.LoggedRecord{
				{
					Time:    fixedNow,
					Level:   slog.LevelInfo,
					Message: "first info log entry",
					Attrs:   []slog.Attr{slog.Group("ga", slog.String("gaka", "gava"), slog.String("gakb", "gavb")), slog.String("ka", "va"), slog.Group("gb", slog.String("gbka", "gbva"))},
				},
			},
			want: []map[string]any{
				{slog.TimeKey: fixedNow, slog.LevelKey: slog.LevelInfo, slog.MessageKey: "first info log entry", "ga": map[string]any{"gaka": "gava", "gakb": "gavb"}, "ka": "va", "gb": map[string]any{"gbka": "gbva"}},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := slogmem.NewLoggedRecords(tc.records).AsSliceOfNestedKeyValuePairs()

			if !cmp.Equal(tc.want, got) {
				t.Errorf("slogmem.NewLoggedRecords(%+v).AsSliceOfNestedKeyValuePairs():\n got %+v\n want %+v\n diff: %s", tc.records, got, tc.want, cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestLoggedRecordsContains(t *testing.T) {
	t.Parallel()

	t.Run("panics when a slog.GroupValue is used instead of dot notation for accessing nested groups", func(t *testing.T) {
		t.Parallel()

		records := []slogmem.LoggedRecord{
			{
				Time:    time.Now(),
				Level:   slog.LevelDebug,
				Message: "some debug message",
				Attrs:   []slog.Attr{slog.Group("r1ka", slog.String("r1gaka", "r1gava"))},
			},
		}
		query := slogmem.RecordQuery{
			Level:   slog.LevelDebug,
			Message: "some debug message",
			Attrs: map[string]slog.Value{
				"r1ka": slog.GroupValue(slog.String("r1gaka", "r1gava")),
			},
		}

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("slogmem.NewLoggedRecords(%+v).Contains(%+v) did not panic when using slog.GroupValue instead of dot notation", records, query)
			}
		}()

		_, _ = slogmem.NewLoggedRecords(records).Contains(query)
	})

	testCases := map[string]struct {
		records []slogmem.LoggedRecord
		query   slogmem.RecordQuery
		want    bool
	}{
		"true is returned when the query matches a record when querying a single record using nil attrs": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs:   nil,
			},
			want: true,
		},
		"true is returned when the query matches a record when querying multiple records using nil attrs": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs:   nil,
			},
			want: true,
		},
		"true is returned when the query matches a record when querying a single record using an empty attrs": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs:   map[string]slog.Value{},
			},
			want: true,
		},
		"true is returned when the query matches a record when querying multiple records using an empty attrs": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs:   map[string]slog.Value{},
			},
			want: true,
		},
		"true is returned when the query with attrs matches a record when querying a single record": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   []slog.Attr{slog.String("a", "aV"), slog.String("b", "bV")},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"a": slog.StringValue("aV"),
					"b": slog.StringValue("bV"),
				},
			},
			want: true,
		},
		"true is returned when the query with attrs matches a record when querying multiple records": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   []slog.Attr{slog.String("a", "aV"), slog.String("b", "bV")},
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   []slog.Attr{slog.String("a", "aV"), slog.String("b", "bV")},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"a": slog.StringValue("aV"),
					"b": slog.StringValue("bV"),
				},
			},
			want: true,
		},
		"true is returned when the query with nested attrs matches a record when querying a single record": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   []slog.Attr{slog.Group("g", slog.String("a", "aV"), slog.Group("g2", slog.String("b", "bV")))},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"g.a":    slog.StringValue("aV"),
					"g.g2.b": slog.StringValue("bV"),
				},
			},
			want: true,
		},
		"true is returned when the query with nested attrs matches a record when querying multiple records": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   []slog.Attr{slog.Group("g", slog.String("a", "aV"), slog.Group("g2", slog.String("b", "bV")))},
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   []slog.Attr{slog.Group("g", slog.String("a", "aV"), slog.Group("g2", slog.String("b", "bV")))},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"g.a":    slog.StringValue("aV"),
					"g.g2.b": slog.StringValue("bV"),
				},
			},
			want: true,
		},
		"false is returned when the message does not match querying a single record": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some other message",
				Attrs:   nil,
			},
			want: false,
		},
		"false is returned when the message does not match querying multiple records": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelDebug,
				Message: "some other message",
				Attrs:   nil,
			},
			want: false,
		},
		"false is returned when the level does not match querying a single record": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "some message",
				Attrs:   nil,
			},
			want: false,
		},
		"false is returned when the level does not match querying multiple records": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelDebug,
					Message: "some message",
					Attrs:   nil,
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "some message",
				Attrs:   nil,
			},
			want: false,
		},
		"false is returned when the attrs do not match querying a single record": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{slog.String("a", "aV"), slog.String("b", "bV")},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"c": slog.StringValue("cV"),
					"d": slog.StringValue("dV"),
				},
			},
			want: false,
		},
		"false is returned when the attrs do not match querying multiple records": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{slog.String("a", "aV"), slog.String("b", "bV")},
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{slog.String("a", "aV"), slog.String("b", "bV")},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"a": slog.StringValue("aV"),
					"c": slog.StringValue("cV"),
				},
			},
			want: false,
		},
		"true is returned when errors logged as attributes using slog.Any are compared with a standard error": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{slog.Any("error", &logValuerStubError{})},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"error": slog.AnyValue(errors.New("some error from logValuerStubError")),
				},
			},
			want: true,
		},
		"true is returned when errors logged as attributes using slog.Any are compared with an error type": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{slog.Any("error", &logValuerStubError{I: 123})},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"error": slog.AnyValue(&logValuerStubError{I: 456}),
				},
			},
			want: true,
		},
		"true is returned when errors logged as attributes using slog.Any are compared with a string": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{slog.Any("error", errors.New("some error"))},
				},
			},
			query: slogmem.RecordQuery{
				Level:   slog.LevelInfo,
				Message: "some message",
				Attrs: map[string]slog.Value{
					"error": slog.AnyValue("some error"),
				},
			},
			want: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, gotDiff := slogmem.NewLoggedRecords(tc.records).Contains(tc.query)

			if got != tc.want {
				t.Errorf("slogmem.NewLoggedRecords(%+v).Contains(%+v) has not returned expected result:\ngot  %t\nwant %t", tc.records, tc.query, got, tc.want)
			}

			if tc.want == true && gotDiff != "" {
				t.Errorf("slogmem.NewLoggedRecords(%+v).Contains(%+v) has returned a diff unexpectedly:\ngot  %s\nwant \"\"", tc.records, tc.query, gotDiff)
			}

			if tc.want == false && gotDiff == "" {
				t.Errorf("slogmem.NewLoggedRecords(%+v).Contains(%+v) has not returned a diff", tc.records, tc.query)
			}
		})
	}
}

func TestLoggedRecordsLen(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		records []slogmem.LoggedRecord
		want    int
	}{
		"returns 0 when logged records is nil": {
			records: nil,
			want:    0,
		},
		"returns 0 when logged records is empty": {
			records: []slogmem.LoggedRecord{},
			want:    0,
		},
		"returns count when logged records is not empty": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{},
				},
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some other message",
					Attrs:   []slog.Attr{},
				},
			},
			want: 2,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := slogmem.NewLoggedRecords(tc.records).Len()

			if got != tc.want {
				t.Errorf("slogmem.NewLoggedRecords(%+v).Len() want: %d, got %d", tc.records, tc.want, got)
			}
		})
	}
}

func TestLoggedRecordsIsEmpty(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		records []slogmem.LoggedRecord
		want    bool
	}{
		"returns true when logged records is nil": {
			records: nil,
			want:    true,
		},
		"returns true when logged records is empty": {
			records: []slogmem.LoggedRecord{},
			want:    true,
		},
		"returns false when logged records is not empty": {
			records: []slogmem.LoggedRecord{
				{
					Time:    time.Now(),
					Level:   slog.LevelInfo,
					Message: "some message",
					Attrs:   []slog.Attr{},
				},
			},
			want: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := slogmem.NewLoggedRecords(tc.records).IsEmpty()

			if got != tc.want {
				t.Errorf("slogmem.NewLoggedRecords(%+v).IsEmpty() did not return %t", tc.records, tc.want)
			}
		})
	}
}

func TestHandlerRespectsCastingLogValuerWhenTestingErrors(t *testing.T) {
	t.Parallel()

	handler := slogmem.NewHandler(slog.LevelDebug)
	logger := slog.New(handler)

	logger.ErrorContext(context.Background(), "Something happened", slog.Any("error", &logValuerStubError{
		I: 123,
		S: "some value",
	}))

	query := slogmem.RecordQuery{
		Level:   slog.LevelError,
		Message: "Something happened",
		Attrs: map[string]slog.Value{
			"error.i": slog.IntValue(123),
			"error.s": slog.StringValue("some value"),
		},
	}

	if ok, diff := handler.Records().Contains(query); !ok {
		t.Errorf("handler does not respect slog.LogValuer casting, diff:\n%s", diff)
	}
}
