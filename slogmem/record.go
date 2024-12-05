package slogmem

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type (
	// LoggedRecord encapsulates the information that was recorded in a single
	// log entry by the [Handler].
	LoggedRecord struct {
		// Time is the time the log was written.
		Time time.Time
		// Level is the [slog.Level] that the log was written as.
		Level slog.Level
		// Message is the message that was passed by the caller for the given log entry.
		Message string
		// Attrs is a slice of [slog.Attr] records that represent the additional
		// attributes that were added to the log entry by the caller as context.
		Attrs []slog.Attr
	}

	// LoggedRecords is a slice of [LoggedRecord] entries that were captured by a [Handler].
	// Adding to LoggedRecords is safe to do concurrently.
	LoggedRecords struct {
		mu      sync.Mutex
		records []LoggedRecord
	}

	// RecordQuery represents the relevant information required in order to query for
	// the existence of a [LoggedRecord] within a set of [LoggedRecords]. Time is not
	// part of the query as it is generally difficult to know when the log was
	// written in order to query for it accurately.
	RecordQuery struct {
		// Level is the [slog.Level] that the log was written as.
		Level slog.Level
		// Message is the message that was passed by the caller for the given log entry.
		Message string
		// Attrs is a map of dot separated keys that each indicate a path to a grouped
		// attribute and the value of that attribute. For example: if an attribute was
		// written as `slog.Group("group", slog.String("key", "value"))` then to query
		// that, we would pass `map[string]slog.Value{"group.key": slog.StringValue("value")}`.
		Attrs map[string]slog.Value
	}
)

// NewLoggedRecords encapsulates the given list of [LoggedRecord] entries within
// a LoggedRecords struct to represent the list of logged records in a way that
// is easy to lookup when asserting logs in tests or similar.
func NewLoggedRecords(records []LoggedRecord) *LoggedRecords {
	return &LoggedRecords{
		mu:      sync.Mutex{},
		records: records,
	}
}

// Contains can be used to check if a LoggedRecords contains a
// [LoggedRecord] that matches the details in the given [RecordQuery].
//
// Contains returns true for the first record that fully matches the given
// [RecordQuery]. If there are no records matching the given query then false
// will be returned.
//
// There is an additional argument returned as a convenience helper which
// provides a diff of the passed [RecordQuery] and the [LoggedRecords] when there
// is no match or a partial match on the message. When Contains returns true, the
// diff (second return value) will be empty. When there is not a full match on a
// record but the message matches, a diff will be produced for each record
// matching that message. Otherwise, a diff over all records that were logged
// will be produced.
//
// NOTE: this diff is nondeterministic, do not rely on its output. This is a
// convenience helper for logging information in failed tests and similar
// scenarios.
func (lr *LoggedRecords) Contains(query RecordQuery) (bool, string) {
	diff, msgMatchDiff := "", ""

	flattenedQuery := flattenRecordQuery(query)

	for i, flattenedRecord := range lr.AsSliceOfNestedKeyValuePairs() {
		if cmp.Equal(flattenedQuery, flattenedRecord, cmpOpts()...) {
			return true, ""
		}

		recordDiff := cmp.Diff(flattenedQuery, flattenedRecord, cmpOpts()...)

		if lr.records[i].Message == query.Message {
			msgMatchDiff += fmt.Sprintln(recordDiff)
		}

		diff += fmt.Sprintln(recordDiff)
	}

	if msgMatchDiff != "" {
		return false, msgMatchDiff
	}

	return false, diff
}

// IsEmpty returns true when no records have been captured.
func (lr *LoggedRecords) IsEmpty() bool { return lr.Len() == 0 }

// Len returns the number of records that have been captured.
func (lr *LoggedRecords) Len() int { return len(lr.records) }

// AsSliceOfNestedKeyValuePairs flattens the LoggedRecords so that they can be
// accessed as a series of key value pair objects representing each recorded log.
//
// This method would be used when formatting the recorded log records as JSON for
// example.
func (lr *LoggedRecords) AsSliceOfNestedKeyValuePairs() []map[string]any {
	const numBaseAttrs = 3 // time, level, message

	flattenedRecords := make([]map[string]any, 0, len(lr.records))

	for _, rec := range lr.records {
		flattenedRecord := make(map[string]any, numBaseAttrs+len(rec.Attrs))

		flattenedRecord[slog.TimeKey] = rec.Time
		flattenedRecord[slog.LevelKey] = rec.Level
		flattenedRecord[slog.MessageKey] = rec.Message

		for _, attr := range rec.Attrs {
			mapAttr(flattenedRecord, attr)
		}

		flattenedRecords = append(flattenedRecords, flattenedRecord)
	}

	return flattenedRecords
}

// append safely appends a [LoggedRecord] to the list of LoggedRecords.
func (lr *LoggedRecords) append(record LoggedRecord) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	lr.records = append(lr.records, record)
}

// flattenRecordQuery is used to flatten a RecordQuery into map for comparison with flattened [LoggedRecords].
func flattenRecordQuery(recordQuery RecordQuery) map[string]any {
	const numBaseAttrs = 2
	flattenedRecordQuery := make(map[string]any, numBaseAttrs+len(recordQuery.Attrs))

	flattenedRecordQuery[slog.LevelKey] = recordQuery.Level
	flattenedRecordQuery[slog.MessageKey] = recordQuery.Message

	for path, value := range recordQuery.Attrs {
		recursiveSetField(flattenedRecordQuery, path, value)
	}

	return flattenedRecordQuery
}

// recursiveSetField sets a field in the given map to the value based on a dot separated fieldPath.
func recursiveSetField(record map[string]any, fieldPath string, value slog.Value) {
	keys := strings.Split(fieldPath, ".")
	currentKey := keys[0]
	remainingKeys := keys[1:]
	remainingPath := strings.Join(remainingKeys, ".")

	if len(keys) == 1 {
		if value.Kind() == slog.KindGroup {
			panic("slog.GroupValue cannot be used as a value when checking attrs, for nested attrs use dot notation instead")
		}

		record[currentKey] = value.Any()

		return
	}

	if currentKeyValue, ok := record[currentKey]; ok {
		ckv, ckvOK := currentKeyValue.(map[string]any)

		// It should be impossible to hit this panic as currentKeyValue cannot be set
		// outside the code in this package. We know the type will always be either a
		// map[string]any or handled above via the value slog.Value passed into the
		// function.
		if !ckvOK {
			panic("unexpected attr value type")
		}

		recursiveSetField(ckv, remainingPath, value)

		return
	}

	nestedGroup := make(map[string]any)
	recursiveSetField(nestedGroup, remainingPath, value)
	record[currentKey] = nestedGroup
}

// mapAttr unpacks any slog.KindGroup attrs and converts the values to any.
func mapAttr(record map[string]any, attr slog.Attr) {
	if attr.Value.Kind() != slog.KindGroup {
		record[attr.Key] = attr.Value.Any()
		return
	}

	mappedGroup := make(map[string]any, len(attr.Value.Group()))
	for _, groupedAttr := range attr.Value.Group() {
		mapAttr(mappedGroup, groupedAttr)
	}

	record[attr.Key] = mappedGroup
}

// cmpOpts is used to ignore Time values and compare errors when using cmp for
// comparisons and diffs internally. There may be a need/want to export this in
// the future but for now it is unexported for better encapsulation.
func cmpOpts() []cmp.Option {
	return []cmp.Option{
		cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
			return k == slog.TimeKey
		}),
		cmpopts.EquateErrors(),
	}
}
