package logger

import (
	"context"
	"log/slog"

	loggerv1 "github.com/omcrgnt/proto/gen/go/logger/v1"
)

// Logger is the logging API.
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

var active Logger = noopLogger{}

// Config configures the process logger (level, format).
// ecfg fills and protovalidates Level and Format before Build is called.
type Config struct {
	Level  *loggerv1.Level  `ecfg:"LEVEL"`
	Format *loggerv1.Format `ecfg:"FORMAT"`
}

// Build returns a Logger for res.Add and SDI.
func (c Config) Build() (any, error) {
	return &logger{
		level:  c.Level.GetValue(),
		format: c.Format.GetValue(),
	}, nil
}

// DefaultLogConfig returns the system Logger config for logger/use registration.
func DefaultLogConfig() Config {
	return Config{
		Level:  &loggerv1.Level{Value: defaultLevelValue},
		Format: &loggerv1.Format{Value: defaultFormatValue},
	}
}

// DefaultLog returns the system Logger resource for tests and legacy callers.
func DefaultLog() any {
	return &logger{level: defaultLevelValue, format: defaultFormatValue}
}

// logger implements [Logger] and lives in res. Created by [DefaultLog] or [Config.Build].
type logger struct {
	level  string
	format string
	out    Output
	slog   *slog.Logger
}

func (l *logger) Deps() []any {
	return []any{(*Output)(nil)}
}

func (l *logger) Inject(args []any) {
	for _, arg := range args {
		if out, ok := arg.(Output); ok {
			l.out = out
			break
		}
	}
	if l.out == nil {
		return
	}
	l.slog = buildSlog(l.level, l.format, l.out)
	active = l
}

func (l *logger) Debug(ctx context.Context, msg string, args ...any) {
	if l.slog == nil {
		return
	}
	l.slog.DebugContext(ctx, msg, args...)
}

func (l *logger) Info(ctx context.Context, msg string, args ...any) {
	if l.slog == nil {
		return
	}
	l.slog.InfoContext(ctx, msg, args...)
}

func (l *logger) Warn(ctx context.Context, msg string, args ...any) {
	if l.slog == nil {
		return
	}
	l.slog.WarnContext(ctx, msg, args...)
}

func (l *logger) Error(ctx context.Context, msg string, args ...any) {
	if l.slog == nil {
		return
	}
	l.slog.ErrorContext(ctx, msg, args...)
}

// Default returns the process logger wired by sdi, or a no-op logger before [sdi.Resolve].
func Default() Logger {
	return active
}
