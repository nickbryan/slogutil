package slogutil

import (
	"io"
	"log/slog"
	"os"
	"time"
)

type (
	// Option is an optional configuration value used to configure a logger.
	Option func(*options)

	options struct {
		level     slog.Leveler
		addSource bool
		now       func() time.Time
		writer    io.Writer
	}
)

// TimeFactoryFunc represents a function that knows how to create [time.Time] values
// to be used by the logger when setting the time value of the log.
type TimeFactoryFunc func() time.Time

// WithLevel will set the log level. The default is [slog.LevelInfo].
func WithLevel(level slog.Leveler) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithSourceAdded sets [slog.HandlerOptions.AddSource]. The default is true.
func WithSourceAdded(addSource bool) Option {
	return func(o *options) {
		o.addSource = addSource
	}
}

// WithTimeFactory sets the [TimeFactoryFunc] on the Logger that will be
// used to determine the time values of the logs. This is useful in rare situations
// when simulation of time is required, for example, in example test functions.
// Prefer using the in-memory handler where possible. The default is [time.Now()].
func WithTimeFactory(factory TimeFactoryFunc) Option {
	return func(o *options) {
		o.now = factory
	}
}

// WithWriter sets the [io.Writer] that the logs are written to. The default is
// [io.Stderr].
func WithWriter(writer io.Writer) Option {
	return func(o *options) {
		o.writer = writer
	}
}

func mapOptionsToDefaults(opts []Option) options {
	mappedDefaultOpts := options{
		level:     slog.LevelInfo,
		addSource: true,
		now:       nil,
		writer:    os.Stderr,
	}

	for _, opt := range opts {
		opt(&mappedDefaultOpts)
	}

	return mappedDefaultOpts
}
