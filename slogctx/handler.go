// Package slogctx provides a context aware [slog.Handler].
package slogctx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nickbryan/slogutil/internal"
)

// Handler extracts attributes from a [context.Context] where they have been
// added via the functions [WithRootAttrs] or [WithAttrs]. All extracted attributes
// will be passed to the embedded [slog.Handler] for further processing.
type Handler struct {
	slog.Handler

	persistentAttrs    internal.AttrGroupTree
	attrExtractors     []Extractor
	rootAttrExtractors []Extractor
}

// Ensure that our [Handler] implements the [slog.Handler] interface.
var _ slog.Handler = &Handler{} //nolint:exhaustruct // Compile time implementation check.

// NewHandler creates a new Handler that extracts attributes from
// [context.Context] where they have been added via the functions
// [WithRootAttrs] and [WithAttrs].
//
// All extracted attributes will be passed to the wrapped [slog.Handler] for
// further processing.
func NewHandler(wrapped slog.Handler) *Handler {
	h := &Handler{
		Handler:            wrapped,
		persistentAttrs:    internal.NewAttrGroupTree(),
		attrExtractors:     make([]Extractor, 0, 1),
		rootAttrExtractors: make([]Extractor, 0, 1),
	}

	h.AddAttrExtractors(newCtxExtractor(ctxKeyWithAttrs{}))
	h.AddRootAttrExtractors(newCtxExtractor(ctxKeyWithRootAttrs{}))

	return h
}

// WithAttrs returns a new Handler whose attributes consist of both the existing
// handler's attributes and those given. If attrs is empty, the existing Handler
// will be returned.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		Handler:            h.Handler,
		persistentAttrs:    h.persistentAttrs.WithAttrs(attrs),
		attrExtractors:     h.attrExtractors,
		rootAttrExtractors: h.rootAttrExtractors,
	}
}

// WithGroup returns a new Handler that will store all future attributes under a
// group with the given name. If name is empty, the receiver Handler will be
// returned and attributes will be stored on the current group, if there is one.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		Handler:            h.Handler,
		persistentAttrs:    h.persistentAttrs.WithGroup(name),
		attrExtractors:     h.attrExtractors,
		rootAttrExtractors: h.rootAttrExtractors,
	}
}

// AddAttrExtractors adds the given list of [Extractor]s
// to the list of [Extractor]s that will run after all other attrs
// have been added to the log record.
func (h *Handler) AddAttrExtractors(extractors ...Extractor) {
	h.attrExtractors = append(h.attrExtractors, extractors...)
}

// AddRootAttrExtractors adds the given list of [Extractor]s
// to the list of [Extractor]s that will run before all other attrs
// have been added to the log record adding them to the root of the
// log record.
func (h *Handler) AddRootAttrExtractors(extractors ...Extractor) {
	h.rootAttrExtractors = append(h.rootAttrExtractors, extractors...)
}

// Handle will extract attributes from [context.Context] where they have been
// added via the functions [WithRootAttrs] and [WithAttrs]. All
// extracted attributes will be passed to the embedded logger for further
// processing.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	// Attributes are ordered as: withRootAttrs, groupedAttrs, recordAttrs, withAttrs
	recordAttrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		recordAttrs = append(recordAttrs, attr)
		return true
	})

	for _, extractor := range h.attrExtractors {
		if attrs := extractor.Extract(ctx); attrs != nil {
			recordAttrs = append(recordAttrs, attrs...)
		}
	}

	// When adding to the root, we order first to ensure we have a scoped copy so that we do not affect other loggers.
	orderedRecordedAttrs := h.persistentAttrs.WithAttrs(recordAttrs).History()

	for _, extractor := range h.rootAttrExtractors {
		if attrs := extractor.Extract(ctx); attrs != nil {
			orderedRecordedAttrs.PushFront(attrs)
		}
	}

	record = slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.AddAttrs(orderedRecordedAttrs.DeduplicatedAttrs()...)

	if err := h.Handler.Handle(ctx, record); err != nil {
		return fmt.Errorf("passing record to inner handler: %w", err)
	}

	return nil
}
