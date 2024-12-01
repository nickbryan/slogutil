// Package slogmem provides a [slog.Handler] that captures log records in memory.
package slogmem

import (
	"context"
	"log/slog"

	"github.com/nickbryan/slogutil/internal"
)

// Handler captures records produced by a call to Handle in-memory so that they can be
// accessed via [LoggedRecords] later for inspection.
type Handler struct {
	persistentAttrs internal.AttrGroupTree
	leveler         slog.Leveler
	loggedRecords   *LoggedRecords
}

// Ensure that our [Handler] implements the [slog.Handler] interface.
var _ slog.Handler = &Handler{} //nolint:exhaustruct // Compile type implementation check.

// NewHandler creates a new in-memory Handler that captures log records which have a
// level greater than or equal to the current level of the given leveler.
func NewHandler(leveler slog.Leveler) *Handler {
	return &Handler{
		persistentAttrs: internal.NewAttrGroupTree(),
		leveler:         leveler,
		loggedRecords:   NewLoggedRecords(make([]LoggedRecord, 0)),
	}
}

// WithAttrs returns a new Handler whose attributes consist of both the existing
// handler's attributes and those given.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		persistentAttrs: h.persistentAttrs.WithAttrs(attrs),
		leveler:         h.leveler,
		loggedRecords:   h.loggedRecords,
	}
}

// WithGroup returns a new Handler that will store all future attributes under a
// group with the given name.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		persistentAttrs: h.persistentAttrs.WithGroup(name),
		leveler:         h.leveler,
		loggedRecords:   h.loggedRecords,
	}
}

// Enabled returns whether the Handler is currently enabled for the given [slog.Level].
// Levels greater than or equal to that of the [Handler]'s [slog.Leveler]'s current
// Level are considered enabled.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.leveler.Level()
}

// Records returns the [LoggedRecords] that were recorded by this Handler.
func (h *Handler) Records() *LoggedRecords {
	return h.loggedRecords
}

// Handle stores the content of the [slog.Record] into the Handler's in-memory
// [LoggedRecords] store.
//
// Handle will only be called when [Enabled] returns true.
func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	recordAttrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		recordAttrs = append(recordAttrs, attr)
		return true
	})

	h.loggedRecords.append(LoggedRecord{
		Time:    record.Time,
		Level:   record.Level,
		Message: record.Message,
		Attrs:   h.persistentAttrs.WithAttrs(recordAttrs).History().DeduplicatedAttrs(),
	})

	return nil
}
