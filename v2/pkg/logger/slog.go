package logger

import "log/slog"

const (
	defaultLevelValue  = "info"
	defaultFormatValue = "json"
)

func buildSlog(level, format string, out Output) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}
	switch format {
	case "text":
		return slog.New(slog.NewTextHandler(out, opts))
	default:
		return slog.New(slog.NewJSONHandler(out, opts))
	}
}

func parseLevel(value string) slog.Leveler {
	switch value {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
