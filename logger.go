package seqlogger

import (
	"log/slog"
)

func New(config Config) *slog.Logger {
	return slog.New(NewSeqHandler(config))
}

func NewWithHandler(handler *SeqHandler) *slog.Logger {
	return slog.New(handler)
}

func MapLogLevel(level slog.Level) string {
	switch {
	case level <= slog.LevelDebug:
		return "Debug"
	case level <= slog.LevelInfo:
		return "Information"
	case level <= slog.LevelWarn:
		return "Warning"
	case level <= slog.LevelError:
		return "Error"
	default:
		return "Fatal"
	}
}
