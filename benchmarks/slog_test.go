package benchmarks

import (
	"io"
	"log/slog"
)

func newSlog(fields ...slog.Attr) *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil).WithAttrs(fields))
}

func newDisabledSlog(fields ...slog.Attr) *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}).WithAttrs(fields))
}

func fakeSlogFields() []slog.Attr {
	return []slog.Attr{
		slog.Int("int", _tenInts[0]),
		slog.Any("ints", _tenInts),
		slog.String("string", _tenStrings[0]),
		slog.Any("strings", _tenStrings),
		slog.Time("time", _tenTimes[0]),
		slog.Any("times", _tenTimes),
		slog.Any("user1", _oneUser),
		slog.Any("user2", _oneUser),
		slog.Any("users", _tenUsers),
		slog.Any("error", errExample),
	}
}

func fakeSlogArgs() []any {
	return []any{
		"int", _tenInts[0],
		"ints", _tenInts,
		"string", _tenStrings[0],
		"strings", _tenStrings,
		"time", _tenTimes[0],
		"times", _tenTimes,
		"user1", _oneUser,
		"user2", _oneUser,
		"users", _tenUsers,
		"error", errExample,
	}
}
