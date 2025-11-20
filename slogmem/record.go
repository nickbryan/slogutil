package slogmem

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"
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

// Contains can be used to check if a LoggedRecords contains a [LoggedRecord]
// that matches the details in the given [RecordQuery].
//
// Contains returns true for the first record that matches the given
// [RecordQuery]. If there are no records matching the given query, then false
// will be returned. A loose match is performed on the attributes in the query.
// Any additional attributes in the record will not be checked.
//
// There is an additional argument returned as a convenience helper that
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
func (lr *LoggedRecords) Contains(query RecordQuery) (ok bool, diff string) {
	paths := append(slices.Collect(maps.Keys(query.Attrs)), slog.MessageKey, slog.LevelKey)
	return lr.compare(query, append(cmpOpts(), includePaths(paths))...)
}

// ContainsExact can be used to check if a LoggedRecords contains a
// [LoggedRecord] that is an exact match of the details in the given
// [RecordQuery].
//
// ContainsExact returns true for the first record that fully matches the given
// [RecordQuery]. If there are no records matching the given query, then false
// will be returned. All attributes in the query must be present in the record.
//
// There is an additional argument returned as a convenience helper that
// provides a diff of the passed [RecordQuery] and the [LoggedRecords] when there
// is no match or a partial match on the message. When ContainsExact returns true, the
// diff (second return value) will be empty. When there is not a full match on a
// record but the message matches, a diff will be produced for each record
// matching that message. Otherwise, a diff over all records that were logged
// will be produced.
//
// NOTE: this diff is nondeterministic, do not rely on its output. This is a
// convenience helper for logging information in failed tests and similar
// scenarios.
func (lr *LoggedRecords) ContainsExact(query RecordQuery) (ok bool, diff string) {
	return lr.compare(query, cmpOpts()...)
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

func (lr *LoggedRecords) compare(query RecordQuery, opts ...cmp.Option) (bool, string) {
	var (
		diff         strings.Builder
		msgMatchDiff strings.Builder
	)

	flattenedQuery := flattenRecordQuery(query)

	for i, flattenedRecord := range lr.AsSliceOfNestedKeyValuePairs() {
		if cmp.Equal(flattenedQuery, flattenedRecord, opts...) {
			return true, ""
		}

		recordDiff := cmp.Diff(flattenedQuery, flattenedRecord, opts...)

		if lr.records[i].Message == query.Message {
			msgMatchDiff.WriteString(fmt.Sprintln(recordDiff))
		}

		diff.WriteString(fmt.Sprintln(recordDiff))
	}

	if msgMatchDiff.Len() > 0 {
		return false, msgMatchDiff.String()
	}

	return false, diff.String()
}

// includePaths returns a cmp.Option that will ignore any paths that do not match the given paths.
func includePaths(paths []string) cmp.Option { //nolint:ireturn // We need to return a cmp.Option here which is an interface.
	include := make([][]string, 0, len(paths))

	for _, p := range paths {
		include = append(include, strings.Split(p, "."))
	}

	return cmp.FilterPath(func(path cmp.Path) bool {
		currentPath := pathToStrings(path)
		if len(currentPath) == 0 {
			return false
		}

		for _, pathToInclude := range include {
			// If the current path is longer than the path to include, it can't be a prefix.
			if len(currentPath) > len(pathToInclude) {
				continue
			}

			// If the current path is a prefix of the path to include, it's a match or a parent, so don't ignore it.
			if cmp.Equal(currentPath, pathToInclude[:len(currentPath)]) {
				return false
			}
		}

		return true // If we got here, the path is not in the include list, so ignore it.
	}, cmp.Ignore()) // Filter out the values that do not match the given paths.
}

// pathToStrings converts a cmp.Path to a slice of strings.
func pathToStrings(path cmp.Path) []string {
	var parts []string

	for _, step := range path {
		if s, ok := step.(cmp.MapIndex); ok {
			parts = append(parts, s.Key().String())
		}
	}

	return parts
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
		cmp.FilterValues(areConcreteErrors, cmp.Comparer(compareErrorStrings)),
		cmp.FilterValues(isStringAndError, cmp.Comparer(compareStringAndError)),
	}
}

// areConcreteErrors reports whether x and y are types that implement error.
// The input types are deliberately of the interface{} type rather than the
// error type so that we can handle situations where the current type is an
// interface{}, but the underlying concrete types both happen to implement
// the error interface.
func areConcreteErrors(x, y interface{}) bool {
	_, ok1 := x.(error)
	_, ok2 := y.(error)

	return ok1 && ok2
}

// compareErrorStrings is used to compare the strings of the errors rather than the types.
// cmp.Diff will produce a value if two errors are logged with the same type but their
// memory address is different (different instances).
func compareErrorStrings(x, y any) bool {
	xAsErr, _ := x.(error)
	yAsErr, _ := y.(error)

	return xAsErr.Error() == yAsErr.Error()
}

// isStringAndError reports whether x is a string and y is an error, or vice versa.
func isStringAndError(x, y any) bool {
	_, xIsString := x.(string)
	_, yIsString := y.(string)
	_, xIsError := x.(error)
	_, yIsError := y.(error)

	return (xIsString && yIsError) || (yIsString && xIsError)
}

// compareStringAndError compares a string to an error's Error() string.
func compareStringAndError(x, y any) bool {
	var (
		str string
		err error
	)

	if s, ok := x.(string); ok {
		str = s
		err, _ = y.(error)
	} else {
		str, _ = y.(string)
		err, _ = x.(error)
	}

	return str == err.Error()
}
